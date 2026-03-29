//go:build windows
// +build windows

package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/eventlog"
)

const (
	sourceName   = "QEMUHostAgent" // name used in Event Viewer → Applications and Services Logs
	registerKind = eventlog.Info | eventlog.Warning | eventlog.Error
)

var (
	elog *eventlog.Log
)

// init registers the app as an event‑log source and opens the log handle.
func init() {
	var err error

	// Register the source once (only the first run will create a registry entry).
	// Subsequent runs will fail with “The source is already registered”.
	if err = eventlog.Install(sourceName, sourceName, false, registerKind); err != nil {
		// Ignore if already registered – this keeps the program idempotent.
		if !os.IsExist(err) {
			log.Printf("Failed to register event source: %v", err)
		}
	}

	elog, err = eventlog.Open(sourceName)
	if err != nil {
		log.Fatalf("Failed to open event log: %v", err)
	}

	// Hook the standard logger to also write to the event log.
	log.SetOutput(os.Stdout)
}

// closeCleanUp closes the event log when the program exits.
func closeCleanUp() {
	if elog != nil {
		elog.Close()
	}
}

// logToEvent writes the message to the Windows Event Log and to stdout.
func logToEvent(severity string, msg string) {
	switch severity {
	case "info":
		if err := elog.Info(1, msg); err != nil {
			log.Printf("EventLog Info error: %v", err)
		}
	case "warning":
		if err := elog.Warning(2, msg); err != nil {
			log.Printf("EventLog Warning error: %v", err)
		}
	case "error":
		if err := elog.Error(3, msg); err != nil {
			log.Printf("EventLog Error error: %v", err)
		}
	default:
		// unknown severity – fall back to info
		if err := elog.Info(1, msg); err != nil {
			log.Printf("EventLog Info error: %v", err)
		}
	}
	// Also print to the console via the standard logger
	log.Printf("[%s] %s", severity, msg)
}

func main() {
	defer closeCleanUp()
	path := `\\.\Global\host.0`
	hostSocket := openSocket(path)
	defer windows.CloseHandle(hostSocket)

}

// GetLastInputInfo returns the tick count value of the last user input.
func getLastInputTime() (uint32, error) {
	const user32 = "user32.dll"

	var (
		//ubi      = syscall.NewCallback(nil) // unused; just to keep Go happy
		luidInfo = struct {
			cbSize uint32
			dwTime uint32
		}{cbSize: 8} // size of the struct
	)

	getLastInputInfo := syscall.NewLazyDLL(user32).NewProc("GetLastInputInfo")

	// Call: BOOL GetLastInputInfo(LPINPUTINFO lpInputInfo);
	r, _, err := getLastInputInfo.Call(uintptr(unsafe.Pointer(&luidInfo)))
	if r == 0 {
		return 0, err
	}
	return luidInfo.dwTime, nil
}

// detectIdle returns true if the last user input was more than 15 minutes ago.
func detectIdle() bool {
	const idleThreshold = 15 * time.Minute // 15 minutes

	idleDuration, b, shouldReturn := getIdleDuration()
	if shouldReturn {
		return b
	}
	return (idleDuration >= idleThreshold)
}

func getIdleDuration() (time.Duration, bool, bool) {
	lastInputTicks, err := getLastInputTime()
	if err != nil {
		// In a real application you might want to log this
		return 0, false, true
	}

	currentTicks := getCurrentTimeTicks()
	if currentTicks < uint64(lastInputTicks) {
		// sanity check – shouldn't happen, but guard it
		return 0, false, true
	}

	idleDuration := time.Duration(currentTicks-uint64(lastInputTicks)) * time.Millisecond
	return idleDuration, false, false
}

func getCurrentTimeTicks() uint64 {
	// GetTickCount64 returns the number of milliseconds that have elapsed
	// since the system was started.
	getTickCount64 := syscall.NewLazyDLL("kernel32.dll").NewProc("GetTickCount64")
	currentTicksUncast, _, _ := getTickCount64.Call()
	currentTicks := uint64(currentTicksUncast)
	return currentTicks
}

func getPowerState(hostSocket windows.Handle) (powerState string) {
	// --- write ----------------------------------------------------------
	var written uint32

	if err := windows.WriteFile(hostSocket, []byte("hello\n"), &written, nil); err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes\n", written)

	// --- read ------------------------------------------------------------
	buf := make([]byte, 1024)
	var n uint32
	if err := windows.ReadFile(hostSocket, buf, &n, nil); err != nil {
		panic(err)
	}
	fmt.Printf("read %d bytes: %q\n", n, buf[:n])
	return
}

func openSocket(path string) (hostSocket windows.Handle) {
	p, err := windows.UTF16PtrFromString(path)
	if err != nil {
		panic(err)
	}

	socketOpen := false
	for !socketOpen {
		hostSocket, err = windows.CreateFile(
			p,
			windows.GENERIC_READ|windows.GENERIC_WRITE,       // desired access
			windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, // sharing
			nil, // default security
			windows.OPEN_EXISTING,
			0,
			0, // no template
		)
		if err != nil {
			errorToLog := fmt.Sprintf("CreateFile failed: %v\n", err)
			if errno, ok := err.(windows.Errno); ok {
				errorToLog += (fmt.Sprintf("Win32 error: %d\n", errno))
			}
			logToEvent("error", errorToLog)
			time.Sleep(60000 * time.Millisecond)
		}
		socketOpen = true
	}
	fmt.Println("opened", path, hostSocket)
	return hostSocket
}
