package storage

import (
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/model"
)

type Storage struct {
	Device          map[string]*model.Device
	DeviceRWLock    sync.RWMutex
	Transaction     map[string]*model.Transaction
	SignatureRWLock sync.RWMutex
}

func NewStorage() *Storage {
	deviceMap := make(map[string]*model.Device, 0)
	signatureMap := make(map[string]*model.Transaction, 0)
	return &Storage{Device: deviceMap, Transaction: signatureMap}
}
