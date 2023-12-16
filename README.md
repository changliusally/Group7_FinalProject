# Project: Gene Structure Simulation in Complex Landscapes

## Overview
 This project aims to simulate gene structure changes in several generations on different landscapes, since complex landscape is an important disturbing factor for stable gene structure of a large population. This project also provide a simple visualization for the distribution of the generation of different genotypes.

## Features
- **Simulation of Gene Structure Changes:** Simulates how gene structures evolve over multiple generations in different landscapes.
- **Population Statistics:** Provides statistics such as allel frequency and genotype frequency in each generation to understand genetic diversity.
- **Visualization:** Offers simple yet informative visualizations of the distribution of different genotypes across generations.
- **Dynamic Landscape Modeling:** Allows for the examination of genetic dynamics in complex landscapes.

  To include details about the functions and their roles in the simulation process in your README file, you can add a new section called "Module Description" or "Simulation Steps." Here's how you can integrate this into the Markdown template I provided earlier:


## Module Description

The simulation is carried out through a series of steps, each handled by different modules in the `src` directory. Below is a description of these modules and their respective roles in the simulation process:

- **main.go:** The main entry point of the simulation, orchestrating the overall process.
- **mating.go:** Manages the mating process within the population.
- **offspring.go:** Oversees the generation of offspring in the simulation.
- **dispersal.go:** Handles the dispersal process of individuals in the population.
- **death.go:** Handles the mortality aspect within the population during the simulation. (If you have a module like this)
- **drawing.go:** Responsible for creating visual outputs. It generates both PNG images and GIF animations to visualize the simulation results.
  
These modules collectively simulate the entire process, including mating, offspring generation, dispersal, and mortality, reflecting the complex dynamics of gene structure changes in varied landscapes.


## Data Preparation
To run the simulation, prepare the input data as follows:

1. **Multiple cdmatrix:** Data matrices representing various conditions.
2. **Initialized Individuals:** Starting configuration of individuals for the simulation.
3. **Parameters File (xls):** Contains essential parameters such as `xyfilename`, `mcruns`, `looptime`, `output_years`, `cdmat`, `mateFrequency`, `matureAge`, `fecundity`,`femaleRate`, `dispersalMethod`, `r`, `K_env`, `K_Grid`, `width`, `Fitness_AA`, `Fitness_Aa`, `Fitness_aa`,`offspringMethod`.



## Running the Simulation
Navigate to the source directory and execute the simulation using the following commands:
```
cd Group7_FinalProject/src
./src ../input inputvars.csv output
```
- `inputvars.csv` is the parameter file.
- `output` specifies the filenames for output data and figures.

## Outputs
- **Final Population Information:** Details about the population at the end of the simulation.
- **Generation-wise Gene Type Distribution:** Figures illustrating the distribution of different gene types in the landscape for each generation.
- **Dynamic Distribution (GIF):** An animated representation showing the temporal changes in gene type distribution across the landscape.


---
