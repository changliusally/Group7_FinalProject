package main

import (
	"math"
	"math/rand"
	"time"
)

// DoOffspring generate offspring for pairs [][]individual.
// For each pair[i], generate a number of []Individual as offspring,
// where the number of offspring for each pair of parents follows a Poisson distribution with lambda as population.fecundity.
// The position of each Individual is the same as the mother's position (pair[i][0].position).
// The genetic genotype has types A and a, from which AA, Aa, or aa are obtained based on the parents. Set age to 0.
// Combine all Individuals into a single slice []Individual. Randomly assign their genders, achieving the female proportion as population.femaleRate.
// Finally, add all Individuals to Population.
// DoOffspring generate offspring for pairs [][]individual.
// For each pair[i], generate a number of []Individual as offspring,
// where the number of offspring for each pair of parents follows a Poisson distribution with lambda as population.fecundity.
// The position of each Individual is the same as the mother's position (pair[i][0].position).
// The genetic genotype has types A and a, from which AA, Aa, or aa are obtained based on the parents. Set age to 0.
// Combine all Individuals into a single slice []Individual. Randomly assign their genders, achieving the female proportion as population.femaleRate.
// Finally, add all Individuals to Population.
func (population *Population) DoOffspring(pairs [][]Individual, landscape Landscape, method string) []Individual {
	var offspring []Individual
	if method == "constant" {
		offspring = population.DoOffspringConstant(pairs, landscape)
	}
	if method == "poisson" {
		offspring = population.DoOffspringPoisson(pairs, landscape)
	}
	if method == "noraml" {
		offspring = population.DoOffspringNormal(pairs, landscape)
	}
	if method == "random" {
		offspring = population.DoOffspringRandom(pairs, landscape)
	}

	return offspring
}

func (population *Population) DoOffspringConstant(pairs [][]Individual, landscape Landscape) []Individual {
	var offspring []Individual

	//update the offsprings of each mating pair
	for _, pair := range pairs {
		// Determine the number of offspring
		numOffspring := population.fecundity
		//create a slice of Individuals representing the offsprings of this mating couple
		pairoffspring := make([]Individual, numOffspring)

		//update each offspring
		for i := 0; i < numOffspring; i++ {
			// Create new offspring
			pairoffspring[i].age = 0

			// Inherits mother's position
			pairoffspring[i].position.x, pairoffspring[i].position.y = RandomGridxy(pair[0].gridIn, landscape)
			//pairoffspring[i].position.y = pair[0].position.y +

			pairoffspring[i].gridIn = pair[0].gridIn

			// decide genetic genotype based on the parents.
			pairoffspring[i].genetics = generateGenetics(pair[0].genetics, pair[1].genetics)

			// Randomly assign sex based on female rate
			if rand.Float64() < population.femaleRate {
				pairoffspring[i].sex = 1 // Female
			} else {
				pairoffspring[i].sex = 0 // Male
			}

			offspring = append(offspring, pairoffspring[i])
			//population.individuals = append(population.individuals, pairoffspring[i])
		}
	}

	return offspring
}

// function poisson takes a lambda parameter and return a integer representing the number
// of offsprings of a mating couple is the population's breeding is poisson distribution
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

// generateGenetics takes two integers as parents' genetics and return an integer
// representing the offspring's genetics
func generateGenetics(father, mother int) int {
	// Simple genetic model: AA (2), Aa (1), aa (0)

	//generate a random number
	randVal := rand.Float64()

	// Both parents are AA
	if father == 2 && mother == 2 {
		return 2 // Offspring is AA
	}

	// One parent is AA and the other is Aa
	if (father == 2 && mother == 1) || (father == 1 && mother == 2) {
		if randVal < 0.5 {
			return 2 // AA
		} else {
			return 1 // Aa
		}
	}

	// Both parents are Aa
	if father == 1 && mother == 1 {
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

	//one parent is AA and the other one is aa
	if (father == 2 && mother == 0) || (father == 0 && mother == 2) {
		return 1
	}

	// One parent is aa and the other is Aa
	if (father == 0 && mother == 1) || (father == 1 && mother == 0) {
		if randVal < 0.5 {
			return 0 // aa
		} else {
			return 1 // Aa
		}
	}
	return 0 // Default case, can be adjusted as needed
}

func (population *Population) DoOffspringPoisson(pairs [][]Individual, landscape Landscape) []Individual {
	rand.Seed(time.Now().UnixNano())
	var offspring []Individual

	for _, pair := range pairs {
		// Determine the number of offspring using Poisson distribution
		numOffspring := Poisson(float64(population.fecundity))
		pairoffspring := make([]Individual, numOffspring)

		for i := 0; i < numOffspring; i++ {

			pairoffspring[i].age = 0

			// Inherits mother's position
			pairoffspring[i].position.x, pairoffspring[i].position.y = RandomGridxy(pair[0].gridIn, landscape)

			// Inherits mother's grid
			pairoffspring[i].gridIn = pair[0].gridIn

			// decide genetic genotype based on the parents.
			pairoffspring[i].genetics = generateGenetics(pair[0].genetics, pair[1].genetics)

			// Randomly assign sex based on female rate
			if rand.Float64() < population.femaleRate {
				pairoffspring[i].sex = 1 // Female
			} else {
				pairoffspring[i].sex = 0 // Male
			}

			offspring = append(offspring, pairoffspring[i])
		}
	}

	// Add the offspring to the population
	//population.individuals = append(population.individuals, offspring...)

	return offspring
}

func Poisson(lambda float64) int {
	L := math.Exp(-lambda)
	k := 0
	p := 1.0

	for p > L {
		k++
		p *= rand.Float64()
	}

	return k - 1
}

func (population *Population) DoOffspringNormal(pairs [][]Individual, landscape Landscape) []Individual {
	rand.Seed(time.Now().UnixNano())
	var offspring []Individual

	for _, pair := range pairs {
		// Determine the number of offspring using Normal distribution
		numOffspring := Normal(float64(population.fecundity), 1)
		if numOffspring < 0 {
			numOffspring = 0 // Ensure non-negative
		}
		pairoffspring := make([]Individual, numOffspring)

		for i := 0; i < numOffspring; i++ {
			pairoffspring[i].age = 0

			// Inherits mother's position
			pairoffspring[i].position.x, pairoffspring[i].position.y = RandomGridxy(pair[0].gridIn, landscape)

			// Inherits mother's grid
			pairoffspring[i].gridIn = pair[0].gridIn

			// decide genetic genotype based on the parents.
			pairoffspring[i].genetics = generateGenetics(pair[0].genetics, pair[1].genetics)

			// Randomly assign sex based on female rate
			if rand.Float64() < population.femaleRate {
				pairoffspring[i].sex = 1 // Female
			} else {
				pairoffspring[i].sex = 0 // Male
			}

			offspring = append(offspring, pairoffspring[i])
		}

	}

	// Add the offspring to the population
	// population.individuals = append(population.individuals, offspring...)

	return offspring
}

func Normal(mean, stddev float64) int {
	return int(rand.NormFloat64()*stddev + mean)
}

func (population *Population) DoOffspringRandom(pairs [][]Individual, landscape Landscape) []Individual {
	rand.Seed(time.Now().UnixNano())
	var offspring []Individual

	for _, pair := range pairs {
		// Determine the number of offspring using a random distribution from 0 to fecundity
		numOffspring := RandomFecundity(float64(population.fecundity))
		pairoffspring := make([]Individual, numOffspring)

		for i := 0; i < numOffspring; i++ {
			pairoffspring[i].age = 0

			// Inherits mother's position
			pairoffspring[i].position.x, pairoffspring[i].position.y = RandomGridxy(pair[0].gridIn, landscape)

			// Inherits mother's grid
			pairoffspring[i].gridIn = pair[0].gridIn

			// decide genetic genotype based on the parents.
			pairoffspring[i].genetics = generateGenetics(pair[0].genetics, pair[1].genetics)

			// Randomly assign sex based on female rate
			if rand.Float64() < population.femaleRate {
				pairoffspring[i].sex = 1 // Female
			} else {
				pairoffspring[i].sex = 0 // Male
			}

			offspring = append(offspring, pairoffspring[i])
		}

	}

	// Add the offspring to the population
	// population.individuals = append(population.individuals, offspring...)

	return offspring
}

func RandomFecundity(fecundity float64) int {
	return int(rand.Float64() * fecundity)
}
