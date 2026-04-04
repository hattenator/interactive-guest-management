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
		const MAX_MESSAGE_SIZE = 8192
		var clientMessage = make([]byte, MAX_MESSAGE_SIZE)
		readSize, err := socketHolder.GuestSocket.Read(clientMessage)
		if err == nil {
			if readSize > 0 {
				for readSize == MAX_MESSAGE_SIZE && err == nil {
					log.Printf("Message Exceeded Maximum Size, Draining")
					// If you see this, we could increase the MAX_MESSAGE_SIZE
					readSize, err = socketHolder.GuestSocket.Read(clientMessage)
				}
				decodedMessage := protocol.CmdMessage{}
				jsonErr := json.Unmarshal(clientMessage[:readSize], &decodedMessage)
				if jsonErr == nil {
					log.Printf("Received Command: %v", decodedMessage)
				} else {
					log.Printf("Json error when receiving from socket: %v", err)
					log.Printf("Buffer was:\n%s", clientMessage[:readSize])
				}

			}
		} else {
			log.Printf("Error receiving commands from socket: %v", err)
		}
		time.Sleep(1000 * time.Millisecond)
	}

}
