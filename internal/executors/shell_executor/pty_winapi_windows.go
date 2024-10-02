package shellexecutor

import (
	"golang.org/x/sys/windows"
	"unicode/utf16"
	"unsafe"
)

var (
	isConPtyAvailable = true
)

func init() {
	kernel32 := windows.NewLazySystemDLL("kernel32.dll")
	necessaryFunctions := []string{
		"CreatePseudoConsole",
		"ClosePseudoConsole",
		"InitializeProcThreadAttributeList",
		"UpdateProcThreadAttribute",
	}

	for _, f := range necessaryFunctions {
		if kernel32.NewProc(f).Find() != nil {
			isConPtyAvailable = false
		}
	}
}

func createProcess(hpc windows.Handle, commandLine, workDir string, env []string) (*windows.ProcessInformation, error) {
	cmdLine, err := windows.UTF16PtrFromString(commandLine)
	if err != nil {
		return nil, err
	}

	var currentDirectory *uint16
	if workDir != "" {
		if currentDirectory, err = windows.UTF16PtrFromString(workDir); err != nil {
			return nil, err
		}
	}

	siEx, err := createStartupInfoEx(hpc)
	if err != nil {
		return nil, err
	}

	flags := uint32(windows.EXTENDED_STARTUPINFO_PRESENT)

	var envBlock *uint16
	if env != nil {
		flags |= uint32(windows.CREATE_UNICODE_ENVIRONMENT)
		envBlock = createEnvBlock(env)
	}

	var pi windows.ProcessInformation
	err = windows.CreateProcess(
		nil, // use this if no args
		cmdLine,
		nil,
		nil,
		false, // inheritHandle
		flags,
		envBlock,
		currentDirectory,
		&siEx.StartupInfo,
		&pi,
	)

	if err != nil {
		return nil, err
	}

	return &pi, nil
}

// createEnvBlock refers to syscall.createEnvBlock in go/src/syscall/exec_windows.go
// Sourced From: https://github.com/creack/pty/pull/155
func createEnvBlock(envv []string) *uint16 {
	if len(envv) == 0 {
		return &utf16.Encode([]rune("\x00\x00"))[0]
	}
	length := 0
	for _, s := range envv {
		length += len(s) + 1
	}
	length += 1

	b := make([]byte, length)
	i := 0
	for _, s := range envv {
		l := len(s)
		copy(b[i:i+l], []byte(s))
		copy(b[i+l:i+l+1], []byte{0})
		i = i + l + 1
	}
	copy(b[i:i+1], []byte{0})

	return &utf16.Encode([]rune(string(b)))[0]
}

func createStartupInfoEx(hpc windows.Handle) (*windows.StartupInfoEx, error) {
	attrList, err := windows.NewProcThreadAttributeList(1)
	if err != nil {
		return nil, err
	}

	err = attrList.Update(windows.PROC_THREAD_ATTRIBUTE_PSEUDOCONSOLE, unsafe.Pointer(hpc), unsafe.Sizeof(hpc))
	if err != nil {
		return nil, err
	}

	siEx := windows.StartupInfoEx{
		StartupInfo: windows.StartupInfo{
			Cb:    uint32(unsafe.Sizeof(windows.StartupInfoEx{})),
			Flags: windows.STARTF_USESTDHANDLES,
		},
		ProcThreadAttributeList: attrList.List(),
	}

	return &siEx, nil
}
