package main

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

type Params struct {
	bitdepth   uint16
	hashlength uint16
	maxlen     int
	size       uint16
	verbose    bool
	withlog    bool
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
func hash_image(filepath string, p Params) ([]byte, error) {
	data, err := ioutil.ReadFile(filepath)

	if err != nil {
		return nil, err
	}

	imgraw, _, err := image.Decode(bytes.NewReader(data))

	if err != nil {
		return nil, err
	}

	imgstd := resize.Resize((uint)(p.size), 0, imgraw, resize.Lanczos3)
	bounds := imgstd.Bounds()
	buf := new(bytes.Buffer)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			g, _, _, _ := color.GrayModel.Convert(imgstd.At(x, y)).RGBA()
			g >>= 16 - p.bitdepth

			binary.Write(buf, binary.LittleEndian, g)
		}
	}

	h := sha1.New()
	h.Write(buf.Bytes())
	shasum := h.Sum(nil)
	sum := make([]byte, p.hashlength)

	for i := uint16(0); i < uint16(len(shasum)); i++ {
		sum[i%p.hashlength] ^= shasum[i]
	}

	return sum, nil
}

/*		--- worker() ---
 * Worker function, waits for jobs in the worker job queue and processes them
 * returning results in the dedicated channel for that job. Returns when the
 * worker job queue is closed. */
func worker(wjobs chan Job, p Params) {
	for {
		j, more := <-wjobs

		if !more {
			return
		}

		hash, err := hash_image(j.path, p)

		if err == nil {
			j.chash <- hash
		} else {
			if p.withlog {
				log.Printf("Failed to hash image '%-"+strconv.Itoa(p.maxlen)+
					"s', %s", j.path, err)
			}
			close(j.chash)
		}
	}
}

/*		--- print_hashes() ---
 * Waits on jobs to complete and prints their hashes. */
func print_hashes(mjobs chan Job, done chan bool, p Params) {
	for {
		j, more := <-mjobs

		if !more {
			done <- true
			return
		}

		hash, more := <-j.chash

		if more {
			if p.verbose {
				fmt.Printf("%-"+strconv.Itoa(p.maxlen)+"s ", j.path)
			}

			fmt.Printf("%x\n", hash)
		}
	}
}

func main() {
	bd := flag.Int("bitdepth", 5, "The number of bits to rescale each pixel "+
		"to. Must be between 1 and 16.")
	flag.IntVar(bd, "b", 5, "Short flag for 'bitdepth'.")

	sz := flag.Int("size", 4, "What 'size' to rescale images to.")
	flag.IntVar(sz, "s", 4, "Short flag for 'size'")

	hl := flag.Int("hashlength", 8, "Hash length in bytes.")
	flag.IntVar(hl, "hl", 8, "Short flag for 'hashlength'")

	v := flag.Bool("verbose", false, "Verbose output. Print image paths "+
		"along with hashes.")
	flag.BoolVar(v, "v", false, "Short flag for 'verbose'.")

	wl := flag.Bool("log", false, "Display error messages on stderr when "+
		"failing to hash an image.")
	flag.BoolVar(wl, "l", false, "Short flag for 'log'.")

	jb := flag.Int("jobs", runtime.NumCPU(), "The maximum number of hashing "+
		"jobs to run in parallel.")
	flag.IntVar(jb, "j", runtime.NumCPU(), "Short flag for 'jobs'.")

	flag.Parse()

	maxlen := 0
	if len(flag.Args()) > 0 {
		for _, fp := range flag.Args() {
			if len(fp) > maxlen {
				maxlen = len(fp)
			}
		}
	}

	params := Params{
		(uint16)(clamp(*bd, 1, 16)),
		(uint16)(clamp(*hl, 1, sha1.Size)),
		maxlen,
		(uint16)(clamp(*sz, 1, 200)),
		(bool)(*v),
		(bool)(*wl),
	}

	njobs := clamp(*jb, 1, 128)
	runtime.GOMAXPROCS(njobs)

	done := make(chan bool)
	wjobs := make(chan Job, njobs)
	mjobs := make(chan Job, njobs)

	go print_hashes(mjobs, done, params)

	for i := 0; i < njobs; i++ {
		go worker(wjobs, params)
	}

	if len(flag.Args()) > 0 {
		for _, filepath := range flag.Args() {
			job := Job{filepath, make(chan []byte)}
			wjobs <- job
			mjobs <- job
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Text()
			job := Job{line, make(chan []byte)}
			wjobs <- job
			mjobs <- job
		}
	}

	close(wjobs)
	close(mjobs)
	<-done
}
