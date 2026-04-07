package icons

import (
	"encoding/binary"
	"fmt"
)

const (
	iconDirSize   = 6
	iconEntrySize = 16
)

// Prebuilt 16x16 tray icons you can pass directly to systray.SetIcon(...).
var (
	RedSquare16    = MustSolidSquareICO(16, 0xFF, 0x00, 0x00, 0xFF)
	YellowSquare16 = MustSolidSquareICO(16, 0xFF, 0xD7, 0x00, 0xFF)
	GreenSquare16  = MustSolidSquareICO(16, 0x00, 0xC8, 0x00, 0xFF)
)

// MustSolidSquareICO is like SolidSquareICO but panics on error.
func MustSolidSquareICO(size int, r, g, b, a byte) []byte {
	ico, err := MakeIcon(size, r, g, b, a)
	if err != nil {
		panic(err)
	}
	return ico
}

// returns the bytes of a Windows .ico file containing a single
// solid-color square icon, suitable for systray.SetIcon(...) on Windows.
//
// size should typically be 16, 20, 24, 32, or 48. For systray, 16x16 or 32x32
// are usually the most useful.
func MakeIcon(size int, r, g, b, a byte) ([]byte, error) {
	if size <= 0 || size > 256 {
		return nil, fmt.Errorf("size must be in range 1..256, got %d", size)
	}

	width := size
	height := size

	// ICO image payload is a DIB/BMP-style image:
	// - BITMAPINFOHEADER
	// - XOR bitmap pixels (BGRA, bottom-up)
	// - AND mask (1bpp, padded to 32-bit rows)
	//
	// For ICO, biHeight is doubled: XOR height + AND mask height.
	const bitmapInfoHeaderSize = 40

	// 32-bit BGRA pixels
	xorRowSize := width * 4
	xorSize := xorRowSize * height

	// 1-bit AND mask, padded to 32 bits per row
	andRowSize := ((width + 31) / 32) * 4
	andSize := andRowSize * height

	imageSize := bitmapInfoHeaderSize + xorSize + andSize
	totalSize := 6 + 16 + imageSize // ICONDIR + ICONDIRENTRY + image payload

	buf := make([]byte, totalSize)

	// ---------------------------------------------------------------------
	// ICONDIR (6 bytes)
	// ---------------------------------------------------------------------
	// WORD idReserved = 0
	// WORD idType = 1 (icon)
	// WORD idCount = 1
	binary.LittleEndian.PutUint16(buf[0:2], 0)
	binary.LittleEndian.PutUint16(buf[2:4], 1)
	binary.LittleEndian.PutUint16(buf[4:6], 1)

	// ---------------------------------------------------------------------
	// ICONDIRENTRY (16 bytes)
	// ---------------------------------------------------------------------
	// BYTE  bWidth        (0 means 256)
	// BYTE  bHeight       (0 means 256)
	// BYTE  bColorCount   (0 for >= 8bpp)
	// BYTE  bReserved     = 0
	// WORD  wPlanes       = 1
	// WORD  wBitCount     = 32
	// DWORD dwBytesInRes
	// DWORD dwImageOffset
	entry := buf[6:22]
	entry[0] = dimByte(width)
	entry[1] = dimByte(height)
	entry[2] = 0
	entry[3] = 0
	binary.LittleEndian.PutUint16(entry[4:6], 1)
	binary.LittleEndian.PutUint16(entry[6:8], 32)
	binary.LittleEndian.PutUint32(entry[8:12], uint32(imageSize))
	binary.LittleEndian.PutUint32(entry[12:16], uint32(6+16))

	// ---------------------------------------------------------------------
	// BITMAPINFOHEADER (40 bytes)
	// ---------------------------------------------------------------------
	img := buf[22:]
	binary.LittleEndian.PutUint32(img[0:4], bitmapInfoHeaderSize)      // biSize
	binary.LittleEndian.PutUint32(img[4:8], uint32(width))             // biWidth
	binary.LittleEndian.PutUint32(img[8:12], uint32(height*2))         // biHeight (XOR + AND)
	binary.LittleEndian.PutUint16(img[12:14], 1)                       // biPlanes
	binary.LittleEndian.PutUint16(img[14:16], 32)                      // biBitCount
	binary.LittleEndian.PutUint32(img[16:20], 0)                       // biCompression = BI_RGB
	binary.LittleEndian.PutUint32(img[20:24], uint32(xorSize+andSize)) // biSizeImage
	binary.LittleEndian.PutUint32(img[24:28], 0)                       // biXPelsPerMeter
	binary.LittleEndian.PutUint32(img[28:32], 0)                       // biYPelsPerMeter
	binary.LittleEndian.PutUint32(img[32:36], 0)                       // biClrUsed
	binary.LittleEndian.PutUint32(img[36:40], 0)                       // biClrImportant

	// ---------------------------------------------------------------------
	// XOR bitmap (BGRA, bottom-up)
	// ---------------------------------------------------------------------
	xorStart := 40
	for y := 0; y < height; y++ {
		// bottom-up row order in DIB
		row := xorStart + (height-1-y)*xorRowSize
		for x := 0; x < width; x++ {
			p := row + x*4
			img[p+0] = b
			img[p+1] = g
			img[p+2] = r
			img[p+3] = a
		}
	}

	// ---------------------------------------------------------------------
	// AND mask
	// ---------------------------------------------------------------------
	// All zero bits => all pixels opaque / defined by alpha.
	// The bytes are already zero-initialized, so nothing more is needed.

	return buf, nil
}

func dimByte(n int) byte {
	if n == 256 {
		return 0
	}
	return byte(n)
}

var (
	PowerOnIcon, _  = MakeIcon(16, 0x00, 0xFF, 0x00, 0xFF) // green
	PowerOffIcon, _ = MakeIcon(16, 0x80, 0x80, 0x80, 0xFF) // grey
	IdleIcon, _     = MakeIcon(16, 0xFF, 0x00, 0x00, 0xFF) // red
	DefaultIcon, _  = MakeIcon(16, 0x00, 0x00, 0xFF, 0xFF) // blue
)
