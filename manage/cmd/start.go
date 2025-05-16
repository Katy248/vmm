package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"vmm/vm"
)

type StartData struct {
	ImageFile  string
	ImageType  vm.ImageType
	RamInGb    int
	IsoFile    string
	SocketFile string
}

func ExecStartDaemonize(data StartData) error {
	command := exec.Command(
		"qemu-system-x86_64",
		"-drive", fmt.Sprintf("file=%s,format=%s", data.ImageFile, data.ImageType),
		"-m", fmt.Sprintf("%dG", data.RamInGb),
		"-cdrom", data.IsoFile,
		"-enable-kvm",
		"-cpu", "host",
		"-qmp", fmt.Sprintf("unix:%s,server,nowait", data.SocketFile),
		"-machine", "q35",
		"-daemonize",
	)
	{
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
	}
	return command.Run()

}
