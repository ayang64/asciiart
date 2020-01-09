package asciiart

import (
	"bytes"
	"fmt"
	"image"
	"io"
)

// Encode returns a byte slice that is sutable for writing to a terminal that
// represents the contents of the supplied image.Image
func Encode(img image.Image) ([]byte, error) {
	return EncodeBuffer([]byte{}, img)
}

// Encode returns a byte slice that is sutable for writing to a terminal that
// represents the contents of the supplied image.Image
//
// It writes to the supplied byte slice and appends to it as needed.  The
// returned byte slice can be reused as the first parameter in subsequent calls
// to prevent allocations.
//
func EncodeBuffer(b []byte, img image.Image) ([]byte, error) {
	buf := bytes.NewBuffer(b)
	buf.Write([]byte("\x1b[;f"))

	esc := func(w io.Writer, prev *[3]uint32, f string, img image.Image, x int, y int) {
		col := img.At(x, y)
		r, g, b, _ := col.RGBA()
		if cur := [3]uint32{r >> 8, g >> 8, b >> 8}; y != 0 || cur == *prev {
			buf.Write([]byte(fmt.Sprintf(f, cur[0], cur[1], cur[2])))
			*prev = cur
		}
	}

	// minor optimization -- store the previous color and avoid emitting escape
	// code if the color hasn't changed.
	for y, prevTop, prevBottom, rect := 0, [3]uint32{0, 0, 0}, [3]uint32{0, 0, 0}, img.Bounds(); y < rect.Max.Y; y += 2 {
		for x := 0; x < rect.Max.X; x++ {
			esc(buf, &prevTop, "\x1b[38;2;%d;%d;%dm", img, x, y)
			esc(buf, &prevBottom, "\x1b[48;2;%d;%d;%dm", img, x, y)
			buf.WriteRune('▀')
		}
	}

	buf.Write([]byte("\x1b[48;2;0;0;0m"))

	return buf.Bytes(), nil
}
