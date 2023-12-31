package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// load file from csv and return the content
// input: a csv file
// output: a slice of slice of string recording the data
func Loadfile(filename string, header bool) [][]string {

	// Open the file
	file, err := os.Open(filename) //
	if err != nil {
		panic("Error opening file:")
	}
	defer file.Close()

	// Create a new reader
	reader := csv.NewReader(file)

	// Read in header
	if header == true {
		headers, err := reader.Read()
		if err != nil {
			panic("Error reading header:")
		}
		fmt.Println(headers)
	}

	// Read in all the records
	records, err := reader.ReadAll()
	if err != nil {
		panic("Error reading all records:")
	}

	return records
}

// read parameters from input file
// input: a slice of string
func ReadInputParameters(parameters []string, datadir string) (Population, Landscape, Model, int, int, int, [][]float64) {
	// initialize the population
	var population Population

	var individuals []Individual

	// first column is the filename of xyfile(individuals information) or a number of individuals
	// if it is a number, then random generate individuals
	xyfile := parameters[0]
	random := false
	// check xyfile whether is a number or a filename
	xyVal, err := strconv.Atoi(xyfile)
	if err == nil {
		// xyfile is a number
		random = true
	} else {
		// xyfile is a filename
		xyFile := datadir + xyfile
		xyParameters := Loadfile(xyFile, true)
		if len(xyParameters[0]) != 6 {
			panic("Error: xyfile column number is not correct")
		}
		individuals = ReadXyfile(xyParameters)
		//fmt.Println("before", individuals)
		population.individuals = individuals
		//fmt.Println("after", population.individuals)
	}

	// fmt.Println("hello0", individuals)
	// second column is the int number of Monte-Carlo run
	mcRun, err1 := strconv.Atoi(parameters[1])
	if err1 != nil {
		panic(err1)
	}
	if mcRun <= 0 {
		panic("Error: Monte-Carlo run number is non-positive")
	}

	// third column is the int number of the simulation looptime
	looptime, err2 := strconv.Atoi(parameters[2])
	if err2 != nil {
		panic(err2)
	}
	if looptime <= 0 {
		panic("Error: simulation looptime is non-positive")
	}

	// fourth column is the int number of output year
	outputYear, err3 := strconv.Atoi(parameters[3])
	if err3 != nil {
		panic(err3)
	}
	if outputYear < 0 {
		panic("Error: output year is negative")
	}

	// fifth column is the filename of cdmatrix
	cdmatrix := parameters[4]
	cdPath := datadir + cdmatrix
	cdmatData := Loadfile(cdPath, false)

	cdmat := ReadCdmatrix(cdmatData)

	// sixth column is the float number of mateFreq
	mateFreq, err4 := strconv.ParseFloat(parameters[5], 64)
	if err4 != nil {
		panic(err4)
	}
	if mateFreq < 0 {
		panic("Error: mateFreq is negative")
	}
	population.mateFreq = mateFreq

	// seventh column is the int number of matureAge
	matureAge, err5 := strconv.Atoi(parameters[6])
	if err5 != nil {
		panic(err5)
	}
	if matureAge < 0 {
		panic("Error: matureAge is negative")
	}
	population.matureAge = matureAge

	// eighth column is the int number of fecundity
	fecundity, err6 := strconv.Atoi(parameters[7])
	if err6 != nil {
		panic(err6)
	}
	if fecundity < 0 {
		panic("Error: fecundity is negative")
	}
	population.fecundity = fecundity

	// ninth column is the float number of femaleRate
	femaleRate, err7 := strconv.ParseFloat(parameters[8], 64)
	if err7 != nil {
		panic(err7)
	}
	if femaleRate < 0 {
		panic("Error: femaleRate is negative")
	}
	population.femaleRate = femaleRate

	// tenth column is the string of disperal method
	dispersalMethod := parameters[9]
	// linear, inverse2
	if dispersalMethod != "linear" && dispersalMethod != "inverse2" {
		panic("Error: dispersal method is wrong")
	}
	population.dispersalMethod = dispersalMethod

	// initialize the landscape
	var landscape Landscape
	var model Model

	// eleventh column is the float number of density dependent growth rate	r	‘1.0’	The growth rate used in the density dependent functions above (‘logistic’, ‘richards’, and ‘rickers’).
	r_env, err8 := strconv.ParseFloat(parameters[10], 64)
	if err8 != nil {
		panic(err8)
	}
	if r_env < 0 {
		panic("Error: density dependent growth rate is negative")
	}
	model.r_env = r_env

	// twelfth column is the int number of environmental carrying capacity	K	. The carrying capacity used in the density dependent functions above (‘logistic’, ‘richards’, and ‘rickers’).
	K_env, err9 := strconv.Atoi(parameters[11])
	if err9 != nil {
		panic(err9)
	}
	if K_env < 0 {
		panic("Error: environmental carrying capacity is non-positive")
	}
	landscape.K_env = K_env

	// thirteenth column is the int number of carring capacity for every grid
	K_grid, err10 := strconv.Atoi(parameters[12])
	if err10 != nil {
		panic(err10)
	}
	if K_grid < 0 {
		panic("Error: carring capacity for every grid is negative")
	}
	landscape.K_grid = K_grid

	// fourteenth column is the int number of the width of the landscape
	width, err11 := strconv.Atoi(parameters[13])
	if err11 != nil {
		panic(err11)
	}
	if width < 0 {
		panic("Error: width of the landscape is negative")
	}
	landscape.width = width

	landGrid := InitializeLand(landscape.width)
	landscape.grid = landGrid

	if random == true {
		// generate the individuals
		individuals = RandomGenerateIndividuals(xyVal, landscape)
		population.individuals = individuals
	} else {
		//fmt.Println("hello1", individuals)
		individuals = FindGrid(landscape, individuals)
		//fmt.Println("hello2", individuals)
		population.individuals = individuals
		//fmt.Println("population", population.individuals)
	}

	// fifteenth to seventeenth column is the float number of the fitness of genotype aa, Aa, AA
	fitness := make([]float64, 3)
	fitness_aa, err12 := strconv.ParseFloat(parameters[14], 64)
	if err12 != nil {
		panic(err12)
	}
	fitness_Aa, err13 := strconv.ParseFloat(parameters[15], 64)
	if err13 != nil {
		panic(err13)
	}
	fitness_AA, err14 := strconv.ParseFloat(parameters[16], 64)
	if err14 != nil {
		panic(err14)
	}
	fitness[0] = fitness_aa
	fitness[1] = fitness_Aa
	fitness[2] = fitness_AA
	if fitness[0] < 0 || fitness[1] < 0 || fitness[2] < 0 {
		panic("Error: fitness is negative")
	}
	population.fitness = fitness

	// eighteenth column is the string of offspring method
	offspringMethod := parameters[17]
	// poisson, constant, random, normal
	if offspringMethod != "poisson" && offspringMethod != "constant" && offspringMethod != "random" && offspringMethod != "normal" {
		panic("Error: offspring method is wrong")
	}

	population.offspringMethod = offspringMethod

	return population, landscape, model, mcRun, looptime, outputYear, cdmat
}

