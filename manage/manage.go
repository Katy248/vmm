package manage

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"vmm/manage/cmd"
	"vmm/manage/qmp"
	"vmm/vm"
)

func Init(vm vm.VM, diskSize int) error {
	cmd := exec.Command(
		"qemu-img", "create",
		"-f", string(vm.ImageType),
		vm.GetImageFile(),
		diskSizeFlag(diskSize),
	)
	{
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
func Start(vm vm.VM, ram int, isoFile string) error {
	return cmd.ExecStartDaemonize(
		cmd.StartData{
			ImageFile:  vm.GetImageFile(),
			ImageType:  vm.ImageType,
			SocketFile: vm.GetSocketFile(),
			RamInGb:    ram,
			IsoFile:    isoFile,
		},
	)
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
		os.Remove(vm.GetImageFile()),
		os.Remove(vm.GetSocketFile()),
	)
}
