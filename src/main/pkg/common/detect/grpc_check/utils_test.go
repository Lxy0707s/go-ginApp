package grpc_check

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"
)

func generateCertificateTemplate(expiry time.Time, IPAddressSAN bool) *x509.Certificate {
	template := &x509.Certificate{
		BasicConstraintsValid: true,
		SubjectKeyId:          []byte{1},
		SerialNumber:          big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Example Org"},
		},
		NotBefore:   time.Now(),
		NotAfter:    expiry,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	template.DNSNames = append(template.DNSNames, "localhost")
	if IPAddressSAN {
		template.IPAddresses = append(template.IPAddresses, net.ParseIP("127.0.0.1"))
		template.IPAddresses = append(template.IPAddresses, net.ParseIP("::1"))
	}

	return template
}

func generateCertificate(template, parent *x509.Certificate, publickey *rsa.PublicKey, privatekey *rsa.PrivateKey) (*x509.Certificate, []byte) {
	derCert, err := x509.CreateCertificate(rand.Reader, template, template, publickey, privatekey)
	if err != nil {
		panic(fmt.Sprintf("Error signing test-certificate: %s", err))
	}
	cert, err := x509.ParseCertificate(derCert)
	if err != nil {
		panic(fmt.Sprintf("Error parsing test-certificate: %s", err))
	}
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derCert})
	return cert, pemCert

}

func generateSignedCertificate(template, parentCert *x509.Certificate, parentKey *rsa.PrivateKey) (*x509.Certificate, []byte, *rsa.PrivateKey) {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("Error creating rsa key: %s", err))
	}
	cert, pemCert := generateCertificate(template, parentCert, &privatekey.PublicKey, parentKey)
	return cert, pemCert, privatekey
}

func generateSelfSignedCertificate(template *x509.Certificate) (*x509.Certificate, []byte, *rsa.PrivateKey) {
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Sprintf("Error creating rsa key: %s", err))
	}
	publickey := &privatekey.PublicKey

	cert, pemCert := generateCertificate(template, template, publickey, privatekey)
	return cert, pemCert, privatekey
}

func generateSelfSignedCertificateWithPrivateKey(template *x509.Certificate, privatekey *rsa.PrivateKey) (*x509.Certificate, []byte) {
	publickey := &privatekey.PublicKey
	cert, pemCert := generateCertificate(template, template, publickey, privatekey)
	return cert, pemCert
}
