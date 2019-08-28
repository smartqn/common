package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"log"

	"github.com/linnv/logx"
)

var label = []byte("qnzs")

//@TODO handle error
func CrypteMsg(rawMsg []byte) string {
	pub := BytesToPublicKey(PublicKey)
	getEncryptMsg := EncryptWithPublicKey(rawMsg, pub)
	ret := hex.EncodeToString(getEncryptMsg)
	return ret
}

func DeCrypteMsg(rawMsg string) []byte {
	encryptMsg, err := hex.DecodeString(rawMsg)
	if err != nil {
		return nil
	}
	pri := BytesToPrivateKey(PrivateKey)
	getDecryptMsg := DecryptWithPrivateKey(encryptMsg, pri)
	return getDecryptMsg
}

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			// Type:  "RSA PRIVATE KEY",
			Type:  "QNZS PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "QNZS PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			logx.Warnf("err: %+v\n", err)
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
	}
	return key
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		log.Println("is encrypted pem block")
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			logx.Warnf("err: %+v\n", err)
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		logx.Warnf("key convert not ok")
	}
	return key
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) []byte {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, label)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
	}
	return ciphertext
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) []byte {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, label)
	if err != nil {
		logx.Warnf("err: %+v\n", err)
	}
	return plaintext
}
