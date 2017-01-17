package captcha

import (
	"io"
	"image/color"
	"time"
	"math/rand"
	crand "crypto/rand"
)

const (
	STD_WIDTH  = 100
	STD_HEIGHT = 40
	MAX_SKEW   = 2
)

const (
	FONT_WIDTH  = 5
	FONT_HEIGHT = 8
	BLACK_CHAR  = 1
)

var font = [][]byte{
	{ // 0
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
	},
	{ // 1
		0, 0, 1, 0, 0,
		0, 1, 1, 0, 0,
		1, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		1, 1, 1, 1, 1,
	},
	{ // 2
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 1, 1,
		0, 1, 1, 0, 0,
		1, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
	},
	{ // 3
		1, 1, 1, 1, 0,
		0, 0, 0, 0, 1,
		0, 0, 0, 1, 0,
		0, 1, 1, 1, 0,
		0, 0, 0, 1, 0,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		1, 1, 1, 1, 0,
	},
	{ // 4
		1, 0, 0, 1, 0,
		1, 0, 0, 1, 0,
		1, 0, 0, 1, 0,
		1, 0, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 0, 0, 1, 0,
		0, 0, 0, 1, 0,
		0, 0, 0, 1, 0,
	},
	{ // 5
		1, 1, 1, 1, 1,
		1, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 1, 0,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		1, 1, 1, 1, 0,
	},
	{ // 6
		0, 0, 1, 1, 1,
		0, 1, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
	},
	{ // 7
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 1, 0,
		0, 0, 1, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
	},
	{ // 8
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
	},
	{ // 9
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 1, 0, 0, 1,
		0, 1, 1, 1, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		1, 1, 1, 1, 0,
	},
}

func init() {
	rand.Seed(int64(time.Second))
}



func setRandomBrightness(c *color.NRGBA, max uint8) {
	minc := min3(c.R, c.G, c.B)
	maxc := max3(c.R, c.G, c.B)
	if maxc > max {
		return
	}
	n := rand.Intn(int(max - maxc)) - int(minc)
	c.R = uint8(int(c.R) + n)
	c.G = uint8(int(c.G) + n)
	c.B = uint8(int(c.B) + n)
}

func min3(x, y, z uint8) (o uint8) {
	o = x
	if y < o {
		o = y
	}
	if z < o {
		o = z
	}
	return
}

func max3(x, y, z uint8) (o uint8) {
	o = x
	if y > o {
		o = y
	}
	if z > o {
		o = z
	}
	return
}

// rnd returns a random number in range [from, to].
func rnd(from, to int) int {
	//println(to+1-from)
	return rand.Intn(to + 1 - from) + from
}

// Standard characters allowed in uniuri string.
var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

// NewLen returns a new random string of the provided length, consisting of
// standard characters.
func NewLen(length int) string {
	return NewLenChars(length, StdChars)
}

// NewLenChars returns a new random string of the provided length, consisting
// of the provided byte slice of allowed characters (maximum 256).
func NewLenChars(length int, chars []byte) string {
	b := make([]byte, length)
	r := make([]byte, length + (length / 4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(crand.Reader, r); err != nil {
			panic("error reading from random source: " + err.Error())
		}
		for _, c := range r {
			if c >= maxrb {
				// Skip this number to avoid modulo bias.
				continue
			}
			b[i] = chars[c % clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
	panic("unreachable")
}

/*
func pic(w http.ResponseWriter, req *http.Request) {
	d := make([]byte, 4)
	s := NewLen(4)
	ss := ""
	d = []byte(s)
	for v := range d {
		d[v] %= 10
		ss += strconv.FormatInt(int64(d[v]), 32)
	}
	w.Header().Set("Content-Type", "image/png")
	NewImage(d, 100, 40).WriteTo(w)
	fmt.Println(ss)
}

func index(w http.ResponseWriter, req *http.Request) {
	str := "<meta charset=\"utf-8\"><h3>golang 图片验证码例子</h3><img border=\"1\" src=\"/pic\" alt=\"图片验证码\" onclick=\"this.src='/pic'\" />"
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(str))
}

func main() {
	http.HandleFunc("/pic", pic)
	http.HandleFunc("/", index)
	s := &http.Server{
		Addr: 		 	":8080",
		ReadTimeout: 	30*time.Second,
		WriteTimeout: 	30*time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
*/
