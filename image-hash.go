package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/nfnt/resize"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No arguments given.")
		return
	}

	filepath := os.Args[1]
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
			g >>= 10

			fmt.Println(g)
			binary.Write(buf, binary.LittleEndian, g)
			//binary.Write(buf, binary.LittleEndian, r>>12)
			//binary.Write(buf, binary.LittleEndian, g>>12)
			//binary.Write(buf, binary.LittleEndian, b>>12)
		}
	}

	h := sha1.New()
	h.Write(buf.Bytes())
	bs := h.Sum(nil)

	fmt.Println("path: ", filepath, " -- ext: ", ext)
	fmt.Printf("hash: %x\n", bs)
}
