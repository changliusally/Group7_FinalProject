package main 

//This is the part for all other functions we need

import (
	"math/rand"
	//"time"
	"fmt"
)

//AdultDeath takes the population and totalDeathCount as input and returns the population
//with certain inividuals deleted. The selection of death individuals is based on the
//population's death rate, female rate, and genetics ratio
func (population *Population) AdultDeath(totalDeathCount int) {
	//check if the totalDeathCount is negative 
	if totalDeathCount <=0 {
		fmt.Println(totalDeathCount)
	} else {
		//check how many individuals are in the current population 
		count := 0
		individualToDie := make([]Individual, 0) //all death candidates 
		for i:= range population.individuals {
			if population.individuals[i].age > 0 {
				individualToDie = append(individualToDie, population.individuals[i])
				count++
			}
		}

		/*
		femaleRate := population.femaleRate 
		//death of female and male should meet population.femaleRate
		femaleDeathCount := int(float64(totalDeathCount) * population.femaleRate)
		// Calculate the number of males to die (total deaths minus female deaths).
		maleDeathCount := totalDeathCount - femaleDeathCount

		//check how many female and male individuals 
		femaleToDie := make([]Individual,0)
		maleToDie := make([]Individual,0)

		for i := range individualToDie {
			if individualToDie[i].sex == 0 { //male 
				maleToDie = append(maleToDie, individualToDie[i])
			} else {
				femaleToDie = append(femaleToDie, individualToDie[i])
			}
		}
		femaleReal := len(femaleToDie)
		maleReal := len(maleToDie)
		*/
		//newIndividual := CopyIndividual(population.individuals)
		if count < totalDeathCount {
			totalDeathCount = count 
		}
		for i := 0; i < totalDeathCount; i++ {
			n := len(individualToDie)
			randVal := rand.Intn(n)
			indi := individualToDie[randVal]
			population.individuals = Delete(population.individuals, indi)
			individualToDie = Delete(individualToDie, indi)
		}
		//population.individuals = newIndividual
	}
	
}

func CopyIndividual(individuals []Individual) []Individual {
	newIndividual := make([]Individual, 0)
	for _, indi := range individuals {
		newIndividual = append(newIndividual, indi)
	}
	return newIndividual
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
	for i := range pop.individuals {
		pop.individuals[i].age += 1
	}

}