package resizer

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
)

type Resizer struct{}

func New() Resizer {
	return Resizer{}
}

func (r Resizer) Fill(data []byte, width uint, height uint) ([]byte, error) {

	if width < 20 || height < 20 {
		return nil, fmt.Errorf("min width and height: 20px")
	}

	if width > 7680 || height > 7680 {
		return nil, fmt.Errorf("max width and height: 7680px")
	}

	format, err := getFormat(data)
	if err != nil {
		return nil, err
	}

	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	resizedImg := imaging.Fill(img, int(width), int(height), imaging.Center, imaging.Lanczos)

	var buf bytes.Buffer
	err = imaging.Encode(&buf, resizedImg, format)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func getFormat(data []byte) (imaging.Format, error) {
	_, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return imaging.JPEG, fmt.Errorf("failed get format: %s", err)
	}

	if format == "jpeg" {
		return imaging.JPEG, nil
	} else if format == "png" {
		return imaging.PNG, nil
	}

	return imaging.JPEG, fmt.Errorf("unsupported format: %s", format)
}
