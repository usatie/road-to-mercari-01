package imgconv

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

// A Encoder writes image to w.
type Encoder interface {
	Encode(w io.Writer, m image.Image) error
}

// The EncoderFunc type is an adapter to allow the use of ordinary functions as Image encoders. If f is a function with the appropriate signature, EncoderFunc(f) is a Encoder that calls f.
type EncoderFunc func(w io.Writer, m image.Image) error

func (e EncoderFunc) Encode(w io.Writer, m image.Image) error { return e(w, m) }

// Create new Encoder from file extension
func newEncoder(ext string, quality int) (Encoder, error) {
	switch ext {
	case "jpg", "jpeg":
		return EncoderFunc(func(w io.Writer, m image.Image) error {
			var options *jpeg.Options
			if quality >= 0 && quality <= 100 {
				options = &jpeg.Options{quality}
			}
			return jpeg.Encode(w, m, options)
		}), nil
	case "gif":
		return EncoderFunc(func(w io.Writer, m image.Image) error {
			return gif.Encode(w, m, nil)
		}), nil
	case "png":
		return EncoderFunc(png.Encode), nil
	default:
		return nil, errors.New("Invalid Extension")
	}
}