// read xyfile
// input: a slice of slice of string
// output: a slice of individuals
func ReadXyfile(individualData [][]string) []Individual {
	// initialize the individuals
	individuals := make([]Individual, len(individualData))
	//fmt.Println(individualData)

	for i := range individualData {
		row := individualData[i]

		// initialize the individual
		var individual Individual

		// initialize the position
		var position OrderedPair

		// first column is the int number of X coordinate
		x, err1 := strconv.ParseFloat(row[0], 64)
		if err1 != nil {
			panic(err1)
		}
		if x < 0 {
			panic("Error: X coordinate is negative")
		}
		position.x = x
		individual.position.x = position.x

		// second column is the int number of Y coordinate
		y, err2 := strconv.ParseFloat(row[1], 64)
		if err2 != nil {
			panic(err2)
		}
		if y < 0 {
			panic("Error: Y coordinate is negative")
		}
		position.y = y
		individual.position.y = position.y

		// third column is the int number of the individual id
		id, err3 := strconv.Atoi(row[2])
		if err3 != nil {
			panic(err3)
		}
		if id < 0 {
			panic("Error: individual id is negative")
		}
		individual.id = id

		// fourth column is the int number of the individual age
		age, err4 := strconv.Atoi(row[3])
		if err4 != nil {
			panic(err4)
		}
		if age < 0 {
			panic("Error: individual age is negative")
		}
		individual.age = age

		// fifth column is the int number of individual sex
		sex, err5 := strconv.Atoi(row[4])
		if err5 != nil {
			panic(err5)
		}
		if sex != 0 && sex != 1 {
			panic("Error: sex wrong")
		}

		individual.sex = sex

		// sixth column is the string of individual genetics
		genetics := row[5]
		if genetics == "aa" {
			individual.genetics = 0
		} else if genetics == "Aa" {
			individual.genetics = 1
		} else if genetics == "AA" {
			individual.genetics = 2
		} else {
			panic("Error: genetics wrong")
		}

		individuals[i] = individual
	}

	return individuals
}

