go-xwd
======

Decoder for the `xwd` (x window dump) file format. Compliant with `image.Image`.

I forked this from `sugyan/ttygif` since they had a good implementation of a decoder for xwd with some bugs.
Here I intend to fix them.

To run the tests, Go 1.16 is required because of the new `embed` package.

The original version did convert the file into a paletted image with only 216 different colors instead of 16b. I'm going to keep it like that for now.

What I changed:
 - added `go.mod`
 - the original reader did not work with colormapped pixmaps. It always expected a raw pixmap.
 - reader does always assume that the color depth is 8bit (and as such reads wrong values)
 - reader fails on small files (e. g. `8colors.xwd`)
 - the reader would allocate a new byte buffer for each pixel (and not reuse the old one)
 - the reader creates a color map from the xwd file which is not used

TODO:
 - Currently, every color is simply converted to an `XWDColor` which has 16b per channel. Most of the times, we are dealing with 8bpc colors instead, so we should invent a new color model that only has 8b per channel

 The xwd format
 --------------

It's pretty safe to say that `xwd` is an obsolete format. It is not portable and there's zero documentation on it. Here's what I found so far:
 - The header consists of 25 uint32 values = 100 bytes
 - xwd version 7 is the newset. It runs on xorg v11 while any older version is used by older versions of the x server [source](https://gitlab.freedesktop.org/xorg/app/xwd/-/blob/master/xwd.c#L52-58)
 - At the end of the header, the window's name gets appended. [The `xwd` writer](https://gitlab.freedesktop.org/xorg/app/xwd/-/blob/master/xwd.c#L477-478) does these calls:
```
fwrite((char *) &header, SIZEOF(XWDheader), 1, out)
fwrite(win_name, win_name_size, 1, out)
```
The size of the window name is equal to header size - 100. [The string is null-terminated](https://gitlab.freedesktop.org/xorg/app/xwd/-/blob/master/xwd.c#L348-349) which we don't need in Go, so we have to strip this one away as well.
 - **Endianness:** The header has two values, `byte_order` and `bitmap_bit_order`, both directly mapped from the `XImage`

The colormap is a list of colors. Each color is defined like this:
```
uint32 pixel
uint16 red
uint16 green
uint16 blue
uint8 flags
uint8 pad
```
The colormap is useless as it only contains grayscale values.
For reading the colormap, I use the hardcoded types and not any of the data in the header.

The image itself is defined as an `XPointer` which is a `char*`
