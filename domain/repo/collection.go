package repo

import (
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain/model"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/storage"
)

type collection struct {
	storage *storage.Storage
}

func NewRepository(db *storage.Storage) *collection {
	return &collection{storage: db}
}

type Collection interface {
	CreateSignatureDevice(device *model.Device) (*model.Device, error)
	GetSignatureDevice(id string) (*model.Device, error)
	GetTransaction(id string) (*model.Transaction, error)
	ListSignatureDevices(id, label, algorithm string) ([]*model.Device, error)
	ListTransactions(deviceID string) ([]*model.Transaction, error)
	SignTransaction(signature *model.Transaction) (*model.Transaction, error)
}

func (c *collection) incrementCounter(id string) error {
	_, exists := c.storage.Device[id]
	if !exists {
		return errors.New("device is missing")
	}
	c.storage.Device[id].SignatureCounter += 1
	return nil
}

func (c *collection) CreateSignatureDevice(device *model.Device) (*model.Device, error) {

	c.storage.DeviceRWLock.Lock()
	defer c.storage.DeviceRWLock.Unlock()

	if _, exists := c.storage.Device[device.ID]; exists {
		return nil, errors.New("device exists already")
	}

	c.storage.Device[device.ID] = device
	return device, nil
}

func (c *collection) ListSignatureDevices(id, label, algorithm string) ([]*model.Device, error) {
	c.storage.DeviceRWLock.RLock()
	defer c.storage.DeviceRWLock.RUnlock()

	var devices []*model.Device

	for _, device := range c.storage.Device {
		if (id == "" || device.ID == id) && (label == "" || device.Label == label) && (algorithm == "" || device.Algorithm == algorithm) {
			devices = append(devices, device)
		}
	}

	return devices, nil
}

func (c *collection) GetSignatureDevice(id string) (*model.Device, error) {
	c.storage.DeviceRWLock.RLock()
	defer c.storage.DeviceRWLock.RUnlock()

	device, exists := c.storage.Device[id]
	if !exists {
		return nil, errors.New("device not found")
	}

	return device, nil
}

func (c *collection) SignTransaction(transaction *model.Transaction) (*model.Transaction, error) {
	c.storage.SignatureRWLock.Lock()
	defer c.storage.SignatureRWLock.Unlock()

	if _, exists := c.storage.Transaction[transaction.ID]; exists {
		return nil, errors.New("transaction already exists")
	}

	c.storage.Transaction[transaction.ID] = transaction
	c.incrementCounter(transaction.DeviceID)
	return transaction, nil
}

func (c *collection) ListTransactions(deviceID string) ([]*model.Transaction, error) {
	c.storage.SignatureRWLock.RLock()
	defer c.storage.SignatureRWLock.RUnlock()

	var signatures []*model.Transaction

	for _, transaction := range c.storage.Transaction {
		if deviceID == "" || transaction.DeviceID == deviceID {
			signatures = append(signatures, transaction)
		}
	}

	return signatures, nil
}

func (c *collection) GetTransaction(id string) (*model.Transaction, error) {
	c.storage.SignatureRWLock.RLock()
	defer c.storage.SignatureRWLock.RUnlock()

	signature, exists := c.storage.Transaction[id]
	if !exists {
		return nil, errors.New("transaction is not found")
	}

	return signature, nil
}
