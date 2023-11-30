package main

import (
	"math"
	"math/rand"
	"time"
)

//DoOffspring generate offspring for pairs [][]individual.
//For each pair[i], generate a number of []Individual as offspring,
// where the number of offspring for each pair of parents follows a Poisson distribution with lambda as population.fecundity.
//The position of each Individual is the same as the mother's position (pair[i][0].position).
//The genetic genotype has types A and a, from which AA, Aa, or aa are obtained based on the parents. Set age to 0.
//Combine all Individuals into a single slice []Individual. Randomly assign their genders, achieving the female proportion as population.femaleRate.
//Finally, add all Individuals to Population.

func DoOffspring(pairs [][]Individual, population Population, generation int) []Individual {
	rand.Seed(time.Now().UnixNano())
	var offspring []Individual

	for _, pair := range pairs {
		// Determine the number of offspring
		numOffspring := population.fecundity
		var pairoffspring []Individual

		for i := 0; i < numOffspring; i++ {
			// Create new offspring
			pairoffspring[i].age = 0

			// Inherits mother's position
			pairoffspring[i].position.x = pair[0].position.x
			pairoffspring[i].position.y = pair[0].position.y

			// decide genetic genotype based on the parents.
			pairoffspring[i].genetics = generateGenetics(pair[0].genetics, pair[1].genetics)

			// make the generation ++
			pairoffspring[i].gen = generation + 1

			// Randomly assign sex based on female rate
			if rand.Float64() < population.femaleRate {
				pairoffspring[i].sex = 1 // Female
			} else {
				pairoffspring[i].sex = 0 // Male
			}

			offspring = append(offspring, pairoffspring[i])
			population.individuals = append(population.individuals, pairoffspring[i])
		}
	}

	return offspring
}

func poisson(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= rand.Float64()
	}
	return k - 1
}

func generateGenetics(father, mother int) int {
	// Simple genetic model: AA (2), Aa (1), aa (0)

	// Both parents are AA
	if father == 2 && mother == 2 {
		return 2 // Offspring is AA
	}

	// One parent is AA and the other is Aa
	if (father == 2 && mother == 1) || (father == 1 && mother == 2) {
		if rand.Float64() < 0.5 {
			return 2 // AA
		} else {
			return 1 // Aa
		}
	}

	// Both parents are Aa
	if father == 1 && mother == 1 {
		randVal := rand.Float64()
		if randVal < 0.25 {
			return 0 // aa
		} else if randVal >= 0.25 && randVal < 0.75 {
			return 1 // Aa
		} else {
			return 2 // AA
		}
	}

	// Both parents are aa
	if father == 0 && mother == 0 {
		return 0 // Offspring is aa
	}

	// One parent is aa and the other is Aa or AA
	if (father == 0 && mother == 1) || (father == 1 && mother == 0) || (father == 0 && mother == 2) || (father == 2 && mother == 0) {
		if rand.Float64() < 0.5 {
			return 0 // aa
		} else {
			if mother == 1 || father == 1 {
				return 1 // Aa
			} else {
				return 2 // AA
			}
		}
	}
	return 0 // Default case, can be adjusted as needed
}
