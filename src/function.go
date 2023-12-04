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
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a new reader
	reader := csv.NewReader(file)

	// Read in header
	if header == true {
		headers, err := reader.Read()
		if err != nil {
			fmt.Println("Error reading header:", err)
		}
	}

	// Read in all the records
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading all records:", err)
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
	mcRun := strconv.Atoi(parameters[1])

	// third column is the int number of the simulation looptime
	looptime := strconv.Atoi(parameters[2])

	// fourth column is the int number of output year
	outputYear := strconv.Atoi(parameters[3])

	// fifth column is the filename of cdmatrix
	cdmatrix := parameters[4]
	cdmatData := Loadfile(cdmatrix, true)
	if cdmatData[0][0] != "0" { // check if the first column is 0
		panic("Error: cdmatrix first column is not 0")
	}
	cdmat := ReadCdmatrix(cdmatData)

	// sixth column is the float number of mateFreq
	population.mateFreq := strconv.ParseFloat(parameters[5], 64)

	// seventh column is the int number of matureAge
	population.matureAge := strconv.Atoi(parameters[6])

	// eighth column is the int number of fecundity
	population.fecundity := strconv.Atoi(parameters[7])

	// ninth column is the float number of femaleRate
	population.femaleRate := strconv.ParseFloat(parameters[8], 64)

	// initialize the landscape
	var landscape Landscape
	var model Model

	// tenth column is the string of population model
	model.popModel := parameters[9]

	// eleventh column is the float number of density dependent growth rate	r	‘1.0’	The growth rate used in the density dependent functions above (‘logistic’, ‘richards’, and ‘rickers’).
	model.r_env := strconv.ParseFloat(parameters[10], 64)

	// twelfth column is the float number of environmental carrying capacity	K	. The carrying capacity used in the density dependent functions above (‘logistic’, ‘richards’, and ‘rickers’).
	landscape.K_env := strconv.ParseFloat(parameters[11], 64)

	// thirteenth column is the float number of carring capacity for every grid
	landscape.K_grid := strconv.ParseFloat(parameters[12], 64)

	// fourteenth column is the float number of the width of the landscape
	landscape.width := strconv.Atoi(parameters[13])

	// fifteenth to seventeenth column is the float number of the fitness of genotype aa, Aa, AA
	fitness := make([]float64, 3)
	population.fitness[0] := strconv.ParseFloat(parameters[14], 64)
	population.fitness[1] := strconv.ParseFloat(parameters[15], 64)
	population.fitness[2] := strconv.ParseFloat(parameters[16], 64)

	return population, landscape, model, mcRun, looptime, outputYear, cdmat
}

// read xyfile
// input: a slice of slice of string
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
