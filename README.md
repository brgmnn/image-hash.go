# Image Hash #

Hash an image based on what it looks like, Image Hash is intended to be used
for identifying duplicate images.

## Installation ##

With a set up Go environment simply:

    go install github.com/brgmnn/image-hash

## Usage ##

To get the hash of an image run:

    $ image-hash /path/to/picture.jpg
    fe1d10cc65aa0bad

You can hash multiple images as well:

    $ image-hash cat.jpg house.jpg
    fe1d10cc65aa0bad
    15496b40ebe0fc82

Pass the `-v` flag to get the image paths as well as their hashes:

    $ image-hash cat.jpg house.jpg
    cat.jpg   fe1d10cc65aa0bad
    house.jpg 15496b40ebe0fc82

## Command line flags ##

##### -b, -bitdepth =`DEPTH`
The bitdepth represents the image bitdepth to rescale to.
Images are converted from colour to `bitdepth` grayscale. Defaults to 5.

##### -hl, hashlength =`LENGTH`
The length of image hashes in bytes. When reducing the
length of a hash, extra bytes are bitwise XORed into the hash. Defaults to 8.

##### -s, -size =`SIZE`
The target image size in pixels when rescaling images. All
images are rescaled to have a width of `size`. Images keep their aspect ratio
and so may have differing heights, but all images are rescaled to the same
width. Defaults to 4.

##### -v, -verbose
Print the image paths as well as their hashes.

## Dependencies ##

Image Hash depends on `github.com/nfnt/resize` which requires at least Go 1.1.
