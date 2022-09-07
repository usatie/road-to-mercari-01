package imgconv

import (
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

// A Decoder reads content from r and decodes it to image.
type Decoder interface {
	Decode(r io.Reader) (image.Image, error)
}

// The DecoderFunc type is an adapter to allow the use of ordinary functions as Image decoders. If f is a function with the appropriate signature, DecoderFunc(f) is a Decoder that calls f.
type DecoderFunc func(r io.Reader) (image.Image, error)

func (f DecoderFunc) Decode(r io.Reader) (image.Image, error) { return f(r) }

// Create new Decoder from file extension
func newDecoder(ext string) (Decoder, error) {
	switch ext {
	case "jpg", "jpeg":
		return DecoderFunc(jpeg.Decode), nil
	case "png":
		return DecoderFunc(png.Decode), nil
	case "gif":
		return DecoderFunc(gif.Decode), nil
	default:
		return nil, errors.New("Invalid Extension")
	}
}
