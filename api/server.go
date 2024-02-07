package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/repo"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/service"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/storage"
	"github.com/gorilla/mux"
)

type Response struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type Server struct {
	httpPort string
}

func NewServer(httpPort string) *Server {
	return &Server{
		httpPort: httpPort,
	}
}

func (s *Server) Run() error {
	mux := mux.NewRouter()
	log := log.New(os.Stdout, "[SIGNING SERVICE] ", log.LstdFlags)

	db := storage.NewStorage()
	repo := repo.NewRepository(db)
	deviceService := service.NewDeviceService(log, repo)
	transactionService := service.NewTransactionService(log, repo)

	mux.Handle("/api/v0/device", http.HandlerFunc(s.callCreateSignatureDevice(deviceService))).Methods(http.MethodPost)
	mux.Handle("/api/v0/device", http.HandlerFunc(s.callListSignatureDevices(deviceService))).Methods(http.MethodGet)
	mux.Handle("/api/v0/device/{id}", http.HandlerFunc(s.callGetSignatureDevices(deviceService))).Methods(http.MethodGet)
	mux.Handle("/api/v0/transaction", http.HandlerFunc(s.callSignTransaction(transactionService))).Methods(http.MethodPost)
	mux.Handle("/api/v0/transaction", http.HandlerFunc(s.callListTransactions(transactionService))).Methods(http.MethodGet)
	mux.Handle("/api/v0/transaction/{id}", http.HandlerFunc(s.callGetTransaction(transactionService))).Methods(http.MethodGet)
	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health)).Methods(http.MethodGet)

	return http.ListenAndServe(s.httpPort, mux)
}

func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

func WriteErrorResponse(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: err.Error(),
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(errorResponse)
}

func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
