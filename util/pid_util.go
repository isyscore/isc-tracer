package util

import "os"

var pid uint16

func GetPid() uint16 {
	if pid != 0 {
		return pid
	}
	pid = uint16(os.Getpid())
	return pid
}
