package vm

type VM struct {
	Name      string
	ImageType ImageType
}

func (vm *VM) GetImageFile() string {
	return vm.Name + "-" + string(vm.ImageType) + ".img"
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
