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
	mateFreq 	float64 //mate frequency: whether integer or a proportion 
	mateThreshold float64 //the distance that two individuals can meet and mate 
	fecundity int //mean offspring numbers: fixed number 
	femaleRate float64 //the percetage of the total offsprings that are female 
}

type Grid struct {
	position      OrderedPair
	costDistance  int 
}

type Landscape struct {
	grid  [][]Grid
	width int
}
