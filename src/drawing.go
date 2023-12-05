package main

import (
	"draw2dimg"
	"image"
	"image/color"
	"math"
)

type Canvas struct {
	gc     *draw2dimg.GraphicContext
	img    image.Image
	width  int // both width and height are in pixels
	height int
}


// function DrawPopulation
func DrawPopulation(population Population, landscape Landscape) Canvas {
	canvasWidth := landscape.width
	gridSize := canvasWidth / 4

	//draw canvas
	c := CreateNewCanvas(canvasWidth, canvasWidth)

	//set grid line
	c.SetStrokeColor(MakeColor(0, 0, 0))
	c.SetLineWidth(1)

	// draw grid line
	for i := 0; i <= 16; i++ {
		x := float64(i) * float64(gridSize)
		c.MoveTo(x, 0)
		c.LineTo(x, float64(canvasWidth))
		c.Stroke()
		c.MoveTo(0, x)
		c.LineTo(float64(canvasWidth), x)
		c.Stroke()
	}

	// traverse through each individual
	for _, individual := range population.individuals {
		// calculate the position
		x := individual.position.x
		y := individual.position.y

		var color color.Color
		if individual.genetics == 0 { // AA
			color = MakeColor(255, 0, 0) // red
		} else if individual.genetics == 1 { // Aa
			color = MakeColor(255, 165, 0) // orange
		} else { // aa
			color = MakeColor(255, 255, 0) // yellow
		}

		//
		c.SetFillColor(color)
		//c.Circle(x+float64(gridSize)/2, y+float64(gridSize)/2, float64(gridSize)/2)
		c.Circle(x, y, float64(landscape.width/800))
		c.Fill()
	}

	return c
}

// Create a new canvas
func CreateNewCanvas(w, h int) Canvas {
	i := image.NewRGBA(image.Rect(0, 0, w, h))
	gc := draw2dimg.NewGraphicContext(i)

	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)
	// fill the background
	gc.Clear()
	gc.SetFillColor(image.Black)

	return Canvas{gc, i, w, h}
}

// Set the line color
func (c *Canvas) SetStrokeColor(col color.Color) {
	c.gc.SetStrokeColor(col)
}

// Set the line width
func (c *Canvas) SetLineWidth(w float64) {
	c.gc.SetLineWidth(w)
}

// Move the current point to (x,y)
func (c *Canvas) MoveTo(x, y float64) {
	c.gc.MoveTo(x, y)
}

// Draw a line from the current point to (x,y), and set the current point to (x,y)
func (c *Canvas) LineTo(x, y float64) {
	c.gc.LineTo(x, y)
}

// Actually draw the lines you've set up with LineTo
func (c *Canvas) Stroke() {
	c.gc.Stroke()
}

// Create a new color
func MakeColor(r, g, b uint8) color.Color {
	return &color.RGBA{r, g, b, 255}
}

// Set the fill color
func (c *Canvas) SetFillColor(col color.Color) {
	c.gc.SetFillColor(col)
}

// Draws an empty circle
// Fill the given circle with the fill color
// Stroke() each time to avoid connected circles
func (c *Canvas) Circle(cx, cy, r float64) {
	c.gc.ArcTo(cx, cy, r, r, 0, -math.Pi*2)
	c.gc.Close()
}

// Fill the area inside the lines you've set up with LineTo, but don't
// draw the lines
func (c *Canvas) Fill() {
	c.gc.Fill()
}
