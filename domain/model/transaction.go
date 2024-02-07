package model

type Transaction struct {
	ID               string
	DeviceID         string
	SignatureCounter int
	Data             []byte
	LastSignature    []byte
}
