package checks

import (
	"context"
	"time"

	"github.com/StackExchange/wmi"
)

// WMI Structs
// https://docs.microsoft.com/en-us/windows/win32/cimwin32prov/win32-networkadapter
type Win32_NetworkAdapter struct {
	AdapterType                 string
	AdapterTypeID               uint16
	AutoSense                   bool
	Availability                uint16
	Caption                     string
	ConfigManagerErrorCode      uint32
	ConfigManagerUserConfig     bool
	CreationClassName           string
	Description                 string
	DeviceID                    string
	ErrorCleared                bool
	ErrorDescription            string
	GUID                        string
	Index                       uint32
	InstallDate                 time.Time
	Installed                   bool
	InterfaceIndex              uint32
	LastErrorCode               uint32
	MACAddress                  string
	Manufacturer                string
	MaxNumberControlled         uint32
	MaxSpeed                    uint64
	Name                        string
	NetConnectionID             string
	NetConnectionStatus         uint16
	NetEnabled                  bool
	NetworkAddresses            []string
	PermanentAddress            string
	PhysicalAdapter             bool
	PNPDeviceID                 string
	PowerManagementCapabilities []uint16
	PowerManagementSupported    bool
	ProductName                 string
	ServiceName                 string
	Speed                       uint64
	Status                      string
	StatusInfo                  uint16
	SystemCreationClassName     string
	SystemName                  string
	TimeOfLastReset             time.Time
}

// https://docs.microsoft.com/en-us/previous-versions/aa394293(v=vs.85)
type Win32_PerfFormattedData_Tcpip_NetworkInterface struct {
	BytesReceivedPerSec             uint32
	BytesSentPerSec                 uint32
	BytesTotalPerSec                uint64
	Caption                         string
	CurrentBandwidth                uint32
	Description                     string
	Frequency_Object                uint64
	Frequency_PerfTime              uint64
	Frequency_Sys100NS              uint64
	Name                            string
	OutputQueueLength               uint32
	PacketsOutboundDiscarded        uint32
	PacketsOutboundErrors           uint32
	PacketsPerSec                   uint32
	PacketsReceivedDiscarded        uint32
	PacketsReceivedErrors           uint32
	PacketsReceivedNonUnicastPerSec uint32
	PacketsReceivedPerSec           uint32
	PacketsReceivedUnicastPerSec    uint32
	PacketsReceivedUnknown          uint32
	PacketsSentNonUnicastPerSec     uint32
	PacketsSentPerSec               uint32
	PacketsSentUnicastPerSec        uint32
	Timestamp_Object                uint64
	Timestamp_PerfTime              uint64
	Timestamp_Sys100NS              uint64
}

type MSFT_NetAdapter struct {
	Name              string
	Status            string
	FullDuplex        bool
	MediaDuplexState  uint32
	MtuSize           uint32
	VlanID            uint16
	TransmitLinkSpeed uint64
	ReceiveLinkSpeed  uint64
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckNet) Run(ctx context.Context) (interface{}, error) {
	var dst []Win32_NetworkAdapter
	// Will return Intel(R) Ethernet Connection (2) I219-V as interface names
	err := wmi.Query("SELECT * FROM Win32_NetworkAdapter", &dst)
	if err != nil {
		return nil, err
	}

	//	var dstTwo []Win32_PerfFormattedData_Tcpip_NetworkInterface
	//	// Win32_PerfFormattedData_Tcpip_NetworkAdapter is a hidden feature??
	//	// I cant find any MS docs about
	//	// Irina found this GitHub issue https://github.com/opserver/Opserver/issues/200#issuecomment-233122437
	//  // Will return Intel[R] Ethernet Connection [2] I219-V as interface names
	//  // To merge both toghether we need to remove all with is not a-zA-Z0-9
	//	err = wmi.Query("SELECT * FROM Win32_PerfFormattedData_Tcpip_NetworkAdapter", &dstTwo)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	js, _ = json.Marshal(dstTwo)
	//	fmt.Println(string(js))

	netResults := make(map[string]*resultNet)

	// Get MTU and Duplex Mode
	var dstNetAdapter []MSFT_NetAdapter
	_ = wmi.QueryNamespace("SELECT * FROM MSFT_NetAdapter ", &dstNetAdapter, `Root\StandardCimv2`)

	for _, nic := range dst {
		// This is the same name as gopsutil and net.Interfaces() use
		// This may be has to be refactored when switching all to WMI
		name := nic.NetConnectionID

		if name != "" {
			duplex := DUPLEX_UNKNOWN
			var mtu int64 = 0

			for _, mfstNic := range dstNetAdapter {
				if mfstNic.Name == name {
					if mfstNic.FullDuplex {
						duplex = DUPLEX_FULL
					} else {
						duplex = DUPLEX_HALF
					}

					mtu = int64(mfstNic.MtuSize)

				}
			}

			netResults[name] = &resultNet{
				Isup:   nic.NetConnectionStatus == 2,
				MTU:    mtu,
				Speed:  int64(nic.Speed) / 1000 / 1000, // bits/s to mbits/s
				Duplex: duplex,
			}
		}

	}

	return netResults, nil
}
