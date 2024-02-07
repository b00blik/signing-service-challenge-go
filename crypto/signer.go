package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/model"
)

type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type RSASigner struct {
	Device *model.Device
}

type ECCSigner struct {
	Device *model.Device
}

type KeyPair struct {
	Public  any
	Private any
}

func (r *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	rsaMarshaler := NewRSAMarshaler()
	keyPair, err := rsaMarshaler.Unmarshal(r.Device.PrivateKey)
	if err != nil {
		return nil, err
	}
	signature, err := rsa.SignPSS(
		rand.Reader,
		keyPair.Private,
		crypto.SHA256,
		dataToBeSigned,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func (e *ECCSigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	eccMarshaler := NewECCMarshaler()
	keyPair, err := eccMarshaler.Decode(e.Device.PrivateKey)
	if err != nil {
		return nil, err
	}
	signature, err := ecdsa.SignASN1(
		rand.Reader,
		keyPair.Private,
		dataToBeSigned,
	)
	if err != nil {
		return nil, err
	}
	return signature, nil
}
