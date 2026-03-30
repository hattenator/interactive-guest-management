package main

import (
	"log"
	"net"
	"time"

	"github.com/hattenator/interactive-guest-management/pkg/protocol"
)

// We need access to read a qemu socket
// And we need access to send tuned-adm commands
// I think tuned should be its own daemon.
func main() {

	//Open a the socket file
	guestSocket := openGuestSocket()
	guestProtocol := protocol.Linux{GuestSocket: guestSocket}
	go guestProtocol.ReceiveCommands()

	for true {
		time.Sleep(5000 * time.Millisecond)
	}
}

func openGuestSocket() net.Conn {

	conn, err := net.Dial("unix", "/var/lib/libvirt/qemu/win11-agent.sock")
	if err != nil {
		log.Fatalf("Failed to open guest socket: %v", err)
	}
	return conn

}
