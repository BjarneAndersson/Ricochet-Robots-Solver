package webserver

import (
	"Ricochet-Robot-Solver/internal/input"
	"Ricochet-Robot-Solver/internal/output"
	"Ricochet-Robot-Solver/internal/precomputation"
	"Ricochet-Robot-Solver/internal/solver"
	"Ricochet-Robot-Solver/internal/tracker"
	"encoding/json"
	"net/http"
)

// loadBoard handles the loading of the board data.
//
// This function resets the application state, loads the board data from a JSON input,
// and updates the application state with the new board data. It sends a success response
// if the board data is loaded successfully, otherwise it sends an error response.
//
// Returns:
//   - http.HandlerFunc: A function that handles HTTP requests and writes responses.
func (s *Server) loadBoard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		s.Logger.Info("Resetting application state")
		s.initAppState()

		s.Logger.Info("Loading board data")

		// transform json data to board object
		board, initBoardState, initRobotOrder, err := input.GetData(*(s.Config))
		if err != nil {
			s.Logger.Error(err, "Loading board data")
			http.Error(w, "Failed to load board data", http.StatusInternalServerError)
			return
		}

		// Set the board in the application state
		s.AppState.Board = board
		s.AppState.InitBoardState = initBoardState
		s.AppState.InitRobotOrder = initRobotOrder

		// Send a success response
		w.WriteHeader(http.StatusOK)
	}
}

// solve handles the solving of the board.
//
// This function solves the board using the solver and tracker, and updates the application state
// with the path and tracking data. It sends a success response if the solving is successful,
// otherwise it logs the error.
//
// Returns:
//   - http.HandlerFunc: A function that handles HTTP requests and writes responses.
func (s *Server) solve() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// PRECOMPUTE
		s.Logger.Info("Precomputing minimal move counts")

		// precompute the minimal moves that a robot has to make in order to get to the target node, if it could stop wherever it wants
		err := precomputation.PrecomputeMinimalMoveCounts(&s.AppState.Board)

		// precompute the robot stopping positions
		robotStoppingPositions, err := precomputation.PrecomputeRobotMoves(&s.AppState.Board)
		if err != nil {
			s.Logger.Error(err, "Precomputing robot moves")
			http.Error(w, "Failed to precompute robot moves", http.StatusInternalServerError)
			return
		}

		// Set the robot stopping positions in the application state
		s.AppState.RobotStoppingPositions = robotStoppingPositions

		// SOLVE
		s.Logger.Info("Solving board")

		// solve the board (with tracking)
		path, trackingData, err := tracker.TrackSolver(solver.Solver, &s.AppState.Board, s.AppState.InitBoardState, &s.AppState.RobotStoppingPositions, *(s.Config))
		if err != nil {
			s.Logger.Error(err, "Solving")
			return
		}

		s.Logger.Info("Solved board")

		// Set the path and tracking data in the application state
		s.AppState.Path = path
		s.AppState.TrackingData = trackingData

		// Send a success response
		w.WriteHeader(http.StatusOK)
	}
}

// outputConsole handles the outputting of the board state and path to the console.
//
// This function outputs extra information based on the configuration, including neighbors,
// robot stopping positions, and the path. It sends a success response if the output is successful,
// otherwise it sends an error response.
//
// Returns:
//   - http.HandlerFunc: A function that handles HTTP requests and writes responses.
func (s *Server) outputConsole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		s.Logger.Info("Outputting to console")

		// output extra information based on config
		err = output.Neighbors(&s.AppState.Board, s.Recorder)
		if err != nil {
			s.Logger.Warn(err, "Outputting neighbors")
			http.Error(w, "Failed to output neighbors", http.StatusInternalServerError)
		}

		err = output.RobotStoppingPositions(&s.AppState.RobotStoppingPositions, s.Recorder)
		if err != nil {
			s.Logger.Warn(err, "Outputting robot stopping positions")
			http.Error(w, "Failed to output robot stopping positions", http.StatusInternalServerError)
		}

		// output the path
		s.AppState.Output.Path, err = output.Path(s.AppState.Path, s.AppState.TrackingData, s.AppState.InitRobotOrder, s.Recorder)
		if err != nil {
			s.Logger.Fatal(err, "Outputting path")
			http.Error(w, "Failed to output path", http.StatusInternalServerError)
			return
		}

		s.AppState.Output.TrackingData = s.AppState.TrackingData
		s.Logger.Warn(nil, s.AppState.Output.TrackingData.Duration)

		outputJson, err := json.Marshal(s.AppState.Output)
		if err != nil {
			s.Logger.Error(err, "Marshalling output JSON")
			http.Error(w, "Failed to marshal output JSON", http.StatusInternalServerError)
			return
		}

		// Send a success response
		w.WriteHeader(http.StatusOK)
		w.Write(outputJson)
	}
}
