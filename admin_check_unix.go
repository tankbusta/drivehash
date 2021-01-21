// +build linux darwin

package drivehash

import "os"

func checkIfAdmin() bool {
	return os.Getuid() == 0
}
