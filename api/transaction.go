package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/service"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/validation"
	"github.com/gorilla/mux"
)

func (s *Server) callSignTransaction(service service.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &validation.SignTransactionRequest{}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			WriteErrorResponse(w, http.StatusInternalServerError, nil)
			return
		}

		if err := json.Unmarshal(body, req); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		response, err := service.SignTransaction(req)

		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		WriteAPIResponse(w, http.StatusOK, response)
	}
}

func (s *Server) callListTransactions(service service.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deviceID := r.URL.Query().Get("device_id")

		req := &validation.ListTransactionRequest{DeviceID: deviceID}

		response, err := service.ListTransaction(req)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		WriteAPIResponse(w, http.StatusOK, response)
	}
}

func (s *Server) callGetTransaction(service service.TransactionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		req := &validation.GetTransactionRequest{ID: vars["id"]}

		response, err := service.GetTransaction(req)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		WriteAPIResponse(w, http.StatusOK, response)
	}
}
