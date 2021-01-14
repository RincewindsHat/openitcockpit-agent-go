package checks

import (
	"context"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckDiskIo) Run(ctx context.Context) (interface{}, error) {

	disks, err := disk.IOCountersWithContext(ctx)
	if err != nil {
		return nil, err
	}
	diskResults := make(map[string]*resultDiskIo)

	for device, iostats := range disks {

		if lastCheckResults, ok := c.lastResults[disk.Name]; ok {
			ReadCount, _ := Wrapdiff(float64(lastCheckResults.ReadCount), float64(iostats.ReadCount))
			WriteCount, _ := Wrapdiff(float64(lastCheckResults.WriteCount), float64(iostats.WriteCount))
			IoTime, _ := Wrapdiff(float64(lastCheckResults.IoTime), float64(iostats.IoTime))
			ReadTime, _ := Wrapdiff(float64(lastCheckResults.ReadTime), float64(iostats.ReadTime))
			WriteTime, _ := Wrapdiff(float64(lastCheckResults.WriteTime), float64(iostats.WriteTime))
			ReadBytes, _ := Wrapdiff(float64(lastCheckResults.ReadBytes), float64(iostats.ReadBytes))
			WriteBytes, _ := Wrapdiff(float64(lastCheckResults.WriteBytes), float64(iostats.WriteBytes))
			Timestamp, _ := Wrapdiff(float64(lastCheckResults.Timestamp), float64(time.Now().Unix()))

			loadPercent := IoTime / (Timestamp * 1000) * 100

			readAvgWait := ReadTime / ReadCount
			readAvgSize := ReadBytes / ReadCount

			writeAvgWait := WriteTime / WriteCount
			writeAvgSize := WriteBytes / WriteCount

			totIos := ReadCount + WriteCount
			totalAvgWait := (ReadTime + WriteTime) / totIos

			if loadPercent <= 101 {
				// Just in case this this has the same bug as Python psutil has^^
				diskstats := &resultDiskIo{
					Timestamp:    time.Now().Unix(),
					ReadBytes:    uint64(ReadCount),
					WriteBytes:   uint64(WriteBytes),
					ReadIops:     uint64(ReadCount),
					WriteIops:    uint64(WriteCount),
					TotalIops:    uint64(totIos),
					ReadCount:    uint64(ReadCount),
					WriteCount:   uint64(WriteCount),
					IoTime:       uint64(IoTime),
					ReadAvgWait:  readAvgWait,
					ReadTime:     uint64(ReadTime),
					ReadAvgSize:  readAvgSize,
					WriteAvgWait: writeAvgWait,
					WriteAvgSize: writeAvgSize,
					WriteTime:    uint64(WriteTime),
					TotalAvgWait: totalAvgWait,
					LoadPercent:  int64(loadPercent),
					Device:       device,
				}

				diskResults[disk.Name] = diskstats
			}

		} else {
			//No previous check results for calculations... wait until check runs again
			diskstats := &resultDiskIo{
				ReadCount:  iostats.ReadCount,
				WriteCount: iostats.WriteCount,
				IoTime:     iostats.IoTime,
				ReadTime:   iostats.ReadTime,
				WriteTime:  iostats.WriteTime,
				ReadBytes:  iostats.ReadBytes,
				WriteBytes: iostats.WriteBytes,
				Timestamp:  time.Now().Unix(),
				Device:     device,
			}

			//Store result for next check run
			diskResults[disk.Name] = diskstats
		}

	}

	c.lastResults = diskResults
	return diskResults, nil
}
