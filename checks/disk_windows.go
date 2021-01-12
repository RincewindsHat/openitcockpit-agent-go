package checks

import (
	"context"
	"fmt"
	"strconv"

	"github.com/it-novum/openitcockpit-agent-go/utils"
	"github.com/leoluk/perflib_exporter/perflib"
)

// Credit to:
// https://github.com/prometheus-community/windows_exporter/blob/9723aa221885f593ac77019566c1ced9d4d746fd/collector/logical_disk.go#L168-L187
// https://docs.microsoft.com/de-de/windows/win32/wmisdk/retrieving-raw-and-formatted-performance-data?redirectedfrom=MSDN
// https://msdn.microsoft.com/en-us/library/ms803973.aspx - LogicalDisk object reference
type Perf_LogicalDisk struct {
	Name                   string
	CurrentDiskQueueLength float64 `perflib:"Current Disk Queue Length"` // Type: Gauge
	DiskReadBytesPerSec    float64 `perflib:"Disk Read Bytes/sec"`       // Type: Counter
	DiskReadsPerSec        float64 `perflib:"Disk Reads/sec"`            // Type: Counter
	DiskWriteBytesPerSec   float64 `perflib:"Disk Write Bytes/sec"`      // Type: Counter
	DiskWritesPerSec       float64 `perflib:"Disk Writes/sec"`           // Type: Counter
	PercentDiskReadTime    float64 `perflib:"% Disk Read Time"`          // Type: Counter
	PercentDiskWriteTime   float64 `perflib:"% Disk Write Time"`         // Type: Counter
	PercentFreeSpace       float64 `perflib:"% Free Space_Base"`         // Type: Gauge - Total disk space in MB (yes) https://docs.microsoft.com/en-us/previous-versions/windows/embedded/ms938601(v=msdn.10)
	PercentFreeSpace_Base  float64 `perflib:"Free Megabytes"`            // Type: Gauge - Free disk space in MB
	PercentIdleTime        float64 `perflib:"% Idle Time"`               // Type: Counter
	SplitIOPerSec          float64 `perflib:"Split IO/Sec"`              // Type: Counter
	AvgDiskSecPerRead      float64 `perflib:"Avg. Disk sec/Read"`        // Type: Counter
	AvgDiskSecPerWrite     float64 `perflib:"Avg. Disk sec/Write"`       // Type: Counter
	AvgDiskSecPerTransfer  float64 `perflib:"Avg. Disk sec/Transfer"`    // Type: Counter
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckDisk) Run(ctx context.Context) (interface{}, error) {

	//Todo can we cache this?
	nametable := perflib.QueryNameTable("Counter 009")
	query := strconv.FormatUint(uint64(nametable.LookupIndex("LogicalDisk")), 10)

	objects, err := perflib.QueryPerformanceData(query)
	diskResults := make([]*resultDisk, 0)
	if err != nil {
		return nil, err
	}

	for _, obj := range objects {
		if obj.Name != "LogicalDisk" {
			continue
		}

		var dst []Perf_LogicalDisk
		err = utils.UnmarshalObject(obj, &dst)
		if err != nil {
			// todo add logging
			fmt.Println(err)
			continue
		}

		for _, disk := range dst {

			//Do the math
			totalDiskSpaceBytes := disk.PercentFreeSpace * 1024 * 1024
			freeDiskSpaceBytes := disk.PercentFreeSpace_Base * 1024 * 1024
			usedDiskSpaceBytes := totalDiskSpaceBytes - freeDiskSpaceBytes

			freeDiskSpacePercentage := freeDiskSpaceBytes / totalDiskSpaceBytes * 100.0
			usedDiskSpacePercentage := 100.0 - freeDiskSpacePercentage

			//Save to struct
			result := &resultDisk{}

			result.Disk.Device = disk.Name
			result.Disk.Mountpoint = disk.Name

			result.Usage.Total = uint64(totalDiskSpaceBytes)
			result.Usage.Used = uint64(usedDiskSpaceBytes)
			result.Usage.Free = uint64(freeDiskSpaceBytes)
			result.Usage.Percent = usedDiskSpacePercentage

			diskResults = append(diskResults, result)
		}
	}

	return diskResults, nil
}
