package main

//Global symbols

//File absolute paths for importing functions
var SRC_PATH = "../src"

// struct
type OrderedPair struct {
	x float64
	y float64
}

type Individual struct {
	id       int
	age      int
	sex      int //0: male 1:famale
	position OrderedPair
	genetics int //single loci, 0:aa 1:Aa 2:AA
	gen      int //which generation
}

type Population struct {
	individuals []Individual
	mutateAge   int
	birthRate   float64
	deathRate   float64
}

type Grid struct {
	position      OrderedPair
	landscapeType int //0:land 1:water
}

type Landscape struct {
	grid  [][]Grid
	width int
}
