package icons

import (
	"bytes"
	"encoding/binary"
	"image"
	"image/color"
	"image/png"
)

const (
	iconDirSize   = 6
	iconEntrySize = 16
)

// MakeIcon renders a 32×32 solid‑colour PNG and embeds it inside a
// MakeIcon renders a 2‑image Windows‑compatible ICO file.
// It creates a 32×32 and a 256×256 solid‑colour PNG and embeds both
// images into a single ICO file so Windows can use the appropriate size.
func MakeIcon(c color.RGBA) []byte {
    // Helper to create PNG data of a given size
    makePNG := func(sz int) []byte {
        img := image.NewRGBA(image.Rect(0, 0, sz, sz))
        for y := 0; y < sz; y++ {
            for x := 0; x < sz; x++ {
                img.Set(x, y, c)
            }
        }
        buf := &bytes.Buffer{}
        _ = png.Encode(buf, img)
        return buf.Bytes()
    }

    // Create 32×32 and 256×256 PNGs
    png32 := makePNG(32)
    png256 := makePNG(256)

    // Assemble ICONDIR header: 2 entries
    iconDir := make([]byte, iconDirSize)
    binary.LittleEndian.PutUint16(iconDir[0:], 0)   // Reserved
    binary.LittleEndian.PutUint16(iconDir[2:], 1)   // Type: 1 = icon
    binary.LittleEndian.PutUint16(iconDir[4:], 2)   // Count: 2 icons

    // Entry 1: 32×32 image - PNG format has specific metadata for PNG data
    entry1 := make([]byte, iconEntrySize)
    entry1[0] = 32
    entry1[1] = 32
    entry1[2] = 0          // Colors: 0 for PNG
    entry1[3] = 0          // Colors: 0 for PNG
    binary.LittleEndian.PutUint16(entry1[4:6], 1)  // Planes: 1 for PNG
    binary.LittleEndian.PutUint16(entry1[6:8], 32) // Bit Count: 32 for PNG
    binary.LittleEndian.PutUint32(entry1[8:], uint32(len(png32)))
    // Offsets are measured from start of file
    offset32 := uint32(iconDirSize + iconEntrySize*2) // header + two entries
    binary.LittleEndian.PutUint32(entry1[12:], offset32)

    // Entry 2: 256×256 image - PNG format has specific metadata for PNG data
    entry2 := make([]byte, iconEntrySize)
    entry2[0] = 0 // 256 encoded as 0 per ICO spec
    entry2[1] = 0
    entry2[2] = 0          // Colors: 0 for PNG
    entry2[3] = 0          // Colors: 0 for PNG
    binary.LittleEndian.PutUint16(entry2[4:6], 1)  // Planes: 1 for PNG
    binary.LittleEndian.PutUint16(entry2[6:8], 32) // Bit Count: 32 for PNG
    binary.LittleEndian.PutUint32(entry2[8:], uint32(len(png256)))
    offset256 := offset32 + uint32(len(png32))
    binary.LittleEndian.PutUint32(entry2[12:], offset256)

    // Assemble the final ICO
    ico := &bytes.Buffer{}
    ico.Write(iconDir)
    ico.Write(entry1)
    ico.Write(entry2)
    ico.Write(png32)
    ico.Write(png256)

    return ico.Bytes()
}

var (
	PowerOnIcon  = MakeIcon(color.RGBA{0x00, 0xFF, 0x00, 0xFF}) // green
	PowerOffIcon = MakeIcon(color.RGBA{0x80, 0x80, 0x80, 0xFF}) // grey
	IdleIcon     = MakeIcon(color.RGBA{0xFF, 0x00, 0x00, 0xFF}) // red
	DefaultIcon  = MakeIcon(color.RGBA{0x00, 0x00, 0xFF, 0xFF}) // blue
)

// Helper functions for testing
func squarePNG(col color.RGBA) []byte {
	// used in tests
	const sz = 32
	img := image.NewRGBA(image.Rect(0, 0, sz, sz))
	for y := 0; y < sz; y++ {
		for x := 0; x < sz; x++ {
			img.Set(x, y, col)
		}
	}
	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}

func makeSquare(col color.RGBA) image.Image {
	const sz = 32
	img := image.NewRGBA(image.Rect(0, 0, sz, sz))
	for y := 0; y < sz; y++ {
		for x := 0; x < sz; x++ {
			img.Set(x, y, col)
		}
	}
	return img
}

func Resize(src image.Image, dim int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, dim, dim))
	for y := 0; y < dim; y++ {
		for x := 0; x < dim; x++ {
			// nearest pixel from src
			sx := x * src.Bounds().Dx() / dim
			sy := y * src.Bounds().Dy() / dim
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}
