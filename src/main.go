package main

import (
	"runtime"
	"time"
	"os"
	"fmt"
)

func main() {

	// set up the timer
	start_time := time.Now()
	foldertime := start_time.strftime("%Y%m%d_%H%M%S")

	// read agv from command line
	// the format: main.exe C:/.../inputfolder inputvars.csv output_folder

	fmt.Println(os.Args[0])

	// check the input folder
	inputFolder, err1 := os.Args[1]
	if err1 != nil {
		panic("Error: input folder not found")
	}

	// read inputvars.csv
	inputFile, err2 := os.Args[2]
	if err2 != nil {
		panic("Error: inputvars.csv not found")
	}

	// check the output folder
	output, err3 := os.Args[3]
	if err3 != nil {
		panic("Error: output folder not found")
	}

	if len(args) >= 4 {
		datadir = inputFolder + '/'
		fileans = datadir + inputFile
		outdir = datadir + output + str(foldertime) + '/'
	} else {
		print("User must specify data directory, input file name, and output file directory (e.g., at command line type main.exe ../data/ inputvariables.csv exampleout).")
		sys.exit()
	}

	// read in the input file
	inputvars = Loadfile(fileans, true)
	if len(inputvars) != 17 {
		panic("Error: inputvars.csv's column number is not correct")
	}
	population, landscape, model, mcRun, looptime, outputYear, cdmat := ReadInputParameters(inputvars[0])

	// Begin Monte-Carlo Looping
	MonteCarloLooping(mcRun, looptime, outputYear, cdmat, population, landscape, model)

}

// Monte-Carlo Looping, can use go routine
func MonteCarloLooping(mcRun int, looptime int, outputYear int, cdmat [][]float64, population Population, landscape Landscape, model Model) [][]Population {
	numProcessors := runtime.NumCPU()

	generation := make([]Population, looptime+1)
	generation[0] = population
	monteResult := make([][]Population, 0)

	output := make(chan []Population, numProcessors)

	for i := 0; i < numProcessors; i++ {

		// begain generation looping
		GenerationLooping(looptime, cdmat, landscape, model, generation, output)

	}

	// get the output
	for i := 0; i < numProcessors; i++ {
		MonteResult[i] = append(MonteResult[i], <-output...)
	}

	return MonteResult

}

// every generation looping, in single processor
func GenerationLooping(looptime int, cdmat [][]float64, landscape Landscape, model Model, generation []population, output chan []generation) {
	for i := 1; i <= looptime; i++ {
		// update the generation
		generation[i] = UpdateGeneration(generation[i-1], landscape, model, cdmat)

	}
	output <- 1

}

// update the generation
func UpdateGeneration(currentPopulation Population, landscape Landscape, model Model, cdmat [][]float64) newPopulation {
	newPopulation := CopyPop(currentPopulation)
	// update the population
	// find the mating pairs for this generation and the total number of new born individuals in this generation
	matingPair, numNewBorn := DoMate(newPopulation)
	newBornIndividuals := newPopulation.DoOffspring(matingPair)

}

// copy the population
func CopyPop(currentPopulation Population) newPopulation {
	var newPopulation Population
	newPopulation.individuals = make([]Individual, len(currentPopulation.individuals))
	for i := range newPopulation.individuals {
		newPopulation.individuals[i] = CopyInd(currentPopulation.individuals[i])
	}
	newPopulation.matureAge = currentPopulation.matureAge
	newPopulation.deathRate = currentPopulation.deathRate
	newPopulation.mateFreq = currentPopulation.mateFreq
	newPopulation.fecundity = currentPopulation.fecundity
	newPopulation.femaleRate = currentPopulation.femaleRate
	newPopulation.fitness = make([]float64, len(currentPopulation.fitness))
	for i := range newPopulation.fitness {
		newPopulation.fitness[i] = currentPopulation.fitness[i]
	}

	return newPopulation
}

// copy the individual
func CopyInd(currentIndividual Individual) newIndividual {
	var newIndividual Individual
	newIndividual.id = currentIndividual.id
	newIndividual.position.x = currentIndividual.position.x
	newIndividual.position.y = currentIndividual.position.y

	newIndividual.sex = currentIndividual.sex
	newIndividual.age = currentIndividual.age
	newIndividual.genetics = currentIndividual.genetics
	newIndividual.grid = currentIndividual.grid

	return newIndividual
}
