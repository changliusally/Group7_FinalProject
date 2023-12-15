package main 

//This is the part for doing adult death after we do offspring disperse. 

import (
	"math/rand"
	//"time"
	"fmt"
)

//AdultDeath takes the population and totalDeathCount as input and returns the population
//with certain inividuals deleted. The selection of death individuals is based on the
//population's death rate, female rate, and genetics ratio
func (population *Population) AdultDeath(totalDeathCount int) {//, pAA, pAa, paa float64) {
	//check if the totalDeathCount is negative 
	if totalDeathCount <=0 {
		fmt.Println(totalDeathCount)
	} else {
		//check how many individuals are in the current population 
		count := 0
		individualToDie := make([]Individual, 0) //all death candidates 
		AA := make([]Individual,0)
		Aa := make([]Individual,0)
		aa := make([]Individual,0)

		// take all aldults as death candidate, divide them into different gene groups
		for i:= range population.individuals {
			if population.individuals[i].age > 0 {
				individualToDie = append(individualToDie, population.individuals[i])
				count++
				if population.individuals[i].genetics == 0 {
					aa = append(aa, population.individuals[i])
				} else if population.individuals[i].genetics == 1 {
					Aa = append(Aa, population.individuals[i])
				} else {
					AA = append(AA, population.individuals[i])
				}
			}
		}

		if count < totalDeathCount {
			totalDeathCount = count 
		}
		// paa := 0.25
		// pAa := 0.375
		// pAA := 0.375

		paa := 1.0 - population.fitness[0]
  		pAa := 1.0 - population.fitness[1]
  		pAA := 1.0 - population.fitness[2]

  		sum := paa + pAa + pAA

  		paa = paa / sum 
  		pAa = pAa / sum
  		pAA = pAA / sum

		lenaa := int(float64(totalDeathCount) * paa)
		lenAa := int(float64(totalDeathCount) * pAa)
		lenAA := int(float64(totalDeathCount) * pAA)

		if lenaa > len(aa) {
			lenaa = len(aa)
		}
		if lenAa > len(Aa) {
			lenAa = len(Aa)
		}
		if lenAA > len(AA) {
			lenAA = len(AA)
		}
		total := totalDeathCount - lenaa - lenAa - lenAA 

		for i := 0; i < lenaa; i++ {
			randVal := rand.Intn(len(aa))
			population.individuals = Delete(population.individuals, aa[randVal])
			individualToDie = Delete(individualToDie, aa[randVal])
			aa = Delete(aa, aa[randVal])
		}
		for i := 0; i < lenAa; i++ {
			randVal := rand.Intn(len(Aa))
			population.individuals = Delete(population.individuals, Aa[randVal])
			individualToDie = Delete(individualToDie, Aa[randVal])
			Aa = Delete(Aa, Aa[randVal])
		}
		for i := 0; i < lenAA; i++ {
			randVal := rand.Intn(len(AA))
			population.individuals = Delete(population.individuals, AA[randVal])
			individualToDie = Delete(individualToDie, AA[randVal])
			AA = Delete(AA, AA[randVal])
		}

		if total > 0 {
			for i := 0; i <= total; i++ {
				rand := rand.Intn(len(individualToDie))
				population.individuals = Delete(population.individuals, individualToDie[rand])
				individualToDie = Delete(individualToDie, individualToDie[rand])
			}
		}
		
	}
	
}

//CopyIndividual takes a slice of individuals as input and returns a
//slice of individual that is exactly same as the input 
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