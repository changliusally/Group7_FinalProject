package main

import (
	"encoding/csv"
	"fmt"
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
		panic("Error opening file:", err)
	}
	defer file.Close()

	// Create a new reader
	reader := csv.NewReader(file)

	// Read in header
	if header == true {
		headers, err := reader.Read()
		if err != nil {
			panic("Error reading header:", err)
		}
		fmt.Println(headers)
	}

	// Read in all the records
	records, err := reader.ReadAll()
	if err != nil {
		panic("Error reading all records:", err)
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
	xyParameters := Loadfile(xyfile, true)
	if len(xyParameters) != 6 {
		panic("Error: xyfile column number is not correct")
	}
	individuals := ReadXyfile(xyParameters)
	population.individuals = individuals

	// second column is the int number of Monte-Carlo run
	mcRun, err1 := strconv.Atoi(parameters[1])
	if err1 <= 0 {
		panic("Error: Monte-Carlo run number is non-positive")
	}

	// third column is the int number of the simulation looptime
	looptime, err2 := strconv.Atoi(parameters[2])
	if err2 <= 0 {
		panic("Error: simulation looptime is non-positive")
	}

	// fourth column is the int number of output year
	outputYear, err3 := strconv.Atoi(parameters[3])
	if err3 <= 0 {
		panic("Error: output year is non-positive")
	}

	// fifth column is the filename of cdmatrix
	cdmatrix := parameters[4]
	cdmatData := Loadfile(cdmatrix, true)
	if cdmatData[0][0] != "0" { // check if the first column is 0
		panic("Error: cdmatrix first column is not 0")
	}
	cdmat := ReadCdmatrix(cdmatData)

	// sixth column is the float number of mateFreq
	population.mateFreq, err4 := strconv.ParseFloat(parameters[5], 64)
	if err4 < 0 {
		panic("Error: mateFreq is negative")
	}

	// seventh column is the int number of matureAge
	population.matureAge, err5 := strconv.Atoi(parameters[6])
	if err5 < 0 {
		panic("Error: matureAge is negative")
	}

	// eighth column is the int number of fecundity
	population.fecundity, err6 := strconv.Atoi(parameters[7])
	if err6 < 0 {
		panic("Error: fecundity is negative")
	}

	// ninth column is the float number of femaleRate
	population.femaleRate, err7 := strconv.ParseFloat(parameters[8], 64)
	if err7 < 0 {
		panic("Error: femaleRate is negative")
	}

	// initialize the landscape
	var landscape Landscape
	var model Model

	// tenth column is the string of population model
	model.popModel := parameters[9]

	// eleventh column is the float number of density dependent growth rate	r	‘1.0’	The growth rate used in the density dependent functions above (‘logistic’, ‘richards’, and ‘rickers’).
	model.r_env, err8 := strconv.ParseFloat(parameters[10], 64)
	if err9 < 0 {
		panic("Error: density dependent growth rate is negative")
	}

	// twelfth column is the float number of environmental carrying capacity	K	. The carrying capacity used in the density dependent functions above (‘logistic’, ‘richards’, and ‘rickers’).
	landscape.K_env, err9 := strconv.ParseFloat(parameters[11], 64)
	if err10 <= 0 {
		panic("Error: environmental carrying capacity is non-positive")
	}

	// thirteenth column is the float number of carring capacity for every grid
	landscape.K_grid, err10 := strconv.ParseFloat(parameters[12], 64)
	if err11 < 0 {
		panic("Error: carring capacity for every grid is negative")
	}

	// fourteenth column is the float number of the width of the landscape
	landscape.width, err11 := strconv.Atoi(parameters[13])
	if err12 < 0 {
		panic("Error: width of the landscape is negative")
	}

	// fifteenth to seventeenth column is the float number of the fitness of genotype aa, Aa, AA
	fitness := make([]float64, 3)
	population.fitness[0], err12 := strconv.ParseFloat(parameters[14], 64)
	population.fitness[1], err13 := strconv.ParseFloat(parameters[15], 64)
	population.fitness[2], err14 := strconv.ParseFloat(parameters[16], 64)
	if err12 < 0 || err13 < 0 || err14 < 0 {
		panic("Error: fitness is negative")
	}

	return population, landscape, model, mcRun, looptime, outputYear, cdmat
}

// read xyfile
// input: a slice of slice of string
// output: a slice of individuals
func ReadXyfile(individualData [][]string) individuals {
	// initialize the individuals
	var individuals []Individual

	for i, row := range individualData {
		// initialize the individual
		individuals[i] = Individual{}

		// first column is the int number of X coordinate
		individuals[i].position.x = strconv.Atoi(row[0])

		// second column is the int number of Y coordinate
		individuals[i].position.y = strconv.Atoi(row[1])

		// third column is the int number of the individual id
		individuals[i].id = strconv.Atoi(row[2])

		// fourth column is the int number of the individual age
		individuals[i].age = strconv.Atoi(row[3])

		// fifth column is the int number of individual sex
		individuals[i].sex = strconv.Atoi(row[4])

		// sixth column is the string of individual genetics
		individuals[i].genetics = row[5]
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
