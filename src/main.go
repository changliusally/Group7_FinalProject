package main

import "time"

func main() {

	// set up the timer
	start_time = time.Now()
	foldertime = start_time.strftime("%Y%m%d_%H%M%S")

	// read agv from command line
	// the format: main.exe C:/.../inputfolder inputvars.csv output_folder

	// check the input folder
	inputFolder, err1 := os.args[1]
	if err1 != nil {
		panic("Error: input folder not found")
	}

	// read inputvars.csv
	inputFile, err2 := os.args[2]
	if err2 != nil {
		panic("Error: inputvars.csv not found")
	}

	// check the output folder
	output, err3 := os.args[3]
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
	population, landscape, model, mcRun, looptime, outputYear, cdmatData := ReadInputParameters(inputvars[0])

}
