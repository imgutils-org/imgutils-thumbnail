// Package thumbnail provides simple thumbnail generation for images.
package thumbnail

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"golang.org/x/image/draw"
)

// Options configures thumbnail generation.
type Options struct {
	Width   int
	Height  int
	Quality int // JPEG quality (1-100), default 85
}

// DefaultOptions returns sensible defaults for thumbnail generation.
func DefaultOptions() Options {
	return Options{
		Width:   150,
		Height:  150,
		Quality: 85,
	}
}

// Generate creates a thumbnail from the source image.
// It maintains aspect ratio, fitting within the specified dimensions.
func Generate(src image.Image, opts Options) image.Image {
	if opts.Width <= 0 {
		opts.Width = 150
	}
	if opts.Height <= 0 {
		opts.Height = 150
	}

	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	// Calculate dimensions maintaining aspect ratio
	ratio := float64(srcW) / float64(srcH)
	var newW, newH int

	if float64(opts.Width)/float64(opts.Height) > ratio {
		newH = opts.Height
		newW = int(float64(newH) * ratio)
	} else {
		newW = opts.Width
		newH = int(float64(newW) / ratio)
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, srcBounds, draw.Over, nil)

	return dst
}

// GenerateFromFile reads an image file and generates a thumbnail.
func GenerateFromFile(path string, opts Options) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return Generate(src, opts), nil
}

// SaveJPEG saves the thumbnail as a JPEG file.
func SaveJPEG(img image.Image, w io.Writer, quality int) error {
	if quality <= 0 || quality > 100 {
		quality = 85
	}
	return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
}

// SavePNG saves the thumbnail as a PNG file.
func SavePNG(img image.Image, w io.Writer) error {
	return png.Encode(w, img)
}

// SaveGIF saves the thumbnail as a GIF file.
func SaveGIF(img image.Image, w io.Writer) error {
	return gif.Encode(w, img, nil)
}

// GenerateAndSave is a convenience function that generates a thumbnail
// and saves it to a file.
func GenerateAndSave(inputPath, outputPath string, opts Options) error {
	thumb, err := GenerateFromFile(inputPath, opts)
	if err != nil {
		return err
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	quality := opts.Quality
	if quality <= 0 {
		quality = 85
	}

	return SaveJPEG(thumb, f, quality)
}
