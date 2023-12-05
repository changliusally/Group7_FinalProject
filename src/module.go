package main 

//This is the part for all other functions we need

import (
	"math/rand"
	"time"
)

//AdultDeath takes the population and totalDeathCount as input and returns the population
//with certain inividuals deleted. The selection of death individuals is based on the
//population's death rate, female rate, and genetics ratio
func (population *Population) AdultDeath(totalDeathCount int) {
	rand.Seed(time.Now().UnixNano())

	// Calculate the total number of deaths based on the population size and death rate.
	//number:=len(population.individuals)
	//totalDeathCount := int(float64(number) * population.deathRate)

	//death of female and male should meet population.femaleRate
	femaleDeathCount := int(float64(totalDeathCount) * population.femaleRate)

	// Calculate the number of males to die (total deaths minus female deaths).
	maleDeathCount := totalDeathCount - femaleDeathCount

	//Create a map to track the number of deaths for each genetic type (0, 1, 2) based on a 1:2:1 ratio.
	geneticsCounts := map[int]int{
		0: totalDeathCount / 4,     // Genotype aa
		1: totalDeathCount / 2,     // Genotype Aa
		2: totalDeathCount / 4,     // Genotype AA
	}

	// choose by random which individual to die and record the individual
	// Slice to store the IDs of individuals who will die.
	var deadIndividuals []int

	// Continue selecting individuals for death until the required numbers are reached for each sex and genetic type.
	for maleDeathCount > 0 || femaleDeathCount > 0 {
		totalIndividual := len(population.individuals)
		index := rand.Intn(totalIndividual)
		individual := population.individuals[index]

		// Skip if the individual is already selected for death.
		//if contains(deadIndividuals, individual.id) || individual.age == 0 {
		if contains(deadIndividuals, individual.id) {
			continue
		}

		if individual.sex == 1 && femaleDeathCount > 0 && geneticsCounts[individual.genetics] > 0 {

			// Check and select female individuals for death.
			femaleDeathCount--
			geneticsCounts[individual.genetics]--
			deadIndividuals = append(deadIndividuals, individual.id)

		} else if individual.sex == 0 && maleDeathCount > 0 && geneticsCounts[individual.genetics] > 0 {

			// Check and select male individuals for death.
			maleDeathCount --
			geneticsCounts[individual.genetics] --
			deadIndividuals = append(deadIndividuals, individual.id)
		}
	}

	// Remove dead individuals from the population
	var updatedIndividuals []Individual
	for _, individual := range population.individuals {
		if !contains(deadIndividuals, individual.id) {
			updatedIndividuals = append(updatedIndividuals, individual)
		}
	}
	population.individuals = updatedIndividuals

}

//contains takes a slice and an element as input as returns a boolean 
//to checks if a slice contains a specific element
func contains(slice []int, element int) bool {
	for _, ele := range slice {
		if ele == element {
			return true
		}
	}
	return false
}


//UpdateAge takes a population as input and update all individuals' age of this population 
func (pop *Population)UpdateAge() {
	for _, inid := range pop.individuals {
		inid.age = inid.age + 1
	}

}