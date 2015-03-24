package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
)

/*		--- clamp() ---
 * Clamps an integer to a minimum of 'lo' or a maximum of 'hi'. */
func clamp(num int, lo int, hi int) int {
	if num > hi {
		return hi
	} else if num < lo {
		return lo
	}

	return num
}

func main() {
	bd := flag.Int("bitdepth", 5, "The number of bits to rescale each pixel "+
		"to. Must be between 1 and 16.")
	wd := flag.Int("size", 4, "What 'size' to rescale images to.")
	v := flag.Bool("v", false, "Verbose output. Print image paths along with hashes.")

	flag.Parse()

	bitdepth := (uint16)(clamp(*bd, 1, 16))
	width := (uint)(clamp(*wd, 1, 200))
	verbose := (bool)(*v)

	if len(flag.Args()) < 1 {
		fmt.Println("No images given.")
		return
	}

	maxlen := 0
	for _, fp := range flag.Args() {
		if len(fp) > maxlen {
			maxlen = len(fp)
		}
	}

	for _, filepath := range flag.Args() {
		data, err := ioutil.ReadFile(filepath)

		if err != nil {
			log.Fatal(err)
		}

		imgraw, _, err := image.Decode(bytes.NewReader(data))

		if err != nil {
			log.Fatal(err)
		}

		imgstd := resize.Resize(width, 0, imgraw, resize.Lanczos3)
		bounds := imgstd.Bounds()

		buf := new(bytes.Buffer)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				g, _, _, _ := color.GrayModel.Convert(imgstd.At(x, y)).RGBA()
				g >>= 16 - bitdepth

				binary.Write(buf, binary.LittleEndian, g)
			}
		}

		h := sha1.New()
		h.Write(buf.Bytes())
		bs := h.Sum(nil)

		if verbose {
			fmt.Printf("%-"+strconv.Itoa(maxlen)+"s ", filepath)
		}

		fmt.Printf("%x\n", bs)
	}
}
