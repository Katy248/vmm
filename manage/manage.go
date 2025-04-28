package manage

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"vmm/manage/qmp"
	"vmm/vm"
)

func Init(vm vm.VM, diskSize int) error {
	cmd := exec.Command(
		"qemu-img", "create",
		"-f", string(vm.ImageType),
		vm.GetImageFile(),
		fmt.Sprintf("%dG", diskSize))
	{
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
func Start(vm vm.VM, ram int, isoFile string) error {
	cmd := exec.Command(
		"qemu-system-x86_64",
		"-drive", fmt.Sprintf("file=%s,format=%s", vm.GetImageFile(), vm.ImageType),
		"-m", fmt.Sprintf("%dG", ram),
		"-pidfile", vm.GetProcessIdFile(),
		"-cdrom", isoFile,
		"-enable-kvm",
		"-cpu", "host",
		"-qmp", fmt.Sprintf("unix:%s,server,nowait", vm.GetSocketFile()),
		"-machine", "q35",
		"-daemonize",
	)
	{
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}

func GetStatus(vm vm.VM) (*qmp.VMStatus, error) {
	connection, err := qmp.New(vm.GetSocketFile())
	if err != nil {
		return nil, err
	}
	defer connection.Close()
	connection.SendQmpCapabilities()

	status, err := connection.GetVMStatus(vm)
	if err != nil {
		return nil, err
	}
	log.Printf("VM status: %s", status.Status)

	return status, nil
}

func Stop(vm vm.VM) error {
	conn, err := qmp.New(vm.GetSocketFile())
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SendQmpCapabilities()

	return conn.SendStop()
}
func Delete(vm vm.VM) error {
	log.Printf("Deleting VM %s", vm.Name)
	return errors.Join(
		os.Remove(vm.GetProcessIdFile()),
		os.Remove(vm.GetImageFile()),
		os.Remove(vm.GetSocketFile()),
	)
}
