//go:build windows
// +build windows

package main

import (
	"image/color"
	"os"
	"testing"
)

func TestResize(t *testing.T) {
	c := color.RGBA{0xAA, 0xBB, 0xCC, 0xFF}
	// `makeSquare` now returns an image.Image directly.
	src := makeSquare(c)

	// Optional sanity check – the center pixel should already have
	// the expected color before resizing.
	expected := color.RGBA{0xAA, 0xBB, 0xCC, 0xFF}
	center := src.At(src.Bounds().Dx()/2, src.Bounds().Dy()/2)
	if actual := color.RGBAModel.Convert(center).(color.RGBA); actual != expected {
		t.Fatalf("unexpected source pixel color: got %v, want %v", actual, expected)
	}

	// Resize to 16×16.
	resized := resize(src, 16)
	// The pixel at the center should equal the original color.
	expected = color.RGBA{0xAA, 0xBB, 0xCC, 0xFF}
	actualColor := resized.At(8, 8)
	actual := color.RGBAModel.Convert(actualColor).(color.RGBA)
	if actual != expected {
		t.Errorf("color mismatch after resize: got %v, want %v", actual, expected)
	}
}

func TestExportIcon(t *testing.T) {
	c := color.RGBA{0xAA, 0xBB, 0xCC, 0xFF}
	icon := MakeIcon(c)
	if len(icon) == 0 {
		t.Fatal("MakeIcon returned empty icon data")
	}

	f, err := os.Create(".\\green.ico")
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
