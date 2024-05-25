module github.com/it-novum/openitcockpit-agent-go

go 1.15

require (
	github.com/andybalholm/crlf v0.0.0-20171020200849-670099aa064f
	github.com/containerd/containerd v1.6.26 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2
	github.com/distatus/battery v0.10.0
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.14+incompatible
	github.com/fsnotify/fsnotify v1.5.3 // indirect
	github.com/godbus/dbus/v5 v5.1.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/leoluk/perflib_exporter v0.1.0
	github.com/lufia/plan9stats v0.0.0-20220326011226-f1430873d8db // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/power-devops/perfstat v0.0.0-20220216144756-c35f1ee13d7c // indirect
	github.com/prometheus/procfs v0.7.3
	github.com/shirou/gopsutil/v3 v3.22.3
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.11.0
	github.com/yusufpapurcu/wmi v1.2.2
	golang.org/x/sys v0.13.0
	golang.org/x/text v0.13.0
	howett.net/plist v1.0.0 // indirect
	libvirt.org/libvirt-go v7.4.0+incompatible
)

replace github.com/shirou/gopsutil/v3 v3.20.12 => github.com/it-novum/gopsutil/v3 v3.21.2-0.20210201093253-6e7f4ffe9947
