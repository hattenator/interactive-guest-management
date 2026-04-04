package linux

import (
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/hattenator/interactive-guest-management/pkg/protocol"
)

type SocketListener struct {
	GuestSocket net.Conn
}

func (socketHolder SocketListener) ReceiveCommands() {
	for true {
		// TODO for security, I should limit the size of messages received.
		var clientMessage []byte
		readSize, err := socketHolder.GuestSocket.Read(clientMessage)
		if err == nil && readSize > 0 {
			decodedMessage := protocol.CmdMessage{}
			jsonErr := json.Unmarshal(clientMessage, &decodedMessage)
			if jsonErr == nil {
				log.Printf("Received Command: %v", decodedMessage)
			}
		}

		time.Sleep(1000 * time.Millisecond)
	}

}
