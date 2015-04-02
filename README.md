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

## Dependencies ##

Image Hash depends on `github.com/nfnt/resize` which requires at least Go 1.1.
