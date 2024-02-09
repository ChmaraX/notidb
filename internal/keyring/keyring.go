package keyring

import (
	"github.com/99designs/keyring"
)

const KeyringServiceName = "notidb"
const KeychainKey = "notidb"

type KeyringManager struct {
	ring keyring.Keyring
}

func NewKeyringManager() (*KeyringManager, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: KeyringServiceName,
		AllowedBackends: []keyring.BackendType{
			keyring.KeychainBackend,
			keyring.SecretServiceBackend,
		},
		KeychainName: "login",
	})

	if err != nil {
		return nil, err
	}

	return &KeyringManager{ring: ring}, nil
}

func (a *KeyringManager) SaveAPIKey(key string) error {
	err := a.ring.Set(keyring.Item{
		Key:  KeychainKey,
		Data: []byte(key),
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *KeyringManager) GetAPIKey() (string, error) {
	item, err := a.ring.Get(KeychainKey)
	if err != nil {
		return "", err
	}

	return string(item.Data), nil
}
