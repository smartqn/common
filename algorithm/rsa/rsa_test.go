package rsa

import (
	"encoding/hex"
	"testing"

	"github.com/linnv/logx"
)

func TestGenerateKeyPair(t *testing.T) {
	// pri, pub := GenerateKeyPair(2048)
	pri, pub := GenerateKeyPair(2001)

	pubStr := PublicKeyToBytes(pub)
	logx.Debugf("pubStr: [%s]\n", pubStr)
	priStr := PrivateKeyToBytes(pri)
	logx.Debugf("priStr: [%s]\n", priStr)
	rawMsg := []byte(`{"aKey":"111"}`)

	encryptMsg := EncryptWithPublicKey(rawMsg, pub)
	logx.Debugf("encryptMsg: %s\n[%s]\n", encryptMsg, hex.EncodeToString(encryptMsg))
	msg, _ := hex.DecodeString(hex.EncodeToString(encryptMsg))
	logx.Debugf("hex.DecodeString(): %s\n", msg)

	decryptMsg := DecryptWithPrivateKey(encryptMsg, pri)
	logx.Debugf("decryptMsg: %s\n", decryptMsg)
	logx.Debugf("rawMsg: %s\n", rawMsg)

	getPri := BytesToPrivateKey(PrivateKey)
	getPub := BytesToPublicKey(PublicKey)

	getEncryptMsg := EncryptWithPublicKey(rawMsg, getPub)
	dd := hex.EncodeToString(getEncryptMsg)
	logx.Debugf("dd: %s\n", dd)
	getDecryptMsg := DecryptWithPrivateKey(getEncryptMsg, getPri)
	logx.Debugf("getEncryptMsg: %s\n", getEncryptMsg)
	logx.Debugf("getDecryptMsg: %s\n", getDecryptMsg)

	// pubStr := PublicKeyToBytes(pub)
	// priStr := PrivateKeyToBytes(pri)
	// logx.Debugf("pubStr: %s\n", pubStr)
	// logx.Debugf("priStr: %s\n", priStr)
	// type args struct {
	// 	bits int
	// }
	// tests := []struct {
	// 	name  string
	// 	args  args
	// 	want  *rsa.PrivateKey
	// 	want1 *rsa.PublicKey
	// }{
	// 	// TODO: Add test cases.
	// }
	// for _, tt := range tests {
	// 	t.Run(tt.name, func(t *testing.T) {
	// 		got, got1 := GenerateKeyPair(tt.args.bits)
	// 		if !reflect.DeepEqual(got, tt.want) {
	// 			t.Errorf("GenerateKeyPair() got = %v, want %v", got, tt.want)
	// 		}
	// 		if !reflect.DeepEqual(got1, tt.want1) {
	// 			t.Errorf("GenerateKeyPair() got1 = %v, want %v", got1, tt.want1)
	// 		}
	// 	})
	// }
}
