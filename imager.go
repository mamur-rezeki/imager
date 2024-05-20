package imager

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/disintegration/imaging"
)

type Imager struct {
	Image     image.Image
	ImageType string
}

func NewImager(img image.Image) (*Imager, error) {
	return &Imager{Image: img}, nil
}

func NewImagerFromFile(location string) (*Imager, error) {
	fp, err := os.Open(location)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	img, imageType, err := image.Decode(fp)
	if err != nil {
		return nil, err
	}

	imgr, err := NewImager(img)
	imgr.ImageType = imageType

	return imgr, err
}

func NewImagerFromBytes(data []byte) (*Imager, error) {
	img, imageType, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return &Imager{Image: img, ImageType: imageType}, nil
}

const (
	IMJPEG string = "jpeg"
	IMJPG  string = "jpg"
	IMGIF  string = "gif"
	IMPNG  string = "png"
	IMWEBP string = "webp"
)

func (i *Imager) Bytes() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	var err error

	switch i.ImageType {
	case IMJPG, IMJPEG:
		err = jpeg.Encode(buf, i.Image, &jpeg.Options{Quality: 100})
	case IMPNG:
		err = png.Encode(buf, i.Image)
	case IMGIF:
		err = gif.Encode(buf, i.Image, &gif.Options{})
	}

	return buf.Bytes(), err
}

func (i *Imager) Save(location string) error {
	return imaging.Save(i.Image, location)
}

type ResizeMode int

const (
	MD_FIT ResizeMode = iota
	MD_CROP
	MD_SCALE
	MD_STRETCH
)

func (i *Imager) Resize(width, height int, modes ...ResizeMode) *Imager {
	mode := MD_FIT
	for _, md := range modes {
		mode = md
	}
	switch mode {
	case MD_SCALE:
		// Resize keeping the aspect ratio
		i.Image = imaging.Resize(i.Image, width, height, imaging.Lanczos)
	case MD_CROP:
		// Crop the image to the center
		i.Image = imaging.CropCenter(i.Image, width, height)
	case MD_FIT:
		// Fit the image within the specified dimensions, maintaining the aspect ratio
		i.Image = imaging.Fit(i.Image, width, height, imaging.Lanczos)
	case MD_STRETCH:
		// Resize to exact dimensions without keeping the aspect ratio
		i.Image = imaging.Resize(i.Image, width, height, imaging.NearestNeighbor)
	}

	return i
}

func (i *Imager) Crop(width, height int, x, y int) *Imager {
	i.Image = imaging.Crop(i.Image, image.Rect(x, y, x+width, y+height))
	return i
}

func (i *Imager) Rotate(degrees int) *Imager {
	i.Image = imaging.Rotate(i.Image, float64(degrees), &image.Uniform{})
	return i
}
