package main

import (
	"Ricochet-Robot-Solver/internal/config"
	"Ricochet-Robot-Solver/internal/input"
	"Ricochet-Robot-Solver/internal/output"
	"Ricochet-Robot-Solver/internal/solver"
	"Ricochet-Robot-Solver/internal/tracker"
	"log"
	_ "net/http/pprof"
	"os"
	"runtime"
	"runtime/pprof"
)

func main() {
	// get program config
	conf, err := config.GetConfig("config")

	// profiler
	if conf.Mode == "profiler" {
		cpu, err := os.Create("profiling/cpu.prof")
		if err != nil {
			log.Fatal(err)
		}
		err = pprof.StartCPUProfile(cpu)
		if err != nil {
			log.Fatal(err)
		}
		defer pprof.StopCPUProfile()
	}

	// transform json data to board object
	board, initBoardState, initRobotOrder, robotStoppingPositions, err := input.GetData(conf.BoardDataLocation)
	if err != nil {
		log.Printf("Error loading board data:\n%v\n", err)
		return
	}

	// output extra information based on config
	if conf.Modes[conf.Mode]["output"].NodeNeighbors == true {
		err = output.Neighbors(&board)
		if err != nil {
			return
		}
	}
	if conf.Modes[conf.Mode]["output"].RobotStoppingPositions == true {
		err = output.RobotStoppingPositions(&robotStoppingPositions)
		if err != nil {
			return
		}
	}

	// solve the board (with tracking)
	path, trackingData, err := tracker.TrackSolver(solver.Solver, &board, initBoardState, &robotStoppingPositions, conf)
	if err != nil {
		log.Printf("\nError solving:\n%v\n", err)
		return
	}

	// output the path
	err = output.Path(path, trackingData, initRobotOrder)
	if err != nil {
		return
	}

	// profiler
	if conf.Mode == "profiler" {
		runtime.GC()
		mem, err := os.Create("profiling/memory.prof")
		if err != nil {
			log.Fatal(err)
		}
		defer func(mem *os.File) {
			err := mem.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(mem)
		if err := pprof.WriteHeapProfile(mem); err != nil {
			log.Fatal(err)
		}
	}
}
