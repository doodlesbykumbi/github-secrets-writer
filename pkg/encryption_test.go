package gsw

import (
	"bytes"
	"crypto/rand"
	"testing"

	"golang.org/x/crypto/nacl/box"
)

func Test_encryptValue(t *testing.T) {
	message := []byte("Attack at dawn.")

	publicKey, privateKey, _ := box.GenerateKey(rand.Reader)

	cryptedMessage, err := encryptValue(message, publicKey)
	if err != nil {
		t.Errorf("failed to encrypt message: %s", err)
	}
	decryptedMessage, err := decryptValue(cryptedMessage, publicKey, privateKey)
	if err != nil {
		t.Errorf("failed to decrypt message: %s", err)
	}

	if bytes.Compare(message, decryptedMessage) != 0 {
		t.Errorf("decrypted message not equal original message")
	}
}
