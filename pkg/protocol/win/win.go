package win

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/hattenator/interactive-guest-management/pkg/protocol"
	"golang.org/x/sys/windows"
)

type SocketListener struct {
	HostSocket windows.Handle
}

func (w SocketListener) SendCommand(command string) {
	msg, err := protocol.NewCmdMessage(command, rand.Uint64(), []byte("secret"))
	if err != nil {
		panic(err)
	}
	jsonBytes, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		panic(err)
	}

	// --- write ----------------------------------------------------------
	var written uint32

	if err := windows.WriteFile(w.HostSocket, jsonBytes, &written, nil); err != nil {
		panic(err)
	}
	fmt.Printf("wrote %d bytes\n", written)

	//TODO put a time limit on the Read and retry
	//TODO catch socket errors and reopen the socket

	// --- read ------------------------------------------------------------
	buf := make([]byte, 1024)
	var n uint32
	if err := windows.ReadFile(w.HostSocket, buf, &n, nil); err != nil {
		panic(err)
	}
	fmt.Printf("read %d bytes: %q\n", n, buf[:n])
	return
}
