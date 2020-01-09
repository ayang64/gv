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
	"time"

	"github.com/ayang64/asciiart"
	"github.com/ayang64/gv/bogoscale"

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

	// convert image to slice of two color unicode glyphs.
	buf, err := asciiart.Encode(bogoscale.Scale(img, width, (height-1)*2))
	if err != nil {
		return err
	}

	if _, err := io.Copy(os.Stdout, bytes.NewReader(buf)); err != nil {
		return err
	}

	return nil
}

func cls() {
	// reset terminal to default foreground and background color.
	os.Stdout.Write([]byte("\x1b[39;m"))
	os.Stdout.Write([]byte("\x1b[49;m"))
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("need more files.")
	}

	for _, path := range os.Args[1:] {
		cls()
		if err := view(path); err != nil {
			log.Fatal(err)
		}
		time.Sleep(time.Second * 2)
	}
	cls()
}
