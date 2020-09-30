package proxyServer

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

type СertSettings struct {
	Folder       string
	CertName     string
	RootCertFile string
	RootKeyFile  string
}

func (cs *СertSettings) LoadCA() (cert tls.Certificate, err error) {
	cert, err = tls.LoadX509KeyPair(cs.RootCertFile, cs.RootKeyFile)
	if os.IsNotExist(err) {
		cert, err = cs.CallGenCAAndWriteToFile()
	}
	if err == nil {
		cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	}
	return
}

func (cs *СertSettings) CallGenCAAndWriteToFile() (cert tls.Certificate, err error) {
	err = os.MkdirAll(cs.Folder, 0700)
	if err != nil {
		return
	}
	certPEM, keyPEM, err := GenCA(cs.CertName)
	if err != nil {
		return
	}
	cert, _ = tls.X509KeyPair(certPEM, keyPEM)
	err = ioutil.WriteFile(cs.RootCertFile, certPEM, 0400)
	if err == nil {
		err = ioutil.WriteFile(cs.RootKeyFile, keyPEM, 0400)
	}
	return cert, err
}

func GenCA(name string) (certPEM, keyPEM []byte, err error) {
	now := time.Now().UTC()
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             now,
		NotAfter:              now.Add(caMaxAge),
		KeyUsage:              caUsage,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
		SignatureAlgorithm:    x509.ECDSAWithSHA512, //ECDSAWithSHA512
	}
	key, err := genKeyPair()
	if err != nil {
		return
	}
	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, key.Public(), key)
	if err != nil {
		return
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return
	}
	certPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})
	keyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "ECDSA PRIVATE KEY",
		Bytes: keyDER,
	})
	return
}
