package types

import (
	"image"
	"image/color"
	_ "image/png"
	"os"
)

var (
	translation = map[color.Color]Bits{
		//blue walls
		color.RGBA{46, 62, 95, 255}:   0x300,
		color.RGBA{67, 92, 142, 255}:  0x310,
		color.RGBA{76, 106, 163, 255}: 0x320,
		color.RGBA{88, 122, 188, 255}: 0x330,
		//blue nonwalls
		color.RGBA{29, 38, 61, 255}: 0x000,
		color.RGBA{39, 50, 83, 255}: 0x010,

		//cave walls
		color.RGBA{84, 83, 35, 255}: 0x301,
		//cave nonwalls
		color.RGBA{39, 46, 26, 255}: 0x001,
		color.RGBA{52, 61, 35, 255}: 0x011,

		//grass walls
		color.RGBA{82, 154, 70, 255}: 0x302,
		//grass nonwalls
		color.RGBA{19, 41, 20, 255}: 0x002,
		color.RGBA{27, 57, 29, 255}: 0x012,
		color.RGBA{35, 76, 38, 255}: 0x022,
		color.RGBA{43, 92, 46, 255}: 0x032,
		color.RGBA{23, 36, 46, 255}: 0x042,
		color.RGBA{31, 47, 26, 255}: 0x052,

		//water walls
		color.RGBA{21, 92, 127, 255}:  0x103,
		color.RGBA{30, 119, 149, 255}: 0x313,
		color.RGBA{34, 135, 178, 255}: 0x123,

		//bridge nonwalls
		color.RGBA{86, 62, 35, 255}:    0x004,
		color.RGBA{143, 140, 117, 255}: 0x014,
	}
)

type Bits int

func (b Bits) Collides() bool {
	return (b>>8)&1 == 1
}

func (b Bits) BlocksVision() bool {
	return (b>>9)&1 == 1
}

func parseMap(path string) [][]Bits {
	fi, err := os.Open(path)
	d(err)

	im, _, err := image.Decode(fi)
	d(err)

	m := make([][]Bits, im.Bounds().Size().Y)
	for j := im.Bounds().Min.Y; j < im.Bounds().Max.Y; j++ {
		m[j-im.Bounds().Min.Y] = make([]Bits, im.Bounds().Size().X)
		for i := im.Bounds().Min.X; i < im.Bounds().Max.X; i++ {
			m[j-im.Bounds().Min.Y][i-im.Bounds().Min.X] = translation[color.RGBAModel.Convert(im.At(i, j))]
		}
	}

	return m
}

func d(err error) {
	if err != nil {
		panic(err)
	}
}
