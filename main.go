//go:build example
// +build example

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/zfedoran/go-wfc-example/assets"
	"github.com/zfedoran/go-wfc/pkg/wfc"
)

const (
	tile_size          = 200
	output_tile_width  = 8
	output_tile_height = 8
	screen_width       = output_tile_width * 40
	screen_height      = output_tile_height * 40
)

var (
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

var (
	constraintFunc = wfc.GetConstraintFunc(2)
)

func init() {
	emptyImage.Fill(color.White)
}

type Game struct {
	count, seed int

	tiles   []image.Image
	mapping map[*wfc.Module]*ebiten.Image

	wave *wfc.Wave
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screen_width, screen_height
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	tile_width := screen_width / output_tile_width
	tile_height := screen_height / output_tile_height

	for _, slot := range g.wave.History {
		x := float32(slot.X*tile_width) - 1
		y := float32(slot.Y*tile_height) - 1

		var path vector.Path
		path.MoveTo(x, y)
		path.LineTo(x+float32(tile_width)+1, y)
		path.LineTo(x+float32(tile_width)+1, y+float32(tile_height)+1)
		path.LineTo(x, y+float32(tile_height)+1)

		op := &ebiten.DrawTrianglesOptions{FillRule: ebiten.FillAll}
		vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
		for i := range vs {
			vs[i].SrcX = 1
			vs[i].SrcY = 1
			vs[i].ColorR = 227 / float32(0xff)
			vs[i].ColorG = 188 / float32(0xff)
			vs[i].ColorB = 61 / float32(0xff)
		}
		screen.DrawTriangles(vs, is, emptySubImage, op)
	}

	for _, slot := range g.wave.PossibilitySpace {
		if len(slot.Superposition) == 1 {
			x := slot.X * tile_width
			y := slot.Y * tile_height

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(tile_width)/(tile_size+5), float64(tile_height)/(tile_size+5))
			op.GeoM.Translate(float64(x), float64(y))

			img := g.getImage(slot.Superposition[0])
			screen.DrawImage(img, op)
		} else {
			num_tiles := len(slot.Superposition)
			num_tiles_sq := int(math.Ceil(math.Sqrt(float64(num_tiles)))) + 1
			module_width := tile_width / num_tiles_sq
			module_height := tile_height / num_tiles_sq
			for j, m := range slot.Superposition {
				x := (j%num_tiles_sq)*module_width + slot.X*tile_width + 2
				y := (j/num_tiles_sq)*module_height + slot.Y*tile_height + 2

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(float64(module_width)/screen_width, float64(module_height)/screen_height)
				op.GeoM.Translate(float64(x), float64(y))

				img := g.getImage(m)
				screen.DrawImage(img, op)
			}
		}
	}

	for i, slot := range g.wave.History {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d", i), slot.X*tile_width, slot.Y*tile_height)
	}

	if g.wave.IsCollapsed() {
		ebitenutil.DebugPrint(screen, fmt.Sprintf(" Seed: 0x%X", g.seed))
	} else {
		ebitenutil.DebugPrint(screen, fmt.Sprintf(" FPS: %0.2f", ebiten.CurrentTPS()))
	}
}

func (g *Game) load() {
	folder := "tiles"
	store := assets.Store

	dir, err := store.ReadDir(folder)
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		item := filepath.Join(folder, file.Name())

		if filepath.Ext(file.Name()) != ".png" {
			continue
		}

		handle, err := store.Open(item)
		if err != nil {
			panic(err)
		}
		defer handle.Close()

		img, err := png.Decode(handle)
		if err != nil {
			panic(err)
		}

		g.tiles = append(g.tiles, img)
	}

	g.wave = wfc.NewWithCustomConstraints(g.tiles, output_tile_width, output_tile_height, constraintFunc)
	g.wave.IsPossibleFn = func(m *wfc.Module, from, to *wfc.Slot, d wfc.Direction) bool {
		time.Sleep(time.Millisecond / 10)
		res := wfc.DefaultIsPossibleFunc(m, from, to, d)
		return res
	}
}

func (g *Game) collapse() {
	for {
		if g.wave.IsCollapsed() {
			time.Sleep(time.Second * 3)
			g.initialize()
		}

		err := g.wave.Collapse(1)
		if err != nil {
			// don't do anything we want to see the contradition
		}
	}
}

func (g *Game) initialize() {
	g.seed = int(time.Now().UnixNano())
	g.wave.Initialize(g.seed)
}

func (g *Game) getImage(m *wfc.Module) *ebiten.Image {
	if g.mapping == nil {
		g.mapping = make(map[*wfc.Module]*ebiten.Image)
	}

	if m == nil {
		return nil
	}

	if g.mapping[m] == nil {
		g.mapping[m] = ebiten.NewImageFromImage(m.Image)
	}

	return g.mapping[m]
}

func main() {
	g := &Game{}

	g.load()
	g.initialize()
	go func() {
		g.collapse()
	}()

	ebiten.SetWindowSize(screen_width, screen_height)
	ebiten.SetWindowTitle("go-wfc-algorithm")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
