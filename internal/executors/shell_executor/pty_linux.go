package shellexecutor

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

func openPty() (master *os.File, slave *os.File, err error) {
	// open PTY master
	master, err = os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		return nil, nil, err
	}
	masterFD := master.Fd()

	// unlock PTY
	var unused int32
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, masterFD, syscall.TIOCSPTLCK, uintptr(unsafe.Pointer(&unused)))
	if errno != 0 {
		master.Close()
		return nil, nil, errno
	}

	// retrieve PTY slave number
	var slaveNumber uint32
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, masterFD, uintptr(syscall.TIOCGPTN), uintptr(unsafe.Pointer(&slaveNumber)))
	if errno != 0 {
		master.Close()
		return nil, nil, errno
	}

	// open PTY slave
	slaveFile := fmt.Sprintf("/dev/pts/%d", slaveNumber)
	slave, err = os.OpenFile(slaveFile, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		master.Close()
		return nil, nil, err
	}

	return master, slave, nil
}

func setPtySize(fd *os.File, rows, cols uint16) error {
	// window size
	ws := &struct {
		Rows uint16
		Cols uint16
		X    uint16 // unused
		Y    uint16 // unused
	}{Rows: rows, Cols: cols, X: 0, Y: 0}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd.Fd(), syscall.TIOCSWINSZ, uintptr(unsafe.Pointer(ws)))
	if errno != 0 {
		return errno
	}
	return nil
}
