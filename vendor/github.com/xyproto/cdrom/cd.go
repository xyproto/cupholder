package cdrom

import (
	"golang.org/x/sys/unix"
	"syscall"
)

type CD struct {
	fileDescriptor int
}

// New opens /dev/cdrom and returns a struct with the file descriptor
func New() (*CD, error) {
	return NewFile("/dev/cdrom")
}

// NewFile opens the given device filename and returns a struct with the file descriptor
func NewFile(deviceFilename string) (*CD, error) {
	fd, err := syscall.Open(deviceFilename, syscall.O_RDONLY|syscall.O_NONBLOCK, 0644)
	if err != nil {
		return nil, err
	}
	return &CD{fileDescriptor: fd}, nil
}

// Ejects performs IOCTL calls, using CDROMEJECT and CDROMEJECT_SW on the current file descriptor
func (cd *CD) Eject() error {
	if _, err := unix.IoctlGetInt(cd.fileDescriptor, CDROMEJECT); err != nil {
		return err
	}
	if _, err := unix.IoctlGetInt(cd.fileDescriptor, CDROMEJECT_SW); err != nil {
		return err
	}
	return nil
}

// Done closes the file descriptor
func (cd *CD) Done() error {
	return syscall.Close(cd.fileDescriptor)
}
