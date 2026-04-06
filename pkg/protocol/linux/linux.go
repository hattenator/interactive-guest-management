package linux

import (
	"encoding/json"
	"log"
	"net"
	"os/exec"
	"regexp"
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
					handleCommand(decodedMessage, socketHolder)
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

func handleCommand(decodedMessage protocol.CmdMessage, socketHolder SocketListener) {

	switch decodedMessage.Command {
	case "GetPowerState":
		log.Printf("Handling GetPowerState")
		GetPowerState(decodedMessage, socketHolder)
	default:
		return
	}

}

func GetPowerState(decodedMessage protocol.CmdMessage, socketHolder SocketListener) {
	//TODO this should hand off to a separate process
	tunedCmd := exec.Command("tuned-adm", "profile")
	//var outputBuf = bytes.Buffer
	//tunedCmd.Stdout = &outputBuf

	output, outputErr := tunedCmd.Output()
	if outputErr == nil {

		outputStr := string(output)
		// Use a regular expression to capture the active profile name
		re := regexp.MustCompile(`(?m)^Current active profile: ([^\s]+)$`)
		matches := re.FindStringSubmatch(outputStr)
		if len(matches) > 1 {
			profile := matches[1]
			log.Printf("Current active profile: %s", profile)
			RespondPowerState(profile, decodedMessage, socketHolder)
		} else {
			log.Printf("Could not parse active profile from output: %s", outputStr)
		}

	} else {
		log.Printf("tuned-adm output error: %v", outputErr)
	}

}

func RespondPowerState(profile string, decodedMessage protocol.CmdMessage, socketHolder SocketListener) {

	response, messageError := protocol.NewCmdMessage(profile, decodedMessage.Nonce, []byte("TODO Salt"))
	if messageError != nil {
		log.Printf("Message Creation error: %v", messageError)
	}
	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		panic(err)
	}
    retryAttempts := 0
    maxRetry := 3
    bytesWritten := 0
    for retryAttempts < maxRetry {
        bw, err := socketHolder.GuestSocket.Write(jsonBytes)
        if err == nil && bw > 0 {
            bytesWritten = bw
            break
        }
        retryAttempts++
        log.Printf("Error sending PowerState Response attempt %d/%d: %v", retryAttempts, maxRetry, err)
        time.Sleep(500 * time.Millisecond)
    }
    if bytesWritten < 1 {
        log.Printf("Incomplete response send after %d attempts", retryAttempts)
        return
    }
    log.Printf("Sent Power State: %s", jsonBytes)
}
