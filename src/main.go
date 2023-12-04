package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
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

	// begin Monte-Carlo Looping
	MonteCarloLoopingMulti(mcRun, looptime, cdmat, population, landscape, model)

	// write the output file

	// darw the output figure

}

// Monte-Carlo Looping, run parallel
func MonteCarloLoopingMulti(mcRun int, looptime int, cdmat [][]float64, population Population, landscape Landscape, model Model) [][]Population {
	// get the number of processors
	numProcessors := runtime.NumCPU()

	// if the number of processors is larger than the number of Monte-Carlo run, then we only use the number of Monte-Carlo run processors
	if mcRun < numProcessors {
		numProcessors = mcRun
	}
	n := mcRun / numProcessors
	monteResult := make([][]Generation, 0)

	output := make(chan []Generation, numProcessors)

	for i := 0; i < numProcessors; i++ {
		if i < numProcessors-1 {
			n += mcRun % numProcessors
		}

		// begain generation looping
		go MonteCarloLoopingSingle(n, looptime, cdmat, landscape, model, generation, output)

	}

	// get the output
	for i := 0; i < numProcessors; i++ {
		MonteResult[i] = append(MonteResult[i], <-output...)
	}

	return MonteResult

}

// Monte-Carlo Looping, run single processor
func MonteCarloLoopingSingle(n int, looptime int, cdmat [][]float64, landscape Landscape, model Model, output chan []Generation) {

	generations := make([]Generation, n)
	// begin generation looping
	for i := 0; i < n; i++ {

		generations[i] = GenerationLooping(looptime, cdmat, landscape, model, population)

	}

	output <- generations

}

// every generation looping
func GenerationLooping(looptime int, cdmat [][]float64, landscape Landscape, model Model, population Population) Generation {
	// initialize the generation, it is the timepoints slice of population
	generation := make(Generation, looptime+1)
	generation[0] = population

	for i := 1; i <= looptime; i++ {
		// update the generation
		generation[i] = UpdateGeneration(generation[i-1], landscape, model, cdmat)
	}

	return generation

}

// update the generation
func UpdateGeneration(currentPopulation Population, landscape Landscape, model Model, cdmat [][]float64) Population {
	newPopulation := CopyPop(currentPopulation)
	// update the population
	// find the mating pairs for this generation and the total number of new born individuals in this generation
	matingPair, numNewBorn := DoMate(newPopulation)
	newBornIndividuals := newPopulation.DoOffspring(matingPair)
	deathCount := newPopulation.DoDispersal(landscape, newBornIndividuals, cdmat)
	newPopulation.AdultDeath(numNewBorn - deathCount)
	newPopulation.UpdateAge()

	return newPopulation

}

// copy the population
func CopyPop(currentPopulation Population) Population {
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
func CopyInd(currentIndividual Individual) Individual {
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
