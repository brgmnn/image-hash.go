package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"
	"strconv"

	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"

	"github.com/nfnt/resize"
)

type Job struct {
	path  string
	chash chan []byte
}

/*		--- clamp() ---
 * Clamps an integer to a minimum of 'lo' or a maximum of 'hi'. */
func clamp(num, lo, hi int) int {
	if num > hi {
		return hi
	} else if num < lo {
		return lo
	}

	return num
}

/*		--- min() ---
 * Returns the smaller integer of two integers. */
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/*		--- hash_image() ---
 * Returns the hash of an image given its path. */
func hash_image(filepath string, width, bitdepth uint16) []byte {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		log.Fatal(err)
	}

	imgraw, _, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		log.Fatal(err)
	}

	imgstd := resize.Resize((uint)(width), 0, imgraw, resize.Lanczos3)
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
	return h.Sum(nil)
}

/*		--- worker() ---
 * Worker function, waits for jobs in the worker job queue and processes them
 * returning results in the dedicated channel for that job. Returns when the
 * worker job queue is closed. */
func worker(wjobs chan Job, width, bitdepth uint16) {
	for {
		j, more := <-wjobs

		if !more {
			return
		}

		j.chash <- hash_image(j.path, width, bitdepth)
	}
}

/*		--- print_hashes() ---
 * Waits on jobs to complete and prints their hashes. */
func print_hashes(mjobs chan Job, done chan bool, verbose bool, maxlen int) {
	for {
		j, more := <-mjobs

		if !more {
			done <- true
			return
		}

		if verbose {
			fmt.Printf("%-"+strconv.Itoa(maxlen)+"s ", j.path)
		}

		fmt.Printf("%x\n", <-j.chash)
	}
}

func main() {
	bd := flag.Int("bitdepth", 5, "The number of bits to rescale each pixel "+
		"to. Must be between 1 and 16.")
	wd := flag.Int("size", 4, "What 'size' to rescale images to.")
	v := flag.Bool("v", false, "Verbose output. Print image paths along "+
		"with hashes.")

	flag.Parse()

	bitdepth := (uint16)(clamp(*bd, 1, 16))
	width := (uint16)(clamp(*wd, 1, 200))
	verbose := (bool)(*v)

	if len(flag.Args()) < 1 {
		fmt.Println("No images given.")
		return
	}

	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)

	maxlen := 0
	for _, fp := range flag.Args() {
		if len(fp) > maxlen {
			maxlen = len(fp)
		}
	}

	done := make(chan bool)
	wjobs := make(chan Job, 4)
	mjobs := make(chan Job, 4)

	go print_hashes(mjobs, done, verbose, maxlen)

	for i := 0; i < cores; i++ {
		go worker(wjobs, width, bitdepth)
	}

	for _, filepath := range flag.Args() {
		job := Job{filepath, make(chan []byte)}
		wjobs <- job
		mjobs <- job
	}

	close(wjobs)
	close(mjobs)
	<-done
}
