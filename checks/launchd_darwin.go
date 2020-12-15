package checks

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/it-novum/openitcockpit-agent-go/utils"
)

// CheckLaunchd gathers information about launchd and launchctl services
type CheckLaunchd struct {
}

// Name will be used in the response as check name
func (c *CheckLaunchd) Name() string {
	return "launchd_services"
}

type resultLaunchdServices struct {
	IsRunning bool
	Pid       int
	Status    int
	Label     string
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckLaunchd) Run(ctx context.Context) (*CheckResult, error) {
	ctx, cancel := context.WithCancel(context.Background())
	timeout := 10 * time.Second

	/* From the man page of launchctl list (macOS 10.15.7):
	 * list [-x] [label]
	 *    With no arguments, list all of the jobs loaded into launchd in three columns. The first column displays the PID of the job if it is running.  The
	 *    second column displays the last exit status of the job. If the number in this column is negative, it represents the negative of the signal which
	 *    stopped the job. Thus, "-15" would indicate that the job was terminated with SIGTERM.  The third column is the job's label. If [label] is specified,
	 *    prints information about the requested job.
	 *
	 *    -x       This flag is no longer supported.
	 */
	result, err := utils.RunCommand(ctx, "launchctl list", timeout)
	cancel()
	if err != nil || result.RC > 0 {
		fmt.Println("Error while executing 'launchctl list'")
		return nil, err
	}

	lines := strings.Split(result.Stdout, "\n")
	if len(lines) > 0 {
		//Remove first line
		lines = lines[1:]
	}

	launchdResults := make([]*resultLaunchdServices, 0, len(lines))
	for _, line := range lines {
		columns := strings.Split(line, "\t")
		if len(columns) == 3 {
			isRunning := false
			pid, err := strconv.Atoi(strings.TrimSpace(columns[0]))
			if err == nil {
				//pid is a number an not "-"
				isRunning = true
			}

			status, _ := strconv.Atoi(strings.TrimSpace(columns[1]))
			label := strings.TrimSpace(columns[2])

			result := &resultLaunchdServices{
				IsRunning: isRunning,
				Pid:       pid,
				Status:    status,
				Label:     label,
			}
			launchdResults = append(launchdResults, result)
		}
	}

	return &CheckResult{Result: launchdResults}, nil

}

// DefaultConfiguration contains the variables for the configuration file and the default values
// can be nil if no configuration is required
func (c *CheckLaunchd) DefaultConfiguration() interface{} {
	return nil
}

// Configure should verify the configuration and set it
// will be run after every reload
// if DefaultConfiguration returns nil, the parameter will also be nil
func (c *CheckLaunchd) Configure(_ interface{}) error {
	return nil
}
