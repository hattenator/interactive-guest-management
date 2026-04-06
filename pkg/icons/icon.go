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
// Windows‑compatible ICO file.
func MakeIcon(c color.RGBA) []byte {
	// Create the 32x32 image
	img := image.NewRGBA(image.Rect(0, 0, 32, 32))
	for y := 0; y < 32; y++ {
		for x := 0; x < 32; x++ {
			img.Set(x, y, c)
		}
	}

	// Encode as PNG
	pngBuf := &bytes.Buffer{}
	if err := png.Encode(pngBuf, img); err != nil {
		return nil
	}
	pngData := pngBuf.Bytes()

	// ICONDIR header
	iconDir := make([]byte, iconDirSize)
	binary.LittleEndian.PutUint16(iconDir[0:], 0)   // Reserved
	binary.LittleEndian.PutUint16(iconDir[2:], 1)   // Type: 1 = icon
	binary.LittleEndian.PutUint16(iconDir[4:], 1)   // Count: 1 icon

	// Icon directory entry (ICONDIRENTRY)
	entry := make([]byte, iconEntrySize)
	entry[0] = 32                        // Width (32)
	entry[1] = 32                        // Height (32)
	entry[2] = 0                         // Color count (0 for PNG)
	entry[3] = 0                         // Reserved
	entry[4] = 0                         // Planes (0 for PNG)
	entry[5] = 0                         // Bits per pixel (0 for PNG)
	binary.LittleEndian.PutUint32(entry[8:], uint32(len(pngData))) // Bytes in resource
	binary.LittleEndian.PutUint32(entry[12:], uint32(iconDirSize+iconEntrySize)) // Offset

	// Assemble the icon file
	ico := &bytes.Buffer{}
	ico.Write(iconDir)
	ico.Write(entry)
	ico.Write(pngData)

	return ico.Bytes()
}

// Power icon data used by the systray.
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
