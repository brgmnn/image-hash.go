package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"image"
	"image/color"
	_ "image/jpeg"

	"github.com/nfnt/resize"
)

func main() {
	var bd int
	flag.IntVar(&bd, "bitdepth", 5, "The bitdepth.")
	bitdepth := (uint16)(bd)

	flag.Parse()
	if len(flag.Args()) < 1 {
		fmt.Println("No arguments given.")
		return
	}

	filepath := flag.Args()[0]
	ext := path.Ext(filepath)

	data, _ := ioutil.ReadFile(filepath)
	imgraw, _, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		log.Fatal(err)
	}

	imgstd := resize.Resize(4, 0, imgraw, resize.Lanczos3)
	bounds := imgstd.Bounds()

	buf := new(bytes.Buffer)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := imgstd.At(x, y)
			g, _, _, _ := color.GrayModel.Convert(c).RGBA()
			g >>= 16 - (uint)(bitdepth)

			fmt.Println(g)
			binary.Write(buf, binary.LittleEndian, g)
		}
	}

	h := sha1.New()
	h.Write(buf.Bytes())
	bs := h.Sum(nil)

	fmt.Println("path: ", filepath, " -- ext: ", ext)
	fmt.Printf("hash: %x\n", bs)
}
