package service

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/model"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/repo"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/validation"
)

type DeviceService interface {
	CreateSignatureDevice(req *validation.CreateSignatureDeviceRequest) (*validation.CreateSignatureDeviceResponse, error)
	ListSignatureDevice(req *validation.ListSignatureDeviceRequest) (*validation.ListSignatureDeviceResponse, error)
	GetSignatureDevice(req *validation.GetSignatureDeviceRequest) (*validation.GetSignatureDeviceResponse, error)
}

type deviceService struct {
	collection repo.Collection
	logger     *log.Logger
}

func NewDeviceService(logger *log.Logger, collection repo.Collection) DeviceService {
	return &deviceService{
		logger:     logger,
		collection: collection,
	}
}

func (d *deviceService) CreateSignatureDevice(req *validation.CreateSignatureDeviceRequest) (*validation.CreateSignatureDeviceResponse, error) {
	var err error
	if err := req.IsValid(); err != nil {
		return nil, err
	}

	device := &model.Device{
		ID:        req.ID,
		Label:     req.Label,
		Algorithm: req.Algorithm,
	}

	eccGenerator := &crypto.ECCGenerator{}
	eccMarshaler := crypto.NewECCMarshaler()

	rsaGenerator := &crypto.RSAGenerator{}
	rsaMarshaler := crypto.NewRSAMarshaler()

	switch req.Algorithm {
	case "ECC":
		device.PrivateKey, err = d.CreateECCKeys(*eccGenerator, eccMarshaler)
		if err != nil {
			return nil, err
		}
	case "RSA":
		device.PrivateKey, err = d.CreateRSAKeys(*rsaGenerator, rsaMarshaler)
		if err != nil {
			return nil, err
		}
	}

	_, err = d.collection.CreateSignatureDevice(device)
	if err != nil {
		return nil, err
	}

	return &validation.CreateSignatureDeviceResponse{Status: "Device Created"}, nil
}

func (d *deviceService) CreateECCKeys(generator crypto.ECCGenerator, marshaler crypto.ECCMarshaler) ([]byte, error) {
	keyPair, err := generator.Generate()

	if err != nil {
		return nil, err
	}

	_, privateKeyBytes, err := marshaler.Encode(*keyPair)

	if err != nil {
		return nil, err
	}

	return privateKeyBytes, nil
}

func (d *deviceService) CreateRSAKeys(generator crypto.RSAGenerator, marshaler crypto.RSAMarshaler) ([]byte, error) {
	keyPair, err := generator.Generate()

	if err != nil {
		return nil, err
	}

	_, privateKeyBytes, err := marshaler.Marshal(*keyPair)

	if err != nil {
		return nil, err
	}

	return privateKeyBytes, nil
}

func (d *deviceService) ListSignatureDevice(req *validation.ListSignatureDeviceRequest) (*validation.ListSignatureDeviceResponse, error) {
	deviceList, err := d.collection.ListSignatureDevices(req.ID, req.Label, req.Algorithm)
	if err != nil {
		return nil, err
	}
	return &validation.ListSignatureDeviceResponse{Device: deviceList}, nil
}

func (d *deviceService) GetSignatureDevice(req *validation.GetSignatureDeviceRequest) (*validation.GetSignatureDeviceResponse, error) {
	device, err := d.collection.GetSignatureDevice(req.ID)
	if err != nil {
		return nil, err
	}
	return &validation.GetSignatureDeviceResponse{Device: device}, nil
}
