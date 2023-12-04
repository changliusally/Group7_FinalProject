package main

//This is the mating function for choosing individual mate pairs. We assume that all pairs do
//sexual mating

import (
	"math/rand"
	"math"
	"fmt"
)

//highest level function: DoMate 
//It takes the population and the number of generation as input and 
//returns a slices of paired individuals couples that is within the appropriate distance 
//and is mature 
func DoMate(population Population) [][]Individual {
	//select individuals within appropriate age of a certain generation 
	var mateCandidateFemale []Individual
	var mateCandidateMale []Individual

	for _, indi := range population.individuals {
		if indi.age >= population.matureAge {
			if indi.sex == 0 {
				mateCandidateMale = append(mateCandidateMale, indi)
			} else {
				mateCandidateFemale = append(mateCandidateFemale, indi)
			}
		}
	}

	//we then randomly select individuals within the female and male candidates 
	n := Min(len(mateCandidateFemale), len(mateCandidateMale)) 
	mateFreq := population.mateFreq * float64(n)
	mateThreshold := population.mateThreshold
	matePair := RandomSelection(mateCandidateFemale, mateCandidateMale, mateFreq, mateThreshold)
	return matePair

}

//RandomSelection take a sclice of female candidates and a slice of male candidates and the mate frequency
//randomly pair them 
//mate without replacement 
func RandomSelection(mateCandidateFemale, mateCandidateMale []Individual, mateFreq, mateThreshold float64) [][]Individual {
	if len(mateCandidateFemale) == 0 || len(mateCandidateMale) == 0 {
		fmt.Println(len(mateCandidateFemale), len(mateCandidateMale))
		panic("No enough male/female individuals in this generation")
	}

	var pairedIndividual [][]Individual

	for i := 0; i < int(mateFreq); i++ {
		//select an individual in male and an individual in female 
		numFemale := len(mateCandidateFemale)
		numMale := len(mateCandidateMale)
		female := rand.Intn(numFemale)
		male := rand.Intn(numMale)

		femaleIndividual := mateCandidateFemale[female]
		maleIndividual := mateCandidateMale[male]
		if Distance(femaleIndividual.position, maleIndividual.position) <= mateThreshold {
			var newPair []Individual
			newPair = append(newPair, femaleIndividual)
			newPair = append(newPair, maleIndividual)
			pairedIndividual = append(pairedIndividual, newPair)
			mateCandidateFemale = Delete(mateCandidateFemale, femaleIndividual)
			mateCandidateMale = Delete(mateCandidateMale, maleIndividual)
		}
	}

	if len(pairedIndividual) == 0 {
		panic("No paired individuals in this generation")
	}
	return pairedIndividual
}

//Min takes two integers as input and returns the smaller integer
func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

//Distance takes two positions as input and returns the distance between these two positions
func Distance(p1, p2 OrderedPair) float64 {
	dX := p1.x - p2.x 
	dY := p1.y - p2.y 
	dist := math.Sqrt(dX*dX + dY*dY)
	return dist 
}

//Delete takes an individualSlice and an individual variable as input and returns 
//a slice of individuals with the individual variable deleted 
func Delete(indSclice []Individual, a Individual) []Individual {
	//check whether this individual is in the slice 
	b := false 
	var deleteIndex int
	for i := range indSclice {
		if IsEqual(indSclice[i], a) {
			b = true
			deleteIndex = i
		}
	}
	if b == false {
		panic("This individual is not in the individual slice")
	}

	indSclice = append(indSclice[:deleteIndex], indSclice[deleteIndex+1:]...)

	return indSclice
}

//IsEqual takes two Individual variables as input and returns a boolean indecating whether two individual
//variables are equal 
func IsEqual(a, b Individual) bool {
	bo := false 
	if (a.position == b.position) && (a.sex == b.sex) && (a.age == b.age) && (a.genetics == b.genetics) {
		bo = true 
	}
	return bo
}