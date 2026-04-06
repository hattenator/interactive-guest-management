package icons

import (
    "bytes"
    "encoding/binary"
    "image/color"
    "image/png"
    "testing"
    "os"
)

// TestMakeIconStructure verifies that MakeIcon produces a valid ICO file
// containing a single 32×32 PNG image.
func TestMakeIconStructure(t *testing.T) {
    generated := MakeIcon(color.RGBA{0x00, 0xFF, 0x00, 0xFF})
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

func TestIconHelpers(t *testing.T) {
    col := color.RGBA{0xAA, 0xBB, 0xCC, 0xFF}
    pngBytes := squarePNG(col)
    if len(pngBytes) == 0 {
        t.Fatalf("squarePNG returned empty data")
    }
    if _, err := png.Decode(bytes.NewReader(pngBytes)); err != nil {
        t.Fatalf("squarePNG returned invalid PNG: %v", err)
    }
    src := makeSquare(col)
    center := src.At(src.Bounds().Dx()/2, src.Bounds().Dy()/2)
    if got := color.RGBAModel.Convert(center).(color.RGBA); got != col {
        t.Fatalf("makeSquare produced wrong pixel: got %v, want %v", got, col)
    }
    resized := Resize(src, 16)
    if resized.Bounds().Dx() != 16 || resized.Bounds().Dy() != 16 {
        t.Fatalf("resized image has wrong size: %dx%d", resized.Bounds().Dx(), resized.Bounds().Dy())
    }
}




func TestExportIcon(t *testing.T) {
	c := color.RGBA{0xAA, 0xBB, 0xCC, 0xFF}
	icon := MakeIcon(c)
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