// random generate individuals
// output: a slice of individuals
func RandomGenerateIndividuals(num int, landscape Landscape) []Individual {
	// initialize the individuals
	var individuals []Individual

	// generate the position, age, sex ,id, genetics for every individual
	for i := 0; i < num; i++ {
		var individual Individual
		individual.age = rand.Intn(4) //0,1,2,3
		individual.sex = rand.Intn(2) //0,1
		individual.id = i
		//individual.genetics = rand.Intn(3) //0,1,2
		randGene := rand.Intn(4)
		if randGene == 0 {
			individual.genetics = 0
		} else if randGene == 1 || randGene == 2 {
			individual.genetics = 1
		} else {
			individual.genetics = 2
		}
		// rand position
		individual.gridIn = rand.Intn(16)
		var position OrderedPair
		position.x, position.y = RandomGridxy(individual.gridIn, landscape)
		individual.position = position
		individuals = append(individuals, individual)
	}
	//fmt.Println(individuals)

	return individuals
}

// read cdmatrix
// input: a slice of slice of string
// output: a slice of slice of float64
func ReadCdmatrix(records [][]string) [][]float64 {
	// Convert string records to [][]float64
	var matrix [][]float64
	for _, row := range records {
		var floatRow []float64
		for _, cell := range row {
			value, err := strconv.ParseFloat(cell, 64)
			if err != nil {
				panic(err)
			}
			floatRow = append(floatRow, value)
		}
		matrix = append(matrix, floatRow)
	}

	return matrix
}

// write output to csv
// input: a slice of slice of generation, a int number of output year, a string of output directory
func WriteOutput(generations []Generation, outputYear int, outdir string, landscape Landscape) {
	// make output directory
	err := os.MkdirAll(outdir, 0777)
	if err != nil {
		panic(err)
	}

	// make every monte carlo run directory
	for i := 0; i < len(generations); i++ {
		err := os.MkdirAll(outdir+"/mc"+strconv.Itoa(i), 0777)
		if err != nil {
			panic(err)
		}
		// write every generation to csv
		for m := 0; m < len(generations[i].population); m++ {
			// if outputyear is 0, then output every generation
			if outputYear == 0 {
				filename := outdir + "/mc" + strconv.Itoa(i) + "/generation" + strconv.Itoa(m) + ".csv"
				WriteCsv(generations[i].population[m].individuals, filename)
				// this is for the generation of pdf rather than gif
				// generate the png for every 10 generations
				if m%10 == 0 {
					photo := outdir + "/mc" + strconv.Itoa(i) + "/generation" + strconv.Itoa(m) + ".png"
					generations[i].population[m].DrawPopulation2(landscape, photo)
				}
			}
			if m == outputYear {
				filename := outdir + "/mc" + strconv.Itoa(i) + "/generation" + strconv.Itoa(m) + ".csv"
				WriteCsv(generations[i].population[m].individuals, filename)
			}
		}

	}

}

