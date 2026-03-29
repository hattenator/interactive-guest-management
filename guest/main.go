package main

import (
	"fmt"

	"golang.org/x/sys/windows"
)

func main() {
	path := `\\.\Global\host.0`
	err, hostSocket := openSocket(path)
	defer windows.CloseHandle(hostSocket)

	// --- write ----------------------------------------------------------
	var written uint32
	if err = windows.WriteFile(hostSocket, []byte("hello\n"), &written, nil); err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes\n", written)

	// --- read ------------------------------------------------------------
	buf := make([]byte, 1024)
	var n uint32
	if err = windows.ReadFile(hostSocket, buf, &n, nil); err != nil {
		panic(err)
	}
	fmt.Printf("read %d bytes: %q\n", n, buf[:n])
}

func openSocket(path string) (error, windows.Handle) {
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		panic(err)
	}

	h, err := windows.CreateFile(
		p,
		windows.GENERIC_READ|windows.GENERIC_WRITE,       // desired access
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, // sharing
		nil, // default security
		windows.OPEN_EXISTING,
		0,
		0, // no template
	)
	if err != nil {
		fmt.Printf("CreateFile failed: %v\n", err)
		if errno, ok := err.(windows.Errno); ok {
			fmt.Printf("Win32 error: %d\n", errno)
		}
		panic(err)
	}

	fmt.Println("opened", path, h)
	return err, h
}
