package cdrom

import (
	"golang.org/x/sys/unix"
	"syscall"
)

type CD struct {
	fileDescriptor int
}

func New() (*CD, error) {
	fd, err := syscall.Open("/dev/cdrom", syscall.O_RDONLY|syscall.O_NONBLOCK, 0644)
	if err != nil {
		return nil, err
	}
	return &CD{fileDescriptor: fd}, nil
}

func (cd *CD) Eject() {
	unix.IoctlGetInt(cd.fileDescriptor, CDROMEJECT)
	unix.IoctlGetInt(cd.fileDescriptor, CDROMEJECT_SW)
}

func (cd *CD) Close() error {
	return syscall.Close(cd.fileDescriptor)
}