// write csv
// input: a slice of Individual, a string of filename
func WriteCsv(individuals []Individual, filename string) {
	// header: XCOORD,YCOORD,ID,age,sex,genetics,gridIn
	// open the csv file
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// the csv writer use "," as the default delimiter
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write the header
	header := []string{"XCOORD", "YCOORD", "ID", "age", "sex", "genetics", "gridIn"}
	err1 := writer.Write(header)
	if err1 != nil {
		panic(err1)
	}

	// write the data
	for _, ind := range individuals {
		record := []string{
			strconv.FormatFloat(ind.position.x, 'f', 2, 64),
			strconv.FormatFloat(ind.position.y, 'f', 2, 64),
			strconv.Itoa(ind.id),
			strconv.Itoa(ind.age),
			strconv.Itoa(ind.sex),
			strconv.Itoa(ind.genetics),
			strconv.Itoa(ind.gridIn),
		}
		err2 := writer.Write(record)
		if err2 != nil {
			panic(err2)
		}
	}

}

// initializeCDmat generates a 16x16 matrix with random values limited by maxLimit.
// input: a float64 number
// output: a slice of slice of float64
func initializeCDmat(maxLimit float64) [][]float64 {
	rand.Seed(time.Now().UnixNano())

	matrix := make([][]float64, 16)
	for i := range matrix {
		matrix[i] = make([]float64, 16)
		for j := range matrix[i] {
			// Generate a random float64 value between 0 and maxLimit
			matrix[i][j] = rand.Float64() * maxLimit
		}
	}

	return matrix
}

// FindGrid takes a lacdscape and a slice of individuals as input and returns the same
// individual slice with each individual's gridIn updated based on the their position
// in the landscape
func FindGrid(landscape Landscape, individuals []Individual) []Individual {
	width := landscape.width

	for i := range individuals {
		position := individuals[i].position
		//check if the position is valid
		if position.x > float64(width) || position.x < 0 {
			panic("Invalid x coordinate")
		}
		if position.y > float64(width) || position.y < 0 {
			panic("Invalid y coordinate")
		}

		//find the grid number
		if position.x >= 0.0 && position.x < float64(width)/4.0 {
			if position.y >= 0.0 && position.y < float64(width)/4.0 {
				individuals[i].gridIn = 0
			} else if position.y >= float64(width)/4.0 && position.y < float64(width)/2.0 {
				individuals[i].gridIn = 4
			} else if position.y >= float64(width)/2.0 && position.y < float64(3*width)/4.0 {
				individuals[i].gridIn = 8
			} else {
				individuals[i].gridIn = 12
			}
		} else if position.x >= float64(width)/4.0 && position.x < float64(width)/2.0 {
			if position.y >= 0 && position.y < float64(width)/4.0 {
				individuals[i].gridIn = 1
			} else if position.y >= float64(width)/4.0 && position.y < float64(width)/2.0 {
				individuals[i].gridIn = 5
			} else if position.y >= float64(width)/2.0 && position.y < float64(3*width)/4.0 {
				individuals[i].gridIn = 9
			} else {
				individuals[i].gridIn = 13
			}
		} else if position.x >= float64(width)/2.0 && position.x < float64(3*width)/4.0 {
			if position.y >= 0 && position.y < float64(width)/4.0 {
				individuals[i].gridIn = 2
			} else if position.y >= float64(width)/4.0 && position.y < float64(width)/2.0 {
				individuals[i].gridIn = 6
			} else if position.y >= float64(width)/2.0 && position.y < float64(3*width)/4.0 {
				individuals[i].gridIn = 10
			} else {
				individuals[i].gridIn = 14
			}
		} else {
			if position.y >= 0 && position.y < float64(width)/4.0 {
				individuals[i].gridIn = 3
			} else if position.y >= float64(width)/4.0 && position.y < float64(width)/2.0 {
				individuals[i].gridIn = 7
			} else if position.y >= float64(width)/2.0 && position.y < float64(3*width)/4.0 {
				individuals[i].gridIn = 11
			} else {
				individuals[i].gridIn = 15
			}
		}

	}
	return individuals

}

