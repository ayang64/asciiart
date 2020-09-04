package asciiart

import (
	"bytes"
	"fmt"
	"image"
	"io"
)

func Encode(w io.Writer, img image.Image) error {

	// minor optimization -- store the previous color and avoid emitting escape
	// code if the color hasn't changed.

	prevTop := [3]uint32{0, 0, 0}
	prevBottom := [3]uint32{0, 0, 0}

	buf := &bytes.Buffer{}

	if _, err := buf.WriteString("\x1b[;f"); err != nil {
		return err
	}

	for y, rect := 0, img.Bounds(); y < rect.Max.Y; y += 2 {
		for x := 0; x < rect.Max.X; x++ {
			col := img.At(x, y)
			r, g, b, _ := col.RGBA()

			curTop := [3]uint32{r >> 8, g >> 8, b >> 8}

			if y == 0 || curTop != prevTop {
				if _, err := buf.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r>>8, g>>8, b>>8)); err != nil {
					return err
				}
				prevTop = curTop
			}

			col = img.At(x, y+1)
			r, g, b, _ = col.RGBA()
			curBottom := [3]uint32{r >> 8, g >> 8, b >> 8}

			if y == 0 || curBottom != prevBottom {
				if _, err := buf.WriteString(fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r>>8, g>>8, b>>8)); err != nil {
					return err
				}
				prevBottom = curBottom
			}

			buf.WriteRune('▀')
		}
	}

	buf.Write([]byte("\x1b[48;2;0;0;0m"))

	if _, err := io.Copy(w, buf); err != nil {
		return err
	}
	return nil
}
