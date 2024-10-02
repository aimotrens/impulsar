package shellexecutor

import "golang.org/x/sys/windows"

// Windows Handle ReadWriteCloser
type winHandleRWC struct {
	handle windows.Handle
}

func (h *winHandleRWC) Read(p []byte) (int, error) {
	var n uint32 = 0
	err := windows.ReadFile(h.handle, p, &n, nil)
	return int(n), err
}

func (h *winHandleRWC) Write(p []byte) (int, error) {
	var n uint32 = 0
	err := windows.WriteFile(h.handle, p, &n, nil)
	return int(n), err
}

func (h *winHandleRWC) Close() error {
	return windows.CloseHandle(h.handle)
}
