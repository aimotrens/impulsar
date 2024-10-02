package shellexecutor

import (
	"golang.org/x/term"
	"io"
	"os"
	"os/exec"
)

func runCmd(cmd *exec.Cmd, interactive bool) error {
	if interactive && isConPtyAvailable && term.IsTerminal(int(os.Stdout.Fd())) {
		err, isPtyErr := runInPty(cmd)

		if !isPtyErr {
			return err
		} else {
			//fmt.Println("Failed to run command in PTY")
		}
	}

	return cmd.Run()
}

func runInPty(cmd *exec.Cmd) (error, bool) {
	pty, err := NewPty()
	if err != nil {
		return err, true
	}
	defer pty.Close()

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return err, true
	}

	if err = pty.SetDimensions(width, height); err != nil {
		return err, true
	}

	go io.Copy(cmd.Stdout, pty)
	go io.Copy(pty, os.Stdin)

	return pty.RunCmd(cmd)
}
