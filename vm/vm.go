package vm

import (
	"log"
	"os"
	"strings"
)

type VM struct {
	Name      string
	ImageType ImageType
	// DiskSize  int
	// uuid      uuid.UUID
}

func (vm *VM) GetImageFile() string {
	return vm.Name + "-" + string(vm.ImageType) + ".img"
}
func (vm *VM) GetProcessIdFile() string {
	return vm.Name + "-pid"
}
func (vm *VM) GetSocketFile() string {
	return vm.Name + "-socket"
}

type ImageType string

const (
	ImageTypeQCOW2 ImageType = "qcow2"
	ImageTypeRAW   ImageType = "raw"
)

func New(name string) VM {
	return VM{
		Name:      name,
		ImageType: ImageTypeRAW,
	}
}

func (vm *VM) GetPid() string {
	data, err := os.ReadFile(vm.GetProcessIdFile())
	if err != nil {
		log.Printf("Error reading process ID file '%s': %v", vm.GetProcessIdFile(), err)
		return ""
	}
	return strings.TrimSpace(string(data))
}
