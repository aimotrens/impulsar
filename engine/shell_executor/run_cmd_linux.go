package shellexecutor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/term"
)

func runCmd(cmd *exec.Cmd) error {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		//fmt.Println("Running command in PTY...")
		err, isPtyErr := runInPty(cmd)

		if !isPtyErr {
			return err
		} else {
			//fmt.Println("Failed to run command in PTY")
		}
	}

	//fmt.Println("Running command in non-PTY...")
	return cmd.Run()
}

func runInPty(cmd *exec.Cmd) (error, bool) {
	// open PTY
	master, slave, err := openPty()
	if err != nil {
		//fmt.Printf("Failed to open PTY: %v\n", err)
		return err, true
	}
	defer master.Close()
	defer slave.Close()

	// set PTY size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		//fmt.Printf("Failed to get terminal size: %v\n", err)
		return err, true
	}

	err = setPtySize(master, uint16(height), uint16(width))
	if err != nil {
		//fmt.Printf("Failed to set PTY size: %v\n", err)
		return err, true
	}

	cmd.Stdout = slave
	cmd.Stderr = slave
	cmd.Stdin = slave

	// run command
	err = cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start command: %v\n", err)
		return err, false
	}

	// link output and input from PTY to stdout
	go io.Copy(os.Stdout, master)
	go io.Copy(master, os.Stdin)

	// wait for command to finish
	err = cmd.Wait()

	return err, false
}

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
