//go:build windows
// +build windows

package main

import (
	"testing"
)

/*
func TestMakeIconStructure(t *testing.T) {
	// Recreate the icon bytes using makeIcon with same color as powerOnIcon
	generated := MakeIcon(color.RGBA{0x00, 0xFF, 0x00, 0xFF})
	if !bytes.Equal(generated, powerOnIcon) {
		t.Fatalf("generated icon bytes do not match expected powerOnIcon")
	}

	// Basic checks of the ICO structure
	if len(generated) < 6+16*3 {
		t.Fatalf("icon data too small: %d bytes", len(generated))
	}

	// Check ICONDIR header
	if generated[0] != 0 || generated[1] != 0 {
		t.Fatalf("reserved bytes not zero")
	}
	if generated[2] != 1 || generated[3] != 0 {
		t.Fatalf("type not 1")
	}
	if generated[4] != 3 || generated[5] != 0 {
		t.Fatalf("count not 3")
	}

	// Validate entries
	expectedSizes := []uint8{16, 32, 48}
	for i, sz := range expectedSizes {
		offset := 6 + i*16
		if generated[offset] != sz {
			t.Fatalf("width of entry %d wrong: %d", i, generated[offset])
		}
		if generated[offset+1] != sz {
			t.Fatalf("height of entry %d wrong: %d", i, generated[offset+1])
		}
		// Planes and bitcount should be zero for PNG
		if generated[offset+8] != 0 || generated[offset+9] != 0 {
			t.Fatalf("planes/bitcount not zero for entry %d", i)
		}
	}

	// Validate that each PNG image is a valid PNG by decoding

	for i := 0; i < 3; i++ {
		// Calculate the position of the ICONDIRENTRY
		entryOffset := 6 + i*16
		// Bytes in resource (PNG length)
		bytesInRes := binary.LittleEndian.Uint32(generated[entryOffset+12 : entryOffset+16])
		// Image offset
		imageOffset := binary.LittleEndian.Uint32(generated[entryOffset+16 : entryOffset+20])
		// Ensure bounds
		if int(imageOffset)+int(bytesInRes) > len(generated) {
			t.Fatalf("entry %d data out of bounds: off=%d len=%d", i, imageOffset, bytesInRes)
		}
		pngData := generated[imageOffset : imageOffset+bytesInRes]
		// Decode PNG
		if _, err := png.Decode(bytes.NewReader(pngData)); err != nil {
			t.Fatalf("entry %d is not a valid PNG: %v", i, err)
		}
	}
}
*/

func TestGetLastInputTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Windows API test in short mode")
	}
	t0, err := getLastInputTime()
	if err != nil {
		t.Fatalf("getLastInputTime returned error: %v", err)
	}
	if t0 == 0 {
		t.Fatalf("getLastInputTime returned zero time, expected >0")
	}
	t.Logf("Idle Start Time is %d", t0)
}

func TestGetIdleDuration(t *testing.T) {
	idleDuration, shouldReturn, failure := getIdleDuration()
	if failure {
		t.Fatalf("getIdleDuration failed")
	}
	if shouldReturn {
		t.Fatalf("getIdleDuration threw an exception")
	}
	if idleDuration <= 0 {
		t.Fatalf("getIdleDuration should have returned positive, not %v", idleDuration)
	}
	t.Logf("Idle Duration Time is %v", idleDuration.Seconds())
}