// write summary to csv
// input: a slice of generation, a string of output directory
func WriteSummary(generations []Generation, outdir string) {

	// open every monte carlo run directory
	for i := 0; i < len(generations); i++ {
		csvOutdir := outdir + "/mc" + strconv.Itoa(i)
		filename := csvOutdir + "_summary.csv"
		// header: Year, TotalIndividuals, TotalFemale, TotalMale, AA, Aa, aa, AAFreq, AaFreq, aaFreq, AFreq, aFreq
		// open the csv file
		file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		// the csv writer use "," as the default delimiter
		writer := csv.NewWriter(file)
		defer writer.Flush()

		// write the header
		header := []string{"Year", "TotalIndividuals", "TotalFemale", "TotalMale", "AA", "Aa", "aa", "AAFreq", "AaFreq", "aaFreq", "AFreq", "aFreq"}
		err1 := writer.Write(header)
		if err1 != nil {
			panic(err1)
		}

		// write every generation to csv
		for m := 0; m < len(generations[i].population); m++ {
			// write the data
			record := []string{
				strconv.Itoa(m),
				strconv.Itoa(len(generations[i].population[m].individuals)),
				strconv.Itoa(CalTotalFemale(generations[i].population[m])),
				strconv.Itoa(CalTotalMale(generations[i].population[m])),
				strconv.Itoa(CalculateAA(generations[i].population[m])),
				strconv.Itoa(CalculateAa(generations[i].population[m])),
				strconv.Itoa(Calculateaa(generations[i].population[m])),
				strconv.FormatFloat(CalculateFrequency(generations[i].population[m], CalculateAA(generations[i].population[m])), 'f', 2, 64),
				strconv.FormatFloat(CalculateFrequency(generations[i].population[m], CalculateAa(generations[i].population[m])), 'f', 2, 64),
				strconv.FormatFloat(CalculateFrequency(generations[i].population[m], Calculateaa(generations[i].population[m])), 'f', 2, 64),
				strconv.FormatFloat(CalAFreq(CalculateFrequency(generations[i].population[m], CalculateAA(generations[i].population[m])), CalculateFrequency(generations[i].population[m], CalculateAa(generations[i].population[m]))), 'f', 2, 64),
				strconv.FormatFloat(CalAFreq(CalculateFrequency(generations[i].population[m], CalculateAa(generations[i].population[m])), CalculateFrequency(generations[i].population[m], Calculateaa(generations[i].population[m]))), 'f', 2, 64),
			}
			err2 := writer.Write(record)
			if err2 != nil {
				panic(err2)
			}
		}
	}
}

// CalTotalFemale takes a population as input and returns the number of female individuals
// in this population
func CalTotalFemale(pop Population) int {
	count := 0
	for i := range pop.individuals {
		if pop.individuals[i].sex == 1 {
			count++
		}
	}
	return count
}

// CalTotalMale takes a population as input and returns the number of male individuals
// in this population
func CalTotalMale(pop Population) int {
	count := 0
	for i := range pop.individuals {
		if pop.individuals[i].sex == 0 {
			count++
		}
	}
	return count
}

// CalculateAA takes a population as input and returns the number of individuals
// in the population whose genetics is AA
func CalculateAA(pop Population) int {
	count := 0
	for i := range pop.individuals {
		if pop.individuals[i].genetics == 2 {
			count++
		}
	}
	return count
}

// CalculateAa takes a population as input and returns the number of individuals
// in the population whose genetics is Aa
func CalculateAa(pop Population) int {
	count := 0
	for i := range pop.individuals {
		if pop.individuals[i].genetics == 1 {
			count++
		}
	}
	return count
}

// Calculateaa takes a population as input and returns the number of individuals
// in the population whose genetics is aa
func Calculateaa(pop Population) int {
	count := 0
	for i := range pop.individuals {
		if pop.individuals[i].genetics == 0 {
			count++
		}
	}
	return count
}

// CalculateFrequency takes a population and the number of individuals in the population
// whose genetics is what we want as input and returns the proportion of individuals
// in the population of wanted
func CalculateFrequency(pop Population, num int) float64 {
	n := len(pop.individuals)
	freq := float64(num) / float64(n)
	return freq
}

// CalAFreq takes the proportion of AA and Aa individuals as input
// and returns the A allel frequency of this population
func CalAFreq(freqAA, freqAa float64) float64 {
	freq := freqAA + freqAa/2.0
	return freq
}
