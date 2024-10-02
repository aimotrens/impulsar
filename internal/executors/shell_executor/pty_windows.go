package shellexecutor

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os/exec"
	"strings"
)

type Pty struct {
	hPc                          windows.Handle
	processInfo                  *windows.ProcessInformation
	cmdIn, cmdOut, ptyIn, ptyOut windows.Handle
	width, height                int
}

func NewPty() (*Pty, error) {
	if !isConPtyAvailable {
		return nil, fmt.Errorf("ConPTY is not available")
	}

	p := Pty{
		width:  80,
		height: 40,
	}

	if err := windows.CreatePipe(&p.ptyIn, &p.cmdIn, nil, 0); err != nil {
		return nil, fmt.Errorf("CreatePipe: %v", err)
	}

	if err := windows.CreatePipe(&p.cmdOut, &p.ptyOut, nil, 0); err != nil {
		p.Close()
		return nil, fmt.Errorf("CreatePipe: %v", err)
	}

	return &p, nil
}

func (p *Pty) SetDimensions(width, height int) error {
	if width < 1 {
		return fmt.Errorf("width must be greater than 0")
	}
	if height < 1 {
		return fmt.Errorf("height must be greater than 0")
	}

	p.width = width
	p.height = height

	return nil
}

func (p *Pty) Write(data []byte) (int, error) {
	return (&winHandleRWC{p.cmdIn}).Write(data)
}

func (p *Pty) Read(data []byte) (int, error) {
	return (&winHandleRWC{p.cmdOut}).Read(data)
}

func (p *Pty) Close() {
	if p.hPc != 0 {
		windows.ClosePseudoConsole(p.hPc)
	}
	if p.ptyIn != 0 {
		_ = windows.CloseHandle(p.ptyIn)
	}
	if p.ptyOut != 0 {
		_ = windows.CloseHandle(p.ptyOut)
	}
	if p.cmdIn != 0 {
		_ = windows.CloseHandle(p.cmdIn)
	}
	if p.cmdOut != 0 {
		_ = windows.CloseHandle(p.cmdOut)
	}
}

func (p *Pty) RunCmd(cmd *exec.Cmd) (error, bool) {
	var err error

	coords := windows.Coord{
		X: int16(p.width),
		Y: int16(p.height),
	}

	err = windows.CreatePseudoConsole(coords, p.ptyIn, p.ptyOut, 0, &p.hPc)
	if err != nil {
		p.Close()
		return err, true
	}

	pathWithArgs := strings.Join(append([]string{cmd.Path}, encapsulate(cmd.Args[1:]...)...), " ")
	p.processInfo, err = createProcess(p.hPc, pathWithArgs, cmd.Dir, cmd.Env)
	if err != nil {
		p.Close()
		return fmt.Errorf("Failed to create console process: %v", err), true
	}

	var exitCode uint32
	_, _ = windows.WaitForSingleObject(p.processInfo.Process, windows.INFINITE)
	_ = windows.GetExitCodeProcess(p.processInfo.Process, &exitCode)
	if exitCode == 0 {
		return nil, false
	}

	return fmt.Errorf("process exited with code %d", exitCode), false
}

func encapsulate(in ...string) []string {
	for i, s := range in {
		//s = strings.ReplaceAll(s, "\"", "\"\"")
		if strings.Contains(s, " ") {
			in[i] = fmt.Sprintf("\"%s\"", s)
		}
	}
	return in
}
