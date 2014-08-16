package types

import (
	"errors"
	"image"
	"image/color"
	_ "image/png"
	"os"
	"path/filepath"
)

const (
	VISION_BLOCK = 0x200
	WALK_BLOCK   = 0x100
	VARIANT      = 0x10
)

var (
	translation = map[color.Color]Bits{
		//blue walls
		color.RGBA{46, 62, 95, 255}:   VISION_BLOCK | WALK_BLOCK | 0*VARIANT,
		color.RGBA{67, 92, 142, 255}:  VISION_BLOCK | WALK_BLOCK | 1*VARIANT,
		color.RGBA{76, 106, 163, 255}: VISION_BLOCK | WALK_BLOCK | 2*VARIANT,
		color.RGBA{88, 122, 188, 255}: VISION_BLOCK | WALK_BLOCK | 3*VARIANT,
		//blue nonwalls
		color.RGBA{29, 38, 61, 255}: 0 * VARIANT,
		color.RGBA{39, 50, 83, 255}: 1 * VARIANT,

		//cave walls
		color.RGBA{84, 83, 35, 255}: VISION_BLOCK | WALK_BLOCK | 0*VARIANT | 1,
		//cave nonwalls
		color.RGBA{39, 46, 26, 255}: 0*VARIANT | 1,
		color.RGBA{52, 61, 35, 255}: 1*VARIANT | 1,

		//grass walls
		color.RGBA{82, 154, 70, 255}: VISION_BLOCK | WALK_BLOCK | 0*VARIANT | 2,
		//grass nonwalls
		color.RGBA{19, 41, 20, 255}: 0*VARIANT | 2,
		color.RGBA{27, 57, 29, 255}: 1*VARIANT | 2,
		color.RGBA{35, 76, 38, 255}: 2*VARIANT | 2,
		color.RGBA{43, 92, 46, 255}: 3*VARIANT | 2,
		color.RGBA{23, 36, 46, 255}: 4*VARIANT | 2,
		color.RGBA{31, 47, 26, 255}: 5*VARIANT | 2,

		//water walls
		color.RGBA{21, 92, 127, 255}:  WALK_BLOCK | 0*VARIANT | 3,
		color.RGBA{30, 119, 149, 255}: WALK_BLOCK | 1*VARIANT | 3 | VISION_BLOCK,
		color.RGBA{34, 135, 178, 255}: WALK_BLOCK | 2*VARIANT | 3,

		//bridge nonwalls
		color.RGBA{86, 62, 35, 255}:    0*VARIANT | 4,
		color.RGBA{143, 140, 117, 255}: 1*VARIANT | 4,
	}
)

type Bits int

func (b Bits) Collides() bool {
	return (b>>8)&1 == 1
}

func (b Bits) BlocksVision() bool {
	return (b>>9)&1 == 1
}

func getName(root, path string) string {
	name, _ := filepath.Rel(root, path)
	name = filepath.ToSlash(name)
	name = name[:len(name)-len(filepath.Ext(name))]
	return name
}

func getPath(root, name, ext string) string {
	return filepath.Join(root, filepath.FromSlash(name)+ext)
}

func parseMap(path string) ([][]Bits, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	im, _, err := image.Decode(fi)
	if err != nil {
		return nil, err
	}

	m := make([][]Bits, im.Bounds().Size().Y)
	for j := im.Bounds().Min.Y; j < im.Bounds().Max.Y; j++ {
		m[j-im.Bounds().Min.Y] = make([]Bits, im.Bounds().Size().X)
		for i := im.Bounds().Min.X; i < im.Bounds().Max.X; i++ {
			m[j-im.Bounds().Min.Y][i-im.Bounds().Min.X] = translation[color.RGBAModel.Convert(im.At(i, j))]
		}
	}

	return m, nil
}

func (c *constant) loadMap(path, name string) (err error) {
	if _, alreadyThere := c.Maps[name]; alreadyThere {
		return errors.New("Duplicate map: " + name)
	}
	c.Maps[name], err = parseMap(path)
	return
}

func (w *World) ApplyPartial(pos Position, module, path string) error {
	part := Partial{
		Pos:  pos,
		Path: filepath.Join(module, PARTIAL_FOLDER, path, PARTIAL_EXT),
	}
	w.Partials = append(w.Partials, part)
	return w.applyPartial(part)
}

func (w *World) applyPartial(part Partial) error {
	m, err := parseMap(filepath.Join(w.root, MODULE_FOLDER, part.Path))
	if err != nil {
		return err
	}

	target := w.Maps[part.Pos.Mapid]
	xoff := part.Pos.X
	yoff := part.Pos.Y
	for y, column := range m {
		for x, b := range column {
			target[yoff+y][xoff+x] = b
		}
	}
	return nil
}
