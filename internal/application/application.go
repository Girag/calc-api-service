package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Girag/calc-api-service/pkg/calculation"
)

type Config struct {
	Port string
	Path string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Port = os.Getenv("CALC_API_PORT")
	config.Path = os.Getenv("CALC_API_PATH")
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.Path == "" {
		config.Path = "/api/v1/calculate"
	}
	return config
}

type Application struct {
	config *Config
}

func NewApp() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

type Request struct {
	Expression string `json:"expression"`
}

type SuccessResponse struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rr, r)

		log.Printf(
			"[%s] %s %s %d %s",
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			rr.statusCode,
			time.Since(start),
		)
	})
}

func HandleError(w http.ResponseWriter, err error) {
	makeErrorResponse := func(errMessage string, statusCode int) {
		response := &ErrorResponse{Error: errMessage}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		data, _ := json.Marshal(response)
		w.Write(data)
	}

	switch {
	case errors.Is(err, ErrMethodNotAllowed):
		makeErrorResponse(err.Error(), http.StatusMethodNotAllowed)
	case errors.Is(err, ErrNotFound):
		makeErrorResponse(err.Error(), http.StatusNotFound)
	case errors.Is(err, ErrBadRequest):
		makeErrorResponse(err.Error(), http.StatusBadRequest)
	case errors.Is(err, calculation.ErrInvalidCharInExpression),
		errors.Is(err, calculation.ErrInvalidExpression),
		errors.Is(err, calculation.ErrDivisionByZero),
		errors.Is(err, calculation.ErrOpeningParenthesisMissing),
		errors.Is(err, calculation.ErrClosingParenthesisMissing):
		makeErrorResponse(err.Error(), http.StatusUnprocessableEntity)
	case errors.Is(err, ErrInternalServerError):
		makeErrorResponse(err.Error(), http.StatusInternalServerError)
	default:
		makeErrorResponse(ErrUnknown.Error(), http.StatusInternalServerError)
	}
}

func (a *Application) CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != a.config.Path {
		fmt.Println("Path not found")
		HandleError(w, ErrNotFound)
		return
	}

	if r.Method != http.MethodPost {
		HandleError(w, ErrMethodNotAllowed)
		return
	}

	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil || request == nil || request.Expression == "" {
		HandleError(w, ErrBadRequest)
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		HandleError(w, err)
	} else {
		response := &SuccessResponse{Result: result}
		w.Header().Set("Content-Type", "application/json")
		data, err := json.Marshal(response)
		if err != nil {
			HandleError(w, err)
			return
		}
		w.Write(data)
	}
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.config.Path != r.URL.Path {
		HandleError(w, ErrNotFound)
		return
	}

	a.CalcHandler(w, r)
}

func (a *Application) RunServer() error {
	mux := http.NewServeMux()
	mux.HandleFunc(a.config.Path, a.CalcHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleError(w, ErrNotFound)
	})

	loggedMux := LoggingMiddleware(mux)

	port := fmt.Sprintf(":%s", a.config.Port)
	log.Printf("Starting server on port %s", a.config.Port)
	return http.ListenAndServe(port, loggedMux)
}
