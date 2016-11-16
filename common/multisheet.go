package common

import (
	"log"

	"engo.io/engo"
	"engo.io/gl"
)

// Subsheet is a class that stores a set of tiles from a file, used by tilemaps and animations
type Subsheet struct {
	parent                *MultiSheet
	width, height         float32         // The dimensions of the total texture
	cellWidth, cellHeight int             // The dimensions of the cells
	offsetX, offsetY      int             // Offset where the subsheet begins relative to 0,0
	cache                 map[int]Texture // The cell cache cells
}

func NewSubsheet(width, height float32, offsetX, offsetY, cellWidth, cellHeight int) *Subsheet {
	return &Subsheet{
		offsetX: offsetX, offsetY: offsetY,
		width: width, height: height,
		cellWidth: cellWidth, cellHeight: cellHeight,
		cache: make(map[int]Texture),
	}
}

// Cell gets the region at the index i, updates and pulls from cache if need be
func (s *Subsheet) Cell(index int) Texture {
	if r, ok := s.cache[index]; ok {
		return r
	}

	cellsPerRow := int(s.Width())
	var x float32 = float32((index % cellsPerRow) * s.cellWidth)
	var y float32 = float32((index / cellsPerRow) * s.cellHeight)
	s.cache[index] = Texture{id: s.parent.texture, width: float32(s.cellWidth), height: float32(s.cellHeight), viewport: engo.AABB{
		engo.Point{x / s.width, y / s.height},
		engo.Point{(x + float32(s.cellWidth)) / s.width, (y + float32(s.cellHeight)) / s.height},
	}}

	return s.cache[index]
}

func (s *Subsheet) Drawable(index int) Drawable {
	return s.Cell(index)
}

func (s *Subsheet) Drawables() []Drawable {
	drawables := make([]Drawable, s.CellCount())

	for i := 0; i < s.CellCount(); i++ {
		drawables[i] = s.Drawable(i)
	}

	return drawables
}

func (s *Subsheet) CellCount() int {
	return int(s.Width()) * int(s.Height())
}

func (s *Subsheet) Cells() []Texture {
	cellsNo := s.CellCount()
	cells := make([]Texture, cellsNo)
	for i := 0; i < cellsNo; i++ {
		cells[i] = s.Cell(i)
	}

	return cells
}

// Width is the amount of tiles on the x-axis of the multisheet
func (s Subsheet) Width() float32 {
	return s.width / float32(s.cellWidth)
}

// Height is the amount of tiles on the y-axis of the multisheet
func (s Subsheet) Height() float32 {
	return s.height / float32(s.cellHeight)
}

type MultiSheet struct {
	texture   *gl.Texture
	subsheets []*Subsheet
}

func NewMultiSheetFromTexture(tr *TextureResource, subsheets []*Subsheet) *MultiSheet {
	m := &MultiSheet{
		texture:   tr.Texture,
		subsheets: subsheets,
	}
	for _, subsheet := range subsheets {
		subsheet.parent = m
	}
	return m
}

// NewMultiSheetFromFile is a simple handler for creating a new spritesheet from a file
// textureName is the name of a texture already preloaded with engo.Files.Add
func NewMultiSheetFromFile(textureName string, subsheets []*Subsheet) *MultiSheet {
	res, err := engo.Files.Resource(textureName)
	if err != nil {
		log.Println("[WARNING] [NewMultiSheetFromFile]: Received error:", err)
		return nil
	}

	img, ok := res.(TextureResource)
	if !ok {
		log.Println("[WARNING] [NewMultiSheetFromFile]: Resource not of type `TextureResource`:", textureName)
		return nil
	}

	return NewMultiSheetFromTexture(&img, subsheets)
}

func (ms *MultiSheet) CellIndex(subsheetIndex, regionIndex int) int {
	index := regionIndex
	for i := 0; i < subsheetIndex; i++ {
		index += i * ms.subsheets[i].CellCount()
	}

	return index
}

func (ms *MultiSheet) Drawables() []Drawable {
	drawables := []Drawable{}
	for _, subsheet := range ms.subsheets {
		for i := 0; i < subsheet.CellCount(); i++ {
			drawables = append(drawables, subsheet.Cell(i))
		}
	}
	return drawables
}
