package main

import (
	
	"time"
	"math/rand"

)

// Dispersal defines the actions that offspring leave their mother and disperse to any other random grid(place) 
// within the allowable range. It includes processes of emigration and immigration and is affected by the 
// differential mortality of each offspring's genotype in the new grid environment.


// DoDispersal is a function than do dispersal of offsprings and population, and removed the redundant offspring when there is no available free grid for them to disperse.
// It will traverse the list of offspring and allocate free grid for it to disperse, and return the death number of offsprings.
func (pop *Population)DoDispersal(land Landscape, offSpring []Individual, probmatrix [][]float64) int{
	death := 0 // record the number of population death when disperse
	dispcount := 0  //the number of dispersed offspring
	offcount := 0  // the number of traversed individuals in population and offsprings
	free := land.K_env-len(pop.individuals)

	// visit every existing individuals in population
	for n := range pop.individuals {
		freegrid := CheckGrid(land, pop.individuals)
		probarray := GetProbArray(pop.individuals, n, probmatrix, freegrid)
		targetGrid := w_choice(freegrid, probarray)
		differentialmortality := DoSelection(pop.individuals[n],targetGrid,pop.fitness)
		differentialmortality_Total := 1 - ((1 - differentialmortality) * (1 - pop.deathRate))
		rand.Seed(time.Now().UnixNano())
		randcheck := rand.Float64()
		if randcheck < differentialmortality_Total {
			death += 1
			continue
		}

		// update the population position
		pop.individuals[n].gridIn = targetGrid
		pop.individuals[n].position.x, pop.individuals[n].position.y = RandomGridxy(targetGrid, land)

	}

	//  makes sure loop stops at carrying capacity (ie, total number of freegrids) or stops at end of offpsring list
	for dispcount < free && offcount < len(offSpring) {

		// visit every offspring
		for i := range offSpring{
			freegrid := CheckGrid(land, pop.individuals)
			probarray := GetProbArray(offSpring, i, probmatrix, freegrid)
			if len(freegrid)!=0 {
				targetGrid := w_choice(freegrid,probarray)
				differentialmortality := DoSelection(offSpring[i],targetGrid,pop.fitness)

				differentialmortality_Total := 1 - ((1 - differentialmortality) * (1 - pop.deathRate))
				
				rand.Seed(time.Now().UnixNano())
				randcheck := rand.Float64()
				if randcheck < differentialmortality_Total {
					offcount += 1
					continue
				}
				
				dispcount += 1
				offcount += 1

				// update population
				offSpring[i].gridIn = targetGrid
				offSpring[i].position.x, offSpring[i].position.y = RandomGridxy(targetGrid, land)
				pop.individuals = append(pop.individuals, offSpring[i])

			} 
		}
	}

	deathcount := offcount - dispcount
	deathcount = deathcount + death
	return deathcount
	
}

// InitializeLand convert the 16 grids into a 4*4 matrix for easy visualization.
// Here, the input width is viewed as dividable by 4.
func InitializeLand(width int) []Grid {
	landGrid := make([]Grid, 0)

	for i := 0; i < 16; i++ { // grid(0-15)
		var grid Grid
		grid.label = i
		index := i+1

		if index % 4 == 1 {
			grid.position.x = 0
			grid.position.y = float64((i/4)*(width/4)) //integer division
		} else if index % 4 == 2 {
			grid.position.x = float64(width/4)
			grid.position.y = float64((i/4)*(width/4))
		} else if index % 4 == 3 {
			grid.position.x = float64(width/2)
			grid.position.y = float64((i/4)*(width/4))
		} else{ //if index % 4 == 0 
			grid.position.x = float64(3*(width/4))
			grid.position.y = float64((i/4)*(width/4))
		}

		landGrid = append(landGrid, grid)
	}

	return landGrid

}

// checkGrid function is to find the available grid that doesn't reach its capacity. 
// It returned the slice of available grids' label.
func CheckGrid(land Landscape, pop []Individual) []int {
	freegrid := make([]int,0)
	for i := 0; i < 16; i++ {
		count := 0
		for _,j := range pop {
			if j.gridIn == i {
				count += 1
			}
		}
		if count < land.K_grid {
			freegrid = append(freegrid, i)
		}
	}

	return freegrid
}
// freegrid are filetered grids that still have available capacity for immigrants
// freegrid = [xycdmatid location of free grid spot in random order] = index of free grid.



// w_choice is a function to do weigted random selection from a list
// this is mainly used for the random selection of the new grid the offspring is going to move into.
func w_choice(freegrid []int, prob []float64) int{
	wtotal := 0.0
	for x := range freegrid {
		wtotal += prob[x]
	}

	// generate a random number between 0 and wtotal
	n := rand.Float64() * (wtotal)
	var index int
	
	for x := range freegrid {
		if n < prob[x]{
			index = x
			break
		} else {
			n = n-prob[x]
		}
	}
	return index 
}



