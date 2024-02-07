package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/service"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/validation"
	"github.com/gorilla/mux"
)

func (s *Server) callCreateSignatureDevice(service service.DeviceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &validation.CreateSignatureDeviceRequest{}
		body, err := io.ReadAll(r.Body)

		if err != nil || len(body) == 0 {
			WriteErrorResponse(w, http.StatusInternalServerError, nil)
			return
		}

		if err := json.Unmarshal(body, req); err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		response, err := service.CreateSignatureDevice(req)

		if err != nil {
			WriteErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		WriteAPIResponse(w, http.StatusCreated, response)
	}
}

func (s *Server) callListSignatureDevices(service service.DeviceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &validation.ListSignatureDeviceRequest{
			ID:        r.URL.Query().Get("id"),
			Label:     r.URL.Query().Get("label"),
			Algorithm: r.URL.Query().Get("algorithm"),
		}

		response, err := service.ListSignatureDevice(req)

		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		WriteAPIResponse(w, http.StatusOK, response)
	}
}

func (s *Server) callGetSignatureDevices(service service.DeviceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		req := &validation.GetSignatureDeviceRequest{ID: vars["id"]}

		response, err := service.GetSignatureDevice(req)
		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		WriteAPIResponse(w, http.StatusOK, response)
	}
}
