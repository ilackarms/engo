package common

import (
	"log"

	"engo.io/engo"
	"engo.io/gl"
	"github.com/ilackarms/sprite-locator/models"
	"encoding/json"
	"image"
	"github.com/engoengine/math"
)

// sort sprites into rows, starting with top left, and going down row by row
func sortSprites(sprites []models.Sprite) []models.Sprite {
	sorted := []models.Sprite{}
	origin := image.Pt(0, 0)
	topLeft := sprites[0]
	minDist := distance(origin, center(topLeft))
	//find topleft sprite
	for _, sprite := range sprites {
		center := center(sprite)
		dist := distance(origin, center(sprite))
		if dist < minDist {
			topLeft = sprite
			minDist = dist
		}
	}

	//next sprite = closest in X to top left, lowest value of Y
	sprite0 := center(topLeft)
	for _, sprite := range sprites {
		dist := float32(center(sprite).X - sprite0.X)
		if dist < minDist {
			topLeft = sprite
			minDist = dist
		}
	}


	return sorted
}

func center(sprite models.Sprite) image.Point {
	return  image.Point{
		X: (sprite.Min.X + sprite.Max.X)/2,
		Y: (sprite.Min.Y + sprite.Max.Y)/2,
	}
}

func distance(p1, p2 image.Point) float32 {
	return math.Sqrt(math.Pow(p2.X - p1.X, 2)+math.Pow(p2.Y - p1.Y, 2))
}


// Spritesheet is a class that stores a set of tiles from a file, used by tilemaps and animations
type Spritesheet struct {
	texture               *gl.Texture     // The original texture
	width, height         float32         // The dimensions of the total texture
	Sprites []models.Sprite
	cache                 map[int]Texture // The cell cache cells
}

func NewSpritesheetFromTexture(tr *TextureResource, metadata *TextResource) *Spritesheet {
	var spriteMetadata models.Spritesheet
	if err := json.Unmarshal([]byte(metadata.Text), &spriteMetadata); err != nil {
		log.Println("[WARNING] [NewSpritesheetFromFile]: Unmarshalling json from ", metadata.URL(), ": ", err)
		return nil
	}

	return &Spritesheet{texture: tr.Texture,
		width: tr.Width, height: tr.Height,
		Sprites: spriteMetadata.Sprites,
		cache: make(map[int]Texture),
	}
}

// NewSpritesheetFromFile is a simple handler for creating a new spritesheet from a file
// textureName is the name of a texture already preloaded with engo.Files.Add
func NewSpritesheetFromFile(textureName, textName string) *Spritesheet {
	res, err := engo.Files.Resource(textureName)
	if err != nil {
		log.Println("[WARNING] [NewSpritesheetFromFile]: Received error:", err)
		return nil
	}

	img, ok := res.(TextureResource)
	if !ok {
		log.Println("[WARNING] [NewSpritesheetFromFile]: Resource not of type `TextureResource`:", textureName)
		return nil
	}

	res, err = engo.Files.Resource(textName)
	if err != nil {
		log.Println("[WARNING] [NewSpritesheetFromFile]: Received error:", err)
		return nil
	}

	txt, ok := res.(TextResource)
	if !ok {
		log.Println("[WARNING] [NewSpritesheetFromFile]: Resource not of type `TextResource`:", textureName)
		return nil
	}

	return NewSpritesheetFromTexture(&img, &txt)
}

// Cell gets the region at the index i, updates and pulls from cache if need be
func (s *Spritesheet) Cell(index int) Texture {
	if r, ok := s.cache[index]; ok {
		return r
	}

	x0 := float32(s.Sprites[index].Min.X)
	y0 := float32(s.Sprites[index].Min.Y)
	x1 := float32(s.Sprites[index].Max.X)
	y1 := float32(s.Sprites[index].Max.Y)
	s.cache[index] = Texture{
		id: s.texture,
		width: x1 - x0,
		height: y1 - y0,
		viewport: engo.AABB{
			engo.Point{
				X: x0,
				Y: y0,
			},
			engo.Point{
				X: x1,
				Y: y1,
			},
	}}

	return s.cache[index]
}

func (s *Spritesheet) Drawable(index int) Drawable {
	return s.Cell(index)
}

func (s *Spritesheet) Drawables() []Drawable {
	drawables := make([]Drawable, s.CellCount())

	for i := 0; i < s.CellCount(); i++ {
		drawables[i] = s.Drawable(i)
	}

	return drawables
}

func (s *Spritesheet) CellCount() int {
	return len(s.Sprites)
}

func (s *Spritesheet) Cells() []Texture {
	cellsNo := s.CellCount()
	cells := make([]Texture, cellsNo)
	for i := 0; i < cellsNo; i++ {
		cells[i] = s.Cell(i)
	}

	return cells
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
