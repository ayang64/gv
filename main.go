// ayan@ayan.net
package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path"

	"github.com/ayang64/gv/internal/bogoscale"

	"golang.org/x/crypto/ssh/terminal"
)

func Encode(w io.Writer, img image.Image) {
	rect := img.Bounds()

	// minor optimization -- store the previous color and avoid emitting escape
	// code if the color hasn't changed.
	prevr, prevg, prevb := uint32(0), uint32(0), uint32(0)

	buf := &bytes.Buffer{}
	os.Stdout.Write([]byte("\x1b[;f"))

	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			col := img.At(x, y)
			r, g, b, _ := col.RGBA()

			if r != prevr || g != prevg || b != prevb {
				buf.Write([]byte(fmt.Sprintf("\x1b[48;2;%d;%d;%dm", r>>8, g>>8, b>>8)))
			}
			buf.Write([]byte{' '})
			prevr, prevg, prevb = r, g, b
		}
	}

	buf.Write([]byte("\x1b[48;2;0;0;0m"))

	io.Copy(w, buf)
}

func view(p string) error {
	// FIXME: it is more reliable to examine contents of file instead of relying
	// on the extension.  A good example is RIFF files that are saved with .jpg
	// extensions.
	decmap := map[string]func(io.Reader) (image.Image, error){
		".jpeg": jpeg.Decode,
		".jpg":  jpeg.Decode,
		".png":  png.Decode,
	}

	decode, exists := decmap[path.Ext(p)]

	if exists == false {
		return fmt.Errorf("no decoder for %s", p)
	}

	r, err := os.Open(p)

	if err != nil {
		return err
	}

	img, err := decode(r)
	r.Close()

	width, height, err := terminal.GetSize(0)

	if err != nil {
		return err
	}

	outputimage := bogoscale.Scale(img, width, height-1)

	Encode(os.Stdout, outputimage)

	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("need more files.")
	}

	if err := view(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}
