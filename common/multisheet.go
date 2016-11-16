package common

import (
	"log"

	"engo.io/engo"
	"engo.io/gl"
)

// Subsheet is an animation set within a larger multisheet
type Subsheet struct {
	width, height         float32 // dimensions of the subsheet
	cellWidth, cellHeight int     // The dimensions of the cells within the subsheet
	offsetX, offsetY      int     // Offset relative to 0,0 where the subsheet begins
}

func NewSubsheet(width, height float32, cellWidth, cellHeight, offsetX, offsetY int) *Subsheet {
	return &Subsheet{
		width:      width,
		height:     height,
		cellWidth:  cellWidth,
		cellHeight: cellHeight,
		offsetX:    offsetX,
		offsetY:    offsetX,
	}
}

// Multisheet is a class that stores a set of tiles from a file, used by tilemaps and animations
type Multisheet struct {
	texture   *gl.Texture         // The original texture
	subsheets map[string]Subsheet // The subsheets that make up the multisheet
	cache     map[int]Texture     // The cell cache cells
}

func NewMultisheetFromTexture(tr *TextureResource, subsheets map[string]*Subsheet) *Multisheet {
	return &Multisheet{texture: tr.Texture,
		subsheets: subsheets,
		cache:     make(map[int]Texture),
	}
}

// NewMultisheetFromFile is a simple handler for creating a new multisheet from a file
// textureName is the name of a texture already preloaded with engo.Files.Add
func NewMultisheetFromFile(textureName string, subsheets map[string]*Subsheet) *Multisheet {
	res, err := engo.Files.Resource(textureName)
	if err != nil {
		log.Println("[WARNING] [NewMultisheetFromFile]: Received error:", err)
		return nil
	}

	img, ok := res.(TextureResource)
	if !ok {
		log.Println("[WARNING] [NewMultisheetFromFile]: Resource not of type `TextureResource`:", textureName)
		return nil
	}

	return NewMultisheetFromTexture(&img, subsheets)
}

// Cell gets the region at the index i, updates and pulls from cache if need be
func (s *Multisheet) Cell(animationName string, index int) Texture {
	if r, ok := s.cache[index]; ok {
		return r
	}

	subsheet, ok := s.subsheets[animationName]
	if !ok {
		log.Fatalf("%s is not a valid aniation for %v", animationName, s)
	}

	cellsPerRow := int(subsheet.width)
	var x float32 = float32((index%cellsPerRow)*subsheet.cellWidth + subsheet.offsetX)
	var y float32 = float32((index/cellsPerRow)*subsheet.cellHeight + subsheet.offsetY)
	s.cache[index] = Texture{
		id:     s.texture,
		width:  float32(subsheet.cellWidth),
		height: float32(subsheet.cellHeight),
		viewport: engo.AABB{
			engo.Point{
				X: x / subsheet.width,
				Y: y / subsheet.height,
			},
			engo.Point{
				X: (x + float32(subsheet.cellWidth)) / subsheet.width,
				Y: (y + float32(subsheet.cellHeight)) / subsheet.height,
			},
		},
	}

	return s.cache[index]
}

func (s *Multisheet) Drawable(animationName string, index int) Drawable {
	return s.Cell(animationName, index)
}

func (s *Multisheet) Drawables() []Drawable {
	drawables := make([]Drawable, s.CellCount())
	var i int
	for animationName, subsheet := range s.subsheets {
		drawables[i] = s.Drawable(animationName, i)
		i++
	}

	for i := 0; i < s.CellCount(); i++ {
		drawables[i] = s.Drawable(i)
	}

	return drawables
}

func (s *Subsheet) CellCount() int {
	return int(s.Width()) * int(s.Height())
}

func (s *Multisheet) Cells() []Texture {
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

/*
type Sprite struct {
	Position *Point
	Scale    *Point
	Anchor   *Point
	Rotation float32
	Color    color.Color
	Alpha    float32
	Region   *Region
}

func NewSprite(region *Region, x, y float32) *Sprite {
	return &Sprite{
		Position: &Point{x, y},
		Scale:    &Point{1, 1},
		Anchor:   &Point{0, 0},
		Rotation: 0,
		Color:    color.White,
		Alpha:    1,
		Region:   region,
	}
}
*/
