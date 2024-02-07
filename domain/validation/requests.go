package validation

import (
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/model"
)

type CreateSignatureDeviceRequest struct {
	ID        string `json:"id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label,omitempty"` //?
}

type ListSignatureDeviceRequest struct {
	ID        string `json:"id,omitempty"`
	Label     string `json:"label,omitempty"`
	Algorithm string `json:"algorithm,omitempty"`
}

type GetSignatureDeviceRequest struct {
	ID string
}

type ListTransactionRequest struct {
	DeviceID string `json:"device_id,omitempty"`
}

type GetTransactionRequest struct {
	ID string
}

type SignTransactionRequest struct {
	DeviceID string `json:"device_id"`
	Data     []byte `json:"data"`
}

type CreateSignatureDeviceResponse struct {
	Status string `json:"status"`
}

type ListSignatureDeviceResponse struct {
	Device []*model.Device `json:"devices"`
}

type GetSignatureDeviceResponse struct {
	Device *model.Device `json:"device"`
}

type SignTransactionResponse struct {
	Transaction string `json:"signature"`
	SignedData  string `json:"signed_data"`
}

type ListTransactionResponse struct {
	Transaction []*model.Transaction `json:"transactions"`
}

type GetTransactionResponse struct {
	Transaction *model.Transaction `json:"transaction"`
}

func (c *CreateSignatureDeviceRequest) IsValid() error {
	if c.ID == "" {
		return errors.New("\"id\" field is missing")
	}
	if c.Algorithm == "" {
		return errors.New("\"algorithm\" field is missing")
	}
	if c.Algorithm != "ECC" && c.Algorithm != "RSA" {
		return errors.New("unknown algorithm")
	}
	return nil
}

func (s *SignTransactionRequest) IsValid() error {
	if s.DeviceID == "" || s.Data == nil {
		return errors.New("\"device_id\" field is missing")
	}
	if s.Data == nil {
		return errors.New("\"data\" field is missing")
	}
	return nil
}
