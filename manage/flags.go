package manage

import "fmt"

// In GB
func diskSizeFlag(diskSize int) string {
	return fmt.Sprintf("%dG", diskSize)
}
