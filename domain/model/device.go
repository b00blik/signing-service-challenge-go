package model

type Device struct {
	ID               string
	Label            string
	Algorithm        string
	PrivateKey       []byte
	SignatureCounter int
}
