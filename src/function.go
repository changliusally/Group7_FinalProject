package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

// load file
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
func ReadInputParameters(parameters []string) (Population, Landscape, Model, int, int, int, [][]float64) {
	// initialize the population
	var population Population

	// first column is the filename of xyfile(individuals information)
	xyfile := parameters[0]
	random := false
	// check xyfile whether is a number or a filename
	_, err := strconv.Atoi(xyfile)
	if err == nil {
		// xyfile is a number
		random = true
	} else {
		// xyfile is a filename
		xyParameters := Loadfile(xyfile, true)
		if len(xyParameters) != 6 {
			panic("Error: xyfile column number is not correct")
		}
		individuals := ReadXyfile(xyParameters)
		population.individuals = individuals
	}

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
	cdmatData := Loadfile(cdmatrix, true)
	if cdmatData[0][0] != "0" { // check if the first column is 0
		panic("Error: cdmatrix first column is not 0")
	}
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

	// initialize the landscape
	var landscape Landscape
	var model Model

	// tenth column is the string of population model
	model.popModel = parameters[9]

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
		individuals := RandomGenerateIndividuals(xyfile, landscape)
		population.individuals = individuals
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

	return population, landscape, model, mcRun, looptime, outputYear, cdmat
}

// read xyfile
// input: a slice of slice of string
// output: a slice of individuals
func ReadXyfile(individualData [][]string) []Individual {
	// initialize the individuals
	var individuals []Individual

	for i, row := range individualData {
		// initialize the individual
		individuals[i] = Individual{}

		// first column is the int number of X coordinate
		individuals[i].position.x, err1 = strconv.Atoi(row[0])
		if err1 != nil {
			panic(err1)
		}
		if individuals[i].position.x < 0 {
			panic("Error: X coordinate is negative")
		}

		// second column is the int number of Y coordinate
		individuals[i].position.y, err2 = strconv.Atoi(row[1])
		if err2 != nil {
			panic(err2)
		}
		if individuals[i].position.y < 0 {
			panic("Error: Y coordinate is negative")
		}

		// third column is the int number of the individual id
		individuals[i].id, err3 = strconv.Atoi(row[2])
		if err3 != nil {
			panic(err3)
		}
		if individuals[i].id < 0 {
			panic("Error: individual id is negative")
		}

		// fourth column is the int number of the individual age
		individuals[i].age, err4 = strconv.Atoi(row[3])
		if err4 != nil {
			panic(err4)
		}
		if individuals[i].age < 0 {
			panic("Error: individual age is negative")
		}

		// fifth column is the int number of individual sex
		individuals[i].sex, err5 = strconv.Atoi(row[4])
		if err5 != nil {
			panic(err5)
		}
		if individual[i].sex != 0 || individual.sex[i] != 1 {
			panic("Error: sex wrong")
		}

		// sixth column is the string of individual genetics
		genetics := row[5]
		if genetics == "aa" {
			individuals[i].genetics = 0
		} else if genetics == "Aa" {
			individuals[i].genetics = 1
		} else if genetics == "AA" {
			individuals[i].genetics = 2
		} else {
			panic("Error: genetics wrong")
		}
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
		inidividuals[i] = Individual{}
		individuals[i].age = rand.Intn(4) //0,1,2,3
		individual[i].sex = rand.Intn(2)  //0,1
		individual[i].id = i
		individual[i].genetics = rand.Intn(3) //0,1,2
		// rand position

	}

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
				return nil, err
			}
			floatRow = append(floatRow, value)
		}
		matrix = append(matrix, floatRow)
	}

	return matrix
}

// write output to csv
// input: a slice of slice of generation, a int number of output year, a string of output directory
func WriteOutput(generations [][]Generation, outputYear int, outdir string) {
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
		for j := 0; j < len(generations[i]); j++ {
			// if outputyear is 0, then output every generation
			if outputYear == 0 {
				filename := outdir + "/mc" + strconv.Itoa(i) + "/generation" + strconv.Itoa(j) + ".csv"
				WriteCsv(generations[i][j].individuals, filename)
			}
			if j == outputYear {
				filename := outdir + "/mc" + strconv.Itoa(i) + "/generation" + strconv.Itoa(j) + ".csv"
				WriteCsv(generations[i][j].individuals, filename)
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
	for _, individual := range individuals {
		record := []string{
			strconv.Itoa(ind.position.x),
			strconv.Itoa(ind.position.y),
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