// getProbArray is a function to retreive the cost distance information of all possible free grid for a specific offspring.
// This output is typically input for the w_choice function.
func GetProbArray(offSpring []Individual, index int, probmatrix [][]float64, freegrid []int) []float64{
	currOff := offSpring[index]
	var probarray []float64 
	
    for _,Gindex := range freegrid {
         probarray = append(probarray, probmatrix[currOff.gridIn][Gindex]) // list of cells [currentGrid, freeGrid]
    }

	return probarray
}


// DoHindexSelection is a function to mimick the situation that offspring's motality will be affected when moving to a new environment.
// It calculates a specific offspring's differential mortality based on its genotype and the environment resistance.
// Hindex is used to describe the degree of how close the genotype is to homozygous dominant AA. (AA = 1, Aa = 0.5, aa = 0)
// xvars are slice of environmental factors ranging from -1 to 1 for each grid.
func DoHindexSelection(dispOff Individual, chosenGrid int, pars []float64,xvars []float64) float64 {
	slope_min := pars[0]
	slope_max := pars[1]
	int_min := pars[2]
	int_max := pars[3]
	X_min := pars[4]
	X_max := pars[5]
	X_val := xvars[chosenGrid]
	// use fitness first for a trial first, AA100, aa50...then we need a fitness matrix, row is grid index, col is genotype.
	
	// Calculate slope
	m := ((slope_min - slope_max) / (X_min - X_max)) * X_val - X_min * ((slope_min - slope_max) / (X_min - X_max)) + slope_min
	
	// Calculate intercept
	b := ((int_max - int_min) / (X_min - X_max)) * X_val - X_min * ((int_max - int_min) / (X_min - X_max)) + int_max
		
	Hindex := 0.5*float64(dispOff.genetics)
	Fitness := m*Hindex + b
	differentialmortality := 1.0 - Fitness
	return differentialmortality
}

// use fitness first for a trial first, AA100, aa50...then we need a fitness matrix, row is grid index, col is genotype.
func DoSelection(dispOff Individual, chosenGrid int, fitness []float64) float64{
	Fitness := fitness[dispOff.genetics]
	differentialmortality := 1.0 - Fitness
	return differentialmortality
}

// randomGridxy takes in the grid index, and will randomly generate a location within its range.
func RandomGridxy(target int, land Landscape) (float64,float64) {
	rand.Seed(time.Now().UnixNano())
	x_min := land.grid[target].position.x
	y_min := land.grid[target].position.y

	x := rand.Float64()*float64(land.width/4) + x_min
	y := rand.Float64()*float64(land.width/4) + y_min

	return x,y

}


// calProb function will convert the cost distance matrix of n*n individuals into probability matrix.
// We have different methods to do this conversion, include linear, nearest neighbor, random mixing and so on. (movement function)
func CalProb(method string, cdmatrix [][]float64) [][]float64 {
	
	max := FindMax(cdmatrix)
	min := FindMin(cdmatrix)
	probMatrix := CopyMatrix(cdmatrix)

	if min < 0 || max < 0 {
		panic("The cost distance cannot have value smaller than 0.")
	}

	if method == "linear" {
		for i:= 0; i < len(cdmatrix); i++{
			for j := 0; j < len(cdmatrix[0]); j++ {
				probMatrix[i][j] = 1-(1.0/max)*cdmatrix[i][j]
			}
		}
	}

	return probMatrix
}



// copyMatrix function takes input a [][]float64, and will return a copy of this matrix.
func CopyMatrix(matrix [][]float64) [][]float64 {
	// Get the dimensions of the matrix
	nrows := len(matrix)
	ncols := len(matrix[0])

	// Create a new matrix with the same dimensions
	copy := make([][]float64, nrows)
	for i := range copy {
		copy[i] = make([]float64, ncols)
	}

	// Copy the values from the original matrix to the new matrix
	for i := 0; i < nrows; i++ {
		for j := 0; j < ncols; j++ {
			copy[i][j] = matrix[i][j]
		}
	}

	return copy
}



//FindMin takes a distance matrix as input and returns the maximum value of the distance matrix
func FindMin(matrix [][]float64) float64 {
	// Assuming the matrix is not empty
	minValue := matrix[0][0]

	for _,row := range matrix {
		for _,value := range row {
			if value < minValue {
				minValue = value
			}
		}
	}

	return minValue
}

//FindMax takes a distance matrix as input and returns the maximum value of the distance matrix
func FindMax(matrix [][]float64) float64 {
	// Assuming the matrix is not empty
	maxValue := matrix[0][0]

	for _,row := range matrix {
		for _,value := range row {
			if value > maxValue {
				maxValue = value
			}
		}
	}

	return maxValue
}