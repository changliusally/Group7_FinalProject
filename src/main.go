package main

import (
	"fmt"
	"gifhelper"
	"os"
	"runtime"
	"strconv"
	"time"
)

func main() {

	// set up the timer
	start_time := time.Now()
	foldertime := start_time.Format("20060102150405")

	// read agv from command line
	// the format: main.exe C:/.../inputfolder inputvars.csv output_folder

	fmt.Println(os.Args[0])

	// check the input folder
	inputFolder := os.Args[1]

	// read inputvars.csv
	inputFile := os.Args[2]

	// check the output folder
	output := os.Args[3]

	var fileans string
	var outdir string
	var datadir string
	if len(os.Args) >= 4 {
		datadir = inputFolder + string('/')
		fileans = datadir + inputFile
		outdir = "output/" + output + foldertime + string('/')
	} else {
		panic("User must specify data directory, input file name, and output file directory (e.g., at command line type main.exe ../data/ inputvariables.csv exampleout).")
	}

	// read in the input file
	inputvars := Loadfile(fileans, true)
	if len(inputvars[0]) != 18 {
		panic("Error: inputvars.csv's column number is not correct")
	}
	population, landscape, model, mcRun, looptime, outputYear, cdmat := ReadInputParameters(inputvars[0], datadir)

	fmt.Println("Input file is read")

	fmt.Println("Now, begin Monte-Carlo Looping")
	// begin Monte-Carlo Looping
	generations := MonteCarloLoopingMulti(mcRun, looptime, cdmat, population, landscape, model)

	fmt.Println("Monte-Carlo Looping is complete")

	// write the output file
	WriteOutput(generations, outputYear, outdir, landscape)
	// write the summary file
	WriteSummary(generations, outdir)

	fmt.Println("Output file is written")

	// darw the output figure
	for i := 0; i < mcRun; i++ {
		images := AnimateSystem(generations[i].population, landscape, 1) //animate the timepoints

		fmt.Println("images drawn!")

		fmt.Println("generate GIF")

		outputFile := "PopulationSimulation_cdmatrix16_" + strconv.Itoa(i) //output file name

		gifhelper.ImagesToGIF(images, outdir+outputFile) //draw the image and store in output folder
	}
	fmt.Println("Simulation complete!")
}

// Monte-Carlo Looping, run parallel
// input: the int number of Monte-Carlo run, the int number of generation looping, the cd matrix, the population, the landscape, the model struct
// output: a slice of generation
func MonteCarloLoopingMulti(mcRun int, looptime int, cdmat [][]float64, population Population, landscape Landscape, model Model) []Generation {
	// get the number of processors
	numProcessors := runtime.NumCPU()

	// if the number of processors is larger than the number of Monte-Carlo run, then we only use the number of Monte-Carlo run processors
	if mcRun < numProcessors {
		numProcessors = mcRun
	}
	n := mcRun / numProcessors
	MonteResult := make([]Generation, 0)

	// create a channel to store the output
	output := make(chan []Generation, numProcessors)

	for i := 0; i < numProcessors; i++ {
		if i < numProcessors-1 {
			n += mcRun % numProcessors
		}

		// begain generation looping
		go MonteCarloLoopingSingle(n, looptime, cdmat, landscape, model, population, output)

	}

	// get the output
	for i := 0; i < numProcessors; i++ {
		MonteResult = append(MonteResult, <-output...)
	}

	return MonteResult

}

// Monte-Carlo Looping, run single processor
// input: the int number of Monte-Carlo run, the int number of generation looping, the cd matrix, the population, the landscape, the model struct
// output: a slice of generation
func MonteCarloLoopingSingle(n int, looptime int, cdmat [][]float64, landscape Landscape, model Model, population Population, output chan []Generation) {

	generations := make([]Generation, n)
	// begin generation looping
	for i := 0; i < n; i++ {

		generations[i] = GenerationLooping(looptime, cdmat, landscape, model, population)

	}

	output <- generations

}

// GenerationLooping in each Monte-Carlo run
// input: the int number of generation looping, the cd matrix, the population, the landscape, the model struct
// output: a generation
func GenerationLooping(looptime int, cdmat [][]float64, landscape Landscape, model Model, population Population) Generation {
	// initialize the generation, it is the timepoints slice of population
	generation := make([]Population, looptime+1)
	generation[0] = population

	for i := 1; i <= looptime; i++ {
		// update the generation
		generation[i] = UpdateGeneration(generation[i-1], landscape, model, cdmat)
	}

	var newGeneration Generation
	newGeneration.population = generation
	return newGeneration

}

// update the generation, each generation looping do mating, offspring, dispersal, and mortality, and update the population
// input: the current population, the landscape, the model struct, the cd matrix
// output: the updated population
func UpdateGeneration(currentPopulation Population, landscape Landscape, model Model, cdmat [][]float64) Population {
	newPopulation := CopyPop(currentPopulation)
	// update the population
	// find the mating pairs for this generation and the total number of new born individuals in this generation
	matingPair, _ := DoMate(newPopulation)
	//fmt.Println("mating finish")
	newBornIndividuals := newPopulation.DoOffspring(matingPair, landscape, newPopulation.offspringMethod)
	//fmt.Println("offspring finish")
	//covert cd matrix to probability matrix
	probMatrix := CalProb(newPopulation.dispersalMethod, cdmat)
	deathCount := newPopulation.DoDispersal(landscape, newBornIndividuals, probMatrix)
	//fmt.Println("disperse finish")
	deathCountNew := len(newBornIndividuals) - deathCount
	newPopulation.AdultDeath(deathCountNew)
	//fmt.Println("mortility finish")
	newPopulation.UpdateAge()
	//newPopulation.UpdateGrid(landscape)
	//newIndividual := FindGrid(landscape, newPopulation.individuals)
	//newPopulation.individuals = newIndividual

	return newPopulation

}

// copy the population
func CopyPop(currentPopulation Population) Population {
	var newPopulation Population
	newPopulation.individuals = make([]Individual, len(currentPopulation.individuals))
	// copy the individuals
	for i := range newPopulation.individuals {
		newPopulation.individuals[i] = CopyInd(currentPopulation.individuals[i])
	}
	// copy other fields
	newPopulation.matureAge = currentPopulation.matureAge
	newPopulation.deathRate = currentPopulation.deathRate
	newPopulation.mateFreq = currentPopulation.mateFreq
	newPopulation.fecundity = currentPopulation.fecundity
	newPopulation.femaleRate = currentPopulation.femaleRate
	newPopulation.dispersalMethod = currentPopulation.dispersalMethod
	newPopulation.offspringMethod = currentPopulation.offspringMethod
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
	newIndividual.gridIn = currentIndividual.gridIn

	return newIndividual
}
