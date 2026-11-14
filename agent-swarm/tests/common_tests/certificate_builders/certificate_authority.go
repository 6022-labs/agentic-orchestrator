package certificate_builders

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"time"
)

// CertificateAuthority is a throwaway private CA that issues server certificates for TLS tests.
type CertificateAuthority struct {
	CertificatePem []byte
	certificate    *x509.Certificate
	key            *ecdsa.PrivateKey
}

func NewCertificateAuthority() *CertificateAuthority {
	ca := &CertificateAuthority{}
	ca.key = ca.newKey()
	template := &x509.Certificate{
		SerialNumber:          ca.newSerial(),
		Subject:               pkix.Name{CommonName: "test-ca"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &ca.key.PublicKey, ca.key)
	if err != nil {
		log.Fatalf("failed to create CA certificate: %v", err)
	}
	if ca.certificate, err = x509.ParseCertificate(der); err != nil {
		log.Fatalf("failed to parse CA certificate: %v", err)
	}
	ca.CertificatePem = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	return ca
}

// IssueServerCertificate returns a PEM certificate and PEM PKCS#8 key valid for dnsNames.
func (ca *CertificateAuthority) IssueServerCertificate(dnsNames ...string) ([]byte, []byte) {
	key := ca.newKey()
	template := &x509.Certificate{
		SerialNumber: ca.newSerial(),
		Subject:      pkix.Name{CommonName: dnsNames[0]},
		DNSNames:     dnsNames,
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, template, ca.certificate, &key.PublicKey, ca.key)
	if err != nil {
		log.Fatalf("failed to create server certificate: %v", err)
	}
	keyDer, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		log.Fatalf("failed to marshal server key: %v", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}),
		pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDer})
}

func (ca *CertificateAuthority) newKey() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Fatalf("failed to generate key: %v", err)
	}
	return key
}

func (ca *CertificateAuthority) newSerial() *big.Int {
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 62))
	if err != nil {
		log.Fatalf("failed to generate serial: %v", err)
	}
	return serial
}
