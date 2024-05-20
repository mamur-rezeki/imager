package imager

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"testing"
)

// createTestImage creates a simple 100x100 red image for testing
func createTestImage() image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{255, 0, 0, 255})
		}
	}
	return img
}

func TestNewImagerFromFile(t *testing.T) {
	img := createTestImage()

	// Save the test image to a temporary file
	file, err := os.CreateTemp("", "test_image_*.jpg")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		t.Fatalf("failed to encode test image: %v", err)
	}
	file.Close()

	// Test NewImagerFromFile
	imgr, err := NewImagerFromFile(file.Name())
	if err != nil {
		t.Fatalf("NewImagerFromFile returned an error: %v", err)
	}

	if imgr.Image == nil {
		t.Fatalf("NewImagerFromFile returned nil image")
	}

	if imgr.ImageType != "jpeg" {
		t.Fatalf("NewImagerFromFile returned incorrect image type: %v", imgr.ImageType)
	}
}

func TestNewImagerFromBytes(t *testing.T) {
	img := createTestImage()

	// Encode the test image to bytes
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		t.Fatalf("failed to encode test image: %v", err)
	}

	// Test NewImagerFromBytes
	imgr, err := NewImagerFromBytes(buf.Bytes())
	if err != nil {
		t.Fatalf("NewImagerFromBytes returned an error: %v", err)
	}

	if imgr.Image == nil {
		t.Fatalf("NewImagerFromBytes returned nil image")
	}

	if imgr.ImageType != "jpeg" {
		t.Fatalf("NewImagerFromBytes returned incorrect image type: %v", imgr.ImageType)
	}
}

func TestResize(t *testing.T) {
	img := createTestImage()

	imgr, err := NewImager(img)
	if err != nil {
		t.Fatalf("NewImager returned an error: %v", err)
	}

	resizedImgr := imgr.Resize(50, 50)
	if resizedImgr.Image.Bounds().Dx() != 50 || resizedImgr.Image.Bounds().Dy() != 50 {
		t.Fatalf("Resize did not return the expected dimensions: got %v", resizedImgr.Image.Bounds())
	}
}

func TestBytes(t *testing.T) {
	img := createTestImage()

	imgr, err := NewImager(img)
	if err != nil {
		t.Fatalf("NewImager returned an error: %v", err)
	}

	imgr.ImageType = IMJPEG
	data, err := imgr.Bytes()
	if err != nil {
		t.Fatalf("Bytes returned an error: %v", err)
	}

	if len(data) == 0 {
		t.Fatalf("Bytes returned empty data")
	}

	// Decode the image to ensure it was encoded correctly
	decodedImg, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Failed to decode image data: %v", err)
	}

	if decodedImg.Bounds() != img.Bounds() {
		t.Fatalf("Decoded image bounds do not match original: got %v", decodedImg.Bounds())
	}
}
