package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"math"
	"math/rand"
)

// adultDeath funciton takes in the constant motality and the population and will return the updated population.
// Constant mortality is applied to each age class, and track total number of deaths for each age and sex class.
func adultDeath(pop Population)  {
	total := len(pop.individuals)
	deathNo := math.Floor(total * pop.deathRate)

	// remove deathNo individuals randomly, immediately modify the original population
	// should include a copyPOP function?
	for i := range death {
		n := rand.Intn(total)
		if n != 0 || n != total-1 {
			pop = append(pop[:n-1],pop[n+1:]...)
		} else if n == 0 {
			pop = pop[1:]
		} else {
			pop = pop[:total-2]
		}
		
		total -= 1
	}

}
// here to make sure population and landscape ([]Grid)

// freegrid haven't be written




// Dispersal defines the actions that offspring leave their mother and disperse to any other random grid(place) 
// within the allowable range. It includes processes of emigration and immigration and is affected by the 
// differential mortality of each offspring's genotype in the new grid environment.

// w_choice is a function to do weigted random selection from a list
// this is mainly used for the random selection of the new grid the offspring is going to move into.
func w_choice(freegrid []Grid) Grid{
	wtotal := 0.0
	for x := range freegrid {
		wtotal += x.costDistance
	}

	// generate a random number between 0 and wtotal
	n := rand.Intn(wtotal)
	var item Grid
	for x := range freegrid {
		if n < x.costDistance{
			item = x
			break
		} else {
			n = n-x.costDistance
		}
	}
	return item //index is better maybe
}
// freegrid are filetered grids that within the dispersal threshold, and in descedning order, so that it will consider 
// grid that has larger probability first. Paper has it in random order.
// freegrid = [xycdmatid location of free grid spot in random order]
// freegrid = index of free grid.



// getProbArray is a function to retreive the cost distance information of all possible free grid for a specific offspring.
// This output is typically input for the w_choice function.
func getProbArray(offSpring, index, xycdmatrix, freegrid) []grid{
	currOff := offSpring[index]
	probarray = probMatrix[currOff.grid][freegrid] // list of cells [currentGrid, freeGrid]
	return probarray
}
// individual should not have position field, has grid would be rather good.


// DoHindexSelection is a function to mimick the situation that offspring's motality will be affected when moving to a new environment.
// It calculates a specific offspring's differential mortality based on its genotype and the environment resistance.
// Hindex is used to describe the degree of how close the genotype is to homozygous dominant AA. (AA = 1, Aa = 0.5, aa = 0)
func DoHindexSelection(dispOff, freegrid, chosenGrid, index, cdevolveans,xvars) float64 {
	
	// use fitness first for a trial first, AA100, aa50...then we need a fitness matrix, row is grid index, col is genotype.
	
	differentialmortality = 1.0 - Fitness
	return differentialmortality
}


type dispInfo struct {
	numGen int  // maybe not needed
	dispOFF Individual
	targetGrid int // the grid index
}
// DoEmigration is a function than do dispersal of offsprings, and removed the redundant offspring when there is no available free grid for them to disperse.
// It will traverse the list of offspring and allocate free grid for it to disperse, and record the successful dispersals.
func DoEmigration(offSpring, freegrid, xycdmatrix, gen, cdevolveans, xvar) []dispInfo {
	disp := make([]dispInfo,0)
	
	dispcount := 0
	offcount := 0

	//  makes sure loop stops at carrying capacity (ie, total number of freegrids) or stops at end of offpsring list
	for dispcount < len(freegrid) && offcount < len(offSpring) {

		// visit every offspring
		for i := range offSpring{
			probarray := getProbArray(offSpring, i, xycdmatrix, freegrid)
			if len(freegrid)!=0 {
				targetGrid = w_choice(probarray)
				differentialmortality := 1-Fitness[offSpring[i].genetics][targetGrid]

				differentialmortality_Total = 1 - ((1 - differentialmortality) * (1 - pop.mortality))
				
				rand.Seed(time.Now().UnixNano())
				randcheck = rand.Float64()
				if randcheck < differentialmortality_Total {
					offcount += 1
					continue//not sure?
				}
				disp[i].dispOFF = offSpring[i]
				disp[i].targetGrid = targetGrid
				dispcount += 1
				offcount += 1

				// update population
				offSpring[i].Grid = targetGrid
				population = append(population, offSpring[i])

				del(freegrid[targetGrid])

			} 
		}
	}
	return disp
}



// ReadcdMatrix function will convert the cost distance matrix of n*n individuals into probability matrix.
// We have different methods to do this conversion, include linear, nearest neighbor, random mixing and so on. (movement function)
func ReadcdMatrix(filename, method string) [][]int {
	// can write a copyMatrix function
	cdmatrix := readCSV(filename)
	threshold := findMax(cdmatrix)
	min := findMin(cdmatrix)
	probMatrix := copyMatrix(cdmatrix)

	if min < 0 || max < 0 {
		panic("The cost distance cannot have value smaller than 0.")
	}

	if method == "linear" {
		for i:= 0; i < len(cdmatrix); i++{
			for j := 0; j < len(cdmatrix[0]); j++ {
				probMatrix[i][j] = 1-(1.0/threshold)*cdmatrix[i][j]
			}
		}
	}

	return probMatrix
}



// copyMatrix function takes input a [][]float64, and will return a copy of this matrix.
func copyMatrix(matrix [][]float64) [][]float64 {

}


// readCSV function helps read the provided cdmatrix in csv file form, and return it in [][]float64 data type.
func readCSV(fileName string) ([][]float64) {
	// Open the CSV file
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all the records from CSV
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

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




func findMin(matrix [][]float64) float64 {
	// Assuming the matrix is not empty
	minValue := math.MaxFloat64

	for _, row := range matrix {
		for _, value := range row {
			if value < minValue {
				minValue = value
			}
		}
	}

	return minValue
}