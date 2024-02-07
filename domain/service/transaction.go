package service

import (
	"fmt"
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/model"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/repo"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/validation"
	"github.com/google/uuid"
)

type TransactionService interface {
	SignTransaction(req *validation.SignTransactionRequest) (*validation.SignTransactionResponse, error)
	ListTransaction(req *validation.ListTransactionRequest) (*validation.ListTransactionResponse, error)
	GetTransaction(req *validation.GetTransactionRequest) (*validation.GetTransactionResponse, error)
}

type transactionService struct {
	collection repo.Collection
	logger     *log.Logger
}

func NewTransactionService(logger *log.Logger, collection repo.Collection) TransactionService {
	return &transactionService{
		logger:     logger,
		collection: collection,
	}
}

func (t *transactionService) SignTransaction(req *validation.SignTransactionRequest) (*validation.SignTransactionResponse, error) {
	if err := req.IsValid(); err != nil {
		return nil, err
	}

	device, err := t.collection.GetSignatureDevice(req.DeviceID)
	if err != nil {
		return nil, err
	}

	signer, err := getSigner(device)
	if err != nil {
		return nil, err
	}

	signature, err := signer.Sign(req.Data)
	if err != nil {
		return nil, err
	}

	lastSignature, err := signer.Sign([]byte(req.DeviceID))
	if err != nil {
		return nil, err
	}

	if device.SignatureCounter > 0 {
		transactions, err := t.collection.ListTransactions(req.DeviceID)
		if err != nil {
			return nil, err
		}
		for _, transaction := range transactions {
			lastSignature = transaction.Data
		}
	}

	signatureModel := &model.Transaction{
		ID:               uuid.New().String(),
		DeviceID:         req.DeviceID,
		SignatureCounter: device.SignatureCounter,
		Data:             signature,
		LastSignature:    lastSignature,
	}

	_, err = t.collection.SignTransaction(signatureModel)
	if err != nil {
		return nil, err
	}

	return &validation.SignTransactionResponse{
		Transaction: string(signature),
		SignedData:  fmt.Sprintf("%d_%s_%s", device.SignatureCounter, req.Data, lastSignature),
	}, nil
}

func (t *transactionService) ListTransaction(req *validation.ListTransactionRequest) (*validation.ListTransactionResponse, error) {
	transactions, err := t.collection.ListTransactions(req.DeviceID)
	if err != nil {
		return nil, err
	}
	return &validation.ListTransactionResponse{Transaction: transactions}, nil
}

func (t *transactionService) GetTransaction(req *validation.GetTransactionRequest) (*validation.GetTransactionResponse, error) {
	transaction, err := t.collection.GetTransaction(req.ID)
	if err != nil {
		return nil, err
	}
	return &validation.GetTransactionResponse{Transaction: transaction}, nil
}

func getSigner(device *model.Device) (crypto.Signer, error) {
	switch device.Algorithm {
	case "ECC":
		return &crypto.ECCSigner{Device: device}, nil
	case "RSA":
		return &crypto.RSASigner{Device: device}, nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", device.Algorithm)
	}
}
