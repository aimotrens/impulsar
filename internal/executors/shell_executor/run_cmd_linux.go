package shellexecutor

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"golang.org/x/term"
)

func runCmd(cmd *exec.Cmd, interactive bool) error {
	if interactive && term.IsTerminal(int(os.Stdout.Fd())) {
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

	// link output and input
	go io.Copy(os.Stdout, master)
	go io.Copy(master, os.Stdin)

	// wait for command to finish
	err = cmd.Wait()

	return err, false
}
