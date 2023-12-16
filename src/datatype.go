package main

//Global symbols

// File absolute paths for importing functions
var SRC_PATH = "../src"

// struct
type OrderedPair struct {
	x float64
	y float64
}

type Individual struct {
	id       int
	position OrderedPair
	sex      int //we set 0 to be male and 1 to be female
	age      int
	genetics int //we set 0 to be recessive and 1 to be dominant for a sigle allel
	//gen int, might be useful in visualization
	gridIn int //which grid this individual is in
}

type Population struct {
	individuals     []Individual
	matureAge       int
	deathRate       float64
	mateFreq        float64   //mate frequency: whether integer or a proportion
	fecundity       int       //mean offspring numbers: fixed number
	femaleRate      float64   //the percetage of the total offsprings that are female
	fitness         []float64 //fitness[0]:aa, fitness[1]:Aa, fitness[1]:AA
	dispersalMethod string    //linear or inverse2
	offspringMethod string    //constant or poisson or random or normal
}

type Generation struct {
	population []Population
}

type Grid struct {
	position OrderedPair
	label    int
}

type Landscape struct {
	grid   []Grid
	width  int
	K_grid int //how many individuals that a grid could process at most
	K_env  int // environmental carrying capacity
}

type Model struct {
	popModel string
	r_env    float64
}
