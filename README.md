go-xwd
======

Decoder for the `xwd` (x window dump) file format. Compliant with `image.Image`.

I forked this from `sugyan/ttygif` since they had a good implementation of a decoder for xwd with some bugs.
Here I intend to fix them.

To run the tests, Go 1.16 is required because of the new `embed` package.

The original version did convert the file into a paletted image with only 216 different colors instead of 16b. I'm going to keep it like that for now.

What I changed:
 - added `go.mod`
 - reader does always assume that the color depth is 8bit (and as such reads wrong values)
 - reader fails on small files (e. g. `8colors.xwd`)
 - the reader would allocate a new byte buffer for each pixel (and not reuse the old one)
 - the reader creates a color map from the xwd file which is not used
