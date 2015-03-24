# Image Hash #

Hash an image based on what it looks like, Image Hash is intended to be used
for identifying duplicate images.

## Installation ##

With a set up Go environment simply:

    go install github.com/brgmnn/image-hash

## Usage ##

To get the hash of an image run:

    $ image-hash /path/to/picture.jpg
    88014ee65ffc1afc519edfde3a561151278281f4

You can hash multiple images as well:

    $ image-hash cat.jpg house.jpg
    88014ee65ffc1afc519edfde3a561151278281f4
    62680e03c4499e23c167e5492fa962a1b646800a

Pass the `-v` flag to get the image paths as well as their hashes:

    $ image-hash cat.jpg house.jpg
    cat.jpg   88014ee65ffc1afc519edfde3a561151278281f4
    house.jpg 62680e03c4499e23c167e5492fa962a1b646800a

## Dependencies ##

Image Hash depends on `github.com/nfnt/resize` which requires at least Go 1.1.
