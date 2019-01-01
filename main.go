// ayan@ayan.net
package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path"

	"github.com/ayang64/gv/internal/asciiart"
	"github.com/ayang64/gv/internal/bogoscale"

	"golang.org/x/crypto/ssh/terminal"
)

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

	outputimage := bogoscale.Scale(img, width, (height-1)*2)

	asciiart.Encode(os.Stdout, outputimage)

	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("need more files.")
	}

	if err := view(os.Args[1]); err != nil {
		log.Fatal(err)
	}

	// reset terminal to default foreground and background color.
	os.Stdout.Write([]byte("\x1b[39;m"))
	os.Stdout.Write([]byte("\x1b[49;m"))
}
