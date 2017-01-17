package captcha

import (
	"image"
	"image/color"
	"io"
	"image/png"
	"math/rand"
)

type Image struct {
	*image.NRGBA
	color   *color.NRGBA
	width   int //a digit width
	height  int //a digit height
	dotsize int
}

func NewImage(digits []byte, width, height int) *Image {
	img := new(Image)
	r := image.Rect(img.width, img.height, STD_WIDTH, STD_HEIGHT)
	img.NRGBA = image.NewNRGBA(r)

	img.color = &color.NRGBA{
		R:uint8(rand.Intn(129)),
		G:uint8(rand.Intn(129)),
		B:uint8(rand.Intn(129)),
		A:0xFF,
	}
	// Draw background (10 random circles of random brightness)
	img.calculate_sizes(width, height, len(digits))
	img.fill_with_circles(10, img.dotsize)

	maxx := width - (img.width+img.dotsize)*len(digits) - img.dotsize
	maxy := height - img.height - img.dotsize*2

	x := rnd(img.dotsize*2, maxx)
	y := rnd(img.dotsize*2, maxy)

	// Draw digits.
	for _, n := range digits {
		img.draw_digit(font[n], x, y)
		x += img.width + img.dotsize
	}

	// Draw strike-through line.
	img.strike_through()
	return img
}

func (img *Image) WriteTo(w io.Writer) (int64, error) {
	return 0, png.Encode(w, img)
}

func (img *Image) calculate_sizes(width, height, ncount int) {

	// Goal: fit all digits inside the image.
	var border int
	if width > height {
		border = height / 5
	} else {
		border = width / 5
	}
	// Convert everything to floats for calculations.
	w := float64(width - border*2)  //268
	h := float64(height - border*2) //48
	// fw takes into account 1-dot spacing between digits.

	fw := float64(FONT_WIDTH) + 1 //6

	fh := float64(FONT_HEIGHT) //8
	nc := float64(ncount)      //7

	// Calculate the width of a single digit taking into account only the
	// width of the image.
	nw := w / nc //38
	// Calculate the height of a digit from this width.
	nh := nw * fh / fw //51

	// Digit too high?

	if nh > h {
		// Fit digits based on height.
		nh = h //nh = 44
		nw = fw / fh * nh
	}
	// Calculate dot size.
	img.dotsize = int(nh / fh)
	// Save everything, making the actual width smaller by 1 dot to account
	// for spacing between digits.
	img.width = int(nw)
	img.height = int(nh) - img.dotsize
}

func (img *Image) fill_with_circles(n, maxradius int) {
	img_color := img.color
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	for i := 0; i < n; i++ {
		setRandomBrightness(img_color, 255)
		r := rnd(1, maxradius)
		img.drawCircle(img_color, rnd(r, maxx-r), rnd(r, maxy-r), r)
	}
}

func (img *Image) draw_horiz_line(color color.Color, fromX, toX, y int) {
	for x := fromX; x <= toX; x++ {
		img.Set(x, y, color)
	}
}

func (img *Image) drawCircle(color color.Color, x, y, radius int) {
	f := 1 - radius
	dfx := 1
	dfy := -2 * radius
	xx := 0
	yy := radius

	img.Set(x, y+radius, color)
	img.Set(x, y-radius, color)
	img.draw_horiz_line(color, x-radius, x+radius, y)

	for xx < yy {
		if f >= 0 {
			yy--
			dfy += 2
			f += dfy
		}
		xx++
		dfx += 2
		f += dfx
		img.draw_horiz_line(color, x-xx, x+xx, y+yy)
		img.draw_horiz_line(color, x-xx, x+xx, y-yy)
		img.draw_horiz_line(color, x-yy, x+yy, y+xx)
		img.draw_horiz_line(color, x-yy, x+yy, y-xx)
	}
}

func (img *Image) strike_through() {
	r := 0
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	y := rnd(maxy/3, maxy-maxy/3)
	for x := 0; x < maxx; x += r {
		r = rnd(1, img.dotsize/3)
		y += rnd(-img.dotsize/2, img.dotsize/2)
		if y <= 0 || y >= maxy {
			y = rnd(maxy/3, maxy-maxy/3)
		}
		img.drawCircle(img.color, x, y, r)
	}
}

func (img *Image) draw_digit(digit []byte, x, y int) {
	skf := rand.Float64() * float64(rnd(-MAX_SKEW, MAX_SKEW))
	xs := float64(x)
	minr := img.dotsize / 2               // minumum radius
	maxr := img.dotsize/2 + img.dotsize/4 // maximum radius
	y += rnd(-minr, minr)
	for yy := 0; yy < FONT_HEIGHT; yy++ {
		for xx := 0; xx < FONT_WIDTH; xx++ {
			if digit[yy*FONT_WIDTH + xx] != BLACK_CHAR {
				continue
			}
			// Introduce random variations.
			or := rnd(minr, maxr)
			ox := x + (xx * img.dotsize) + rnd(0, or/2)
			oy := y + (yy * img.dotsize) + rnd(0, or/2)

			img.drawCircle(img.color, ox, oy, or)
		}
		xs += skf
		x = int(xs)
	}
}
