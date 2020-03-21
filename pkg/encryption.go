package gsw

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/box"
)

const publicKeyLength = 32
const nonceLength = 24

func createNonce(ePublicKey, publicKey *[publicKeyLength]byte) *[nonceLength]byte {
	h, _ := blake2b.New(nonceLength, nil)
	h.Write(ePublicKey[:])
	h.Write(publicKey[:])
	var nonce = &[nonceLength]byte{}
	copy(nonce[:], h.Sum(nil))

	return nonce
}

func encryptValue(value []byte, publicKey *[publicKeyLength]byte) ([]byte, error) {
	ePublicKey, ePrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	crypted := make([]byte, 0, box.Overhead+len(value))

	nonce := createNonce(ePublicKey, publicKey)

	_ = box.Seal(crypted, value, nonce, publicKey, ePrivateKey)

	return append(ePublicKey[:], crypted[0:cap(crypted)]...), nil
}

func decryptValue(
	encryptedValue []byte,
	publicKey, privateKey *[publicKeyLength]byte,
) ([]byte, error) {
	var ePublicKey = &[publicKeyLength]byte{}
	copy(ePublicKey[:], encryptedValue[:publicKeyLength])

	nonce := createNonce(ePublicKey, publicKey)

	output := make([]byte, 0, len(encryptedValue[publicKeyLength:])-box.Overhead)
	_, ok := box.Open(output, encryptedValue[publicKeyLength:], nonce, ePublicKey, privateKey)
	if !ok {
		return nil, fmt.Errorf("unable to decrypt value")
	}

	return output[0:cap(output)], nil
}
