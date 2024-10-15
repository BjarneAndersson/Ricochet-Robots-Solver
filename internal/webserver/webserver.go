package webserver

import (
	"Ricochet-Robot-Solver/internal/application"
	"Ricochet-Robot-Solver/internal/config"
	"html/template"
	"net/http"
	_ "os"

	// External packages
	"github.com/gorilla/mux"

	// Internal packages
	logging "Ricochet-Robot-Solver/internal/log"
)

// Server represents the webserver.
//
// Fields:
//   - Router: A pointer to a mux.Router instance for handling routes.
//   - URL: The URL where the server will be hosted.
//   - Port: The port on which the server will listen.
//   - Logger: A pointer to a logging.Logger instance for logging purposes.
type Server struct {
	*mux.Router

	URL  string
	Port string

	AppState application.AppState

	Config *config.Config

	Logger   *logging.Logger
	Recorder *logging.SolveRecorder
}

// NewServer creates a new webserver instance.
// It initializes the router, URL, port, and logger, sets up the routes, and returns the server instance.
//
// Parameters:
//   - url: The URL where the server will be hosted.
//   - port: The port on which the server will listen.
//   - logger: A pointer to a logging.Logger instance for logging purposes.
//
// Returns:
//   - *Server: A pointer to the newly created Server instance.
func NewServer(url string, port string, logger *logging.Logger, solveRecorder *logging.SolveRecorder, cfg *config.Config) *Server {
	// create the webserver
	s := &Server{
		Router:   mux.NewRouter(),
		URL:      url,
		Port:     port,
		Config:   cfg,
		Logger:   logger,
		Recorder: solveRecorder,
	}

	s.initAppState()

	s.Host(url)
	s.setupRoutes()

	return s
}

func (s *Server) initAppState() {
	s.AppState = application.AppState{}
}

// setupRoutes sets up the routes for the web server.
// It defines the routes for the main page, APIs, and static assets.
func (s *Server) setupRoutes() {
	// pages
	s.HandleFunc("/", s.renderBoard()).Methods("GET")

	s.Logger.Info("Root route set up")

	// apis
	s.HandleFunc("/api/change-quarter-boards", s.changeAllQuarterBoards()).Methods("POST")
	s.HandleFunc("/api/solve", s.solve()).Methods("POST")
	s.HandleFunc("/api/output/console", s.outputConsole()).Methods("POST")

	// static assets
	s.PathPrefix("/").Handler(http.StripPrefix("/assets", http.FileServer(http.Dir("./website/static"))))
}

// root is the root route handler for the web server.
// It handles HTTP GET requests to the root URL ("/").
//
// Returns:
//   - http.HandlerFunc: A function that handles HTTP requests and writes responses.
func (s *Server) root() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// logging
		var err error

		s.Logger.Info("Root route accessed")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// load the board
		err = s.loadBoard()
		if err != nil {
			s.Logger.Error(err, "Loading board data")
			http.Error(w, "Failed to load board data", http.StatusInternalServerError)
			return
		}

		// parse the template
		html, err := template.ParseFiles("./website/templates/index.html")
		if err != nil {
			s.Logger.Error(err, "Error parsing template")
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}

		// render the template
		err = html.Execute(w, nil)
		if err != nil {
			s.Logger.Error(err, "Error formatting template")
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
