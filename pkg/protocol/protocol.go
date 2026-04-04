package protocol

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// CmdMessage holds the command payload and related metadata.
// It is marshalled to JSON for transmission.
// The `Signature` field is the 32‑byte digest produced by a fast hash
// algorithm (SHA‑256 in this example – replace with any other fast
// algorithm such as SHA‑512 or xxhash if required).
//
// NOTE: Because a byte array is not a JSON string, it is encoded as a
// JSON array of numbers (i.e. []uint8). If you prefer a hex string you
// can use the [encoding/hex] package and store the signature as a
// string instead.
type CmdMessage struct {
	Command   string `json:"command"`
	Nonce     uint64 `json:"nonce"`
	Timestamp uint64 `json:"timestamp"`
	Signature string `json:"signature"`
}

// NewCmdMessage creates a CmdMessage with the supplied command string.
// It automatically fills in the nonce and timestamp, and generates a
// digital signature using the supplied key.
func NewCmdMessage(cmd string, nonce uint64, key []byte) (*CmdMessage, error) {
	m := &CmdMessage{
		Command:   cmd,
		Nonce:     nonce,
		Timestamp: uint64(time.Now().Unix()),
	}
	// Use HMAC‑SHA256 to generate the signature with the provided key.
	h := hmac.New(sha256.New, key)
	h.Write([]byte(cmd))
	h.Write(uint64ToBytes(nonce))
	h.Write(uint64ToBytes(m.Timestamp))
	sigBytes := h.Sum(nil)
	// Convert the 32‑byte digest to a Base64 string.
	m.Signature = base64.StdEncoding.EncodeToString(sigBytes)
	return m, nil
}

// uint64ToBytes converts a uint64 to an 8 byte slice in big‑endian order.
func uint64ToBytes(v uint64) []byte {
	b := make([]byte, 8)
	for i := uint(0); i < 8; i++ {
		b[7-i] = byte(v >> (i * 8))
	}
	return b
}

func Example() {
	// Example usage
	nonce := uint64(42)
	msg, err := NewCmdMessage("do_something", nonce, []byte("secret"))
	if err != nil {
		panic(err)
	}
	jsonBytes, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
}
