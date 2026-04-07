package icons

import (
	"bytes"
	"encoding/binary"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

// TestMakeIconStructure verifies that MakeIcon produces a valid ICO file
// containing a single 32×32 PNG image.
func TestMakeIconStructure(t *testing.T) {
	t.Skip("this test is obsolete")

	generated, err := MakeIcon(16, 0x00, 0xFF, 0x00, 0xFF)
	if err != nil {
		t.Fatalf("Error generating icon: %v", err)
	}
	if len(generated) < 22 { // at least header+one entry
		t.Fatalf("generated icon is too short: %d bytes", len(generated))
	}
	// Validate header
	if generated[0] != 0 || generated[1] != 0 {
		t.Fatalf("reserved bytes not zero: %v", generated[:2])
	}
	if generated[2] != 1 || generated[3] != 0 {
		t.Fatalf("type not 1, got %d", generated[2])
	}
	count := binary.LittleEndian.Uint16(generated[4:6])
	if count == 0 || count > 2 {
		t.Fatalf("unexpected icon count %d", count)
	}
	// Helper to decode PNG from data slice
	decodePNG := func(data []byte) error {
		_, err := png.Decode(bytes.NewReader(data))
		return err
	}
	offset := uint32(iconDirSize + iconEntrySize*uint32(count))
	// Iterate entries
	for i := 0; i < int(count); i++ {
		eIdx := 6 + iconEntrySize*i
		width := int(generated[eIdx])
		height := int(generated[eIdx+1])
		if width == 0 {
			width = 256
		}
		if height == 0 {
			height = 256
		}
		if width != 32 && width != 256 {
			t.Fatalf("unexpected entry width %d", width)
		}
		if height != 32 && height != 256 {
			t.Fatalf("unexpected entry height %d", height)
		}
		// bytes in resource
		br := binary.LittleEndian.Uint32(generated[eIdx+8 : eIdx+12])
		// start of image data
		start := offset
		if start+uint32(br) > uint32(len(generated)) {
			t.Fatalf("entry data overflow")
		}
		data := generated[start : start+uint32(br)]
		if err := decodePNG(data); err != nil {
			t.Fatalf("PNG decode error: %v", err)
		}
		offset += uint32(br)
	}
}

func TestExportIconsToFiles(t *testing.T) {
	t.Parallel()

	outDir := t.TempDir()

	tests := []struct {
		name string
		data []byte
	}{
		{name: "red.ico", data: RedSquare16},
		{name: "yellow.ico", data: YellowSquare16},
		{name: "green.ico", data: GreenSquare16},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(outDir, tc.name)

			if len(tc.data) < 22 {
				t.Fatalf("icon too short: %d bytes", len(tc.data))
			}

			// Basic ICO header checks.
			if got := binary.LittleEndian.Uint16(tc.data[0:2]); got != 0 {
				t.Fatalf("bad reserved field: got %d want 0", got)
			}
			if got := binary.LittleEndian.Uint16(tc.data[2:4]); got != 1 {
				t.Fatalf("bad type field: got %d want 1", got)
			}
			if got := binary.LittleEndian.Uint16(tc.data[4:6]); got != 1 {
				t.Fatalf("bad image count: got %d want 1", got)
			}

			if err := os.WriteFile(path, tc.data, 0o644); err != nil {
				t.Fatalf("write %s: %v", path, err)
			}

			info, err := os.Stat(path)
			if err != nil {
				t.Fatalf("stat %s: %v", path, err)
			}
			if info.Size() == 0 {
				t.Fatalf("%s is empty", path)
			}
		})
	}
}

func TestSolidSquareICO_CustomSize(t *testing.T) {
	t.Parallel()

	icon, err := MakeIcon(32, 0x20, 0x80, 0xFF, 0xFF)
	if err != nil {
		t.Fatalf("SolidSquareICO failed: %v", err)
	}

	if len(icon) < 22 {
		t.Fatalf("icon too short: %d bytes", len(icon))
	}

	// Width/height byte in ICONDIRENTRY.
	if got := icon[6]; got != 32 {
		t.Fatalf("bad width byte: got %d want 32", got)
	}
	if got := icon[7]; got != 32 {
		t.Fatalf("bad height byte: got %d want 32", got)
	}
}

func TestExportIcon(t *testing.T) {

	icon, err := MakeIcon(16, 0xAA, 0xBB, 0xCC, 0xFF)
	if len(icon) == 0 {
		t.Fatal("MakeIcon returned empty icon data")
	}

	f, err := os.Create("./green.ico")
	if err != nil {
		t.Fatalf("Exporting Icon failed to open file: %v", err)
	}
	_, err = f.Write(icon)
	if err != nil {
		f.Close()
		t.Fatalf("Exporting Icon failed to write data: %v", err)
	}
	if err = f.Close(); err != nil {
		t.Fatalf("Failed to close file: %v", err)
	}
}
