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
	id int 
	position OrderedPair
	sex int //we set 0 to be male and 1 to be female
	age int
	genetics int //we set 0 to be recessive and 1 to be dominant for a sigle allel
	gen int 
	grid int //which grid this individual is in 
}

type Population struct {
	individuals []Individual
	matureAge   int
	deathRate   float64
	mateFreq 	float64 //mate frequency: whether integer or a proportion 
	mateThreshold float64 //the distance that two individuals can meet and mate 
	fecundity int //mean offspring numbers: fixed number 
	femaleRate float64 //the percetage of the total offsprings that are female 
}

type Grid struct {
	position OrderedPair
	label int 
	kValue int //how many individuals that a grid could process at most 

}

type Landscape struct {
	grid  [][]Grid
	width int
}
