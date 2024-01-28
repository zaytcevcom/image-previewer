package resizer

import (
	"bytes"
	"github.com/disintegration/imaging"
	"image/color"
	"testing"
)

func TestFill(t *testing.T) {
	r := New()

	imageBytesJPEG, err := generateImage(imaging.JPEG)
	if err != nil {
		t.Errorf("Unable to encode image: %v", err)
	}

	imageBytesPNG, err := generateImage(imaging.PNG)
	if err != nil {
		t.Errorf("Unable to encode image: %v", err)
	}

	imageBytesGIF, err := generateImage(imaging.GIF)
	if err != nil {
		t.Errorf("Unable to encode image: %v", err)
	}

	tests := []struct {
		name     string
		data     []byte
		width    uint
		height   uint
		hasError bool
	}{
		{
			name:     "width and height are too small",
			data:     imageBytesJPEG,
			width:    10,
			height:   10,
			hasError: true,
		},
		{
			name:     "width and height are too big",
			data:     imageBytesJPEG,
			width:    8000,
			height:   8000,
			hasError: true,
		},
		{
			name:     "width and height are just at the maximum",
			data:     imageBytesJPEG,
			width:    7680,
			height:   7680,
			hasError: false,
		},
		{
			name:     "width and height are just over the minimum",
			data:     imageBytesJPEG,
			width:    20,
			height:   20,
			hasError: false,
		},
		{
			name:     "width and height within normal range",
			data:     imageBytesJPEG,
			width:    100,
			height:   100,
			hasError: false,
		},
		{
			name:     "invalid image data",
			data:     []byte("invalid"),
			width:    100,
			height:   100,
			hasError: true,
		},
		{
			name:     "PNG image",
			data:     imageBytesPNG,
			width:    100,
			height:   100,
			hasError: false,
		},
		{
			name:     "GIF image",
			data:     imageBytesGIF,
			width:    100,
			height:   100,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := r.Fill(tt.data, tt.width, tt.height)
			if tt.hasError && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("did not expect an error but got %v", err)
			}
		})
	}
}

func generateImage(format imaging.Format) ([]byte, error) {
	img := imaging.New(50, 50, color.Black)
	buf := new(bytes.Buffer)
	err := imaging.Encode(buf, img, format)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}