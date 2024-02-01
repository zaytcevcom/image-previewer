package resizer

import (
	"bytes"
	"image/color"
	"testing"

	"github.com/disintegration/imaging"
)

var (
	imageBytesJPEG, _ = generateImage(imaging.JPEG)
	imageBytesPNG, _  = generateImage(imaging.PNG)
	imageBytesGIF, _  = generateImage(imaging.GIF)
)

func TestFill(t *testing.T) {
	r := New()

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
			name:     "invalid image data",
			data:     []byte("invalid"),
			width:    100,
			height:   100,
			hasError: true,
		},
		{
			name:     "JPEG image",
			data:     imageBytesJPEG,
			width:    100,
			height:   100,
			hasError: false,
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
