package thriftserver

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

func LoadX509KeyPair(certFile, keyFile, password string) (tls.Certificate, error) {
	certPEMByte, err := ioutil.ReadFile(certFile)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyPEMByte, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return tls.Certificate{}, err
	}

	keyPEMBlock, rest := pem.Decode(keyPEMByte)
	if len(rest) > 0 {
		return tls.Certificate{}, errors.New("Decode key failed")
	}

	if x509.IsEncryptedPEMBlock(keyPEMBlock) {
		keyDePEMByte, err := x509.DecryptPEMBlock(keyPEMBlock, []byte(password))
		if err != nil {
			return tls.Certificate{}, err
		}

		// 解析其中的private key
		key, err := x509.ParsePKCS1PrivateKey(keyDePEMByte)
		if err != nil {
			return tls.Certificate{}, err
		}
		// 编码成新的PEM结构
		keyNewPemByte := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})

		return tls.X509KeyPair(certPEMByte, keyNewPemByte)
	} else {
		return tls.X509KeyPair(certPEMByte, keyPEMByte)
	}
}
