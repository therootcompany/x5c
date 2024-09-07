package x5c

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"
	"time"
)

type Hex []byte

func (h Hex) MarshalJSON() ([]byte, error) {
	return json.Marshal(hex.EncodeToString(h))
}

// type Certificate struct {
// 	Raw                     Hex `json:"raw"`
// 	RawTBSCertificate       Hex `json:"raw_tbs_certificate"`
// 	RawSubjectPublicKeyInfo Hex `json:"raw_subject_public_key_info"`
// 	RawSubject              Hex `json:"raw_subject"`
// 	RawIssuer               Hex `json:"raw_issuer"`

// 	Signature          Hex                     `json:"signature"`
// 	SignatureAlgorithm x509.SignatureAlgorithm `json:"signature_algorithm"`
// 	PublicKeyAlgorithm x509.PublicKeyAlgorithm `json:"public_key_algorithm"`
// 	PublicKey          any                     `json:"public_key"`

// 	Version                     int                     `json:"version"`
// 	SerialNumber                *big.Int                `json:"serial_number"`
// 	Issuer                      pkix.Name               `json:"issuer"`
// 	Subject                     pkix.Name               `json:"subject"`
// 	NotBefore                   time.Time               `json:"not_before"`
// 	NotAfter                    time.Time               `json:"not_after"`
// 	KeyUsage                    x509.KeyUsage           `json:"key_usage"`
// 	Extensions                  []pkix.Extension        `json:"extensions"`
// 	ExtraExtensions             []pkix.Extension        `json:"extra_extensions"`
// 	UnhandledCriticalExtensions []asn1.ObjectIdentifier `json:"unhandled_critical_extensions"`

// 	ExtKeyUsage        []x509.ExtKeyUsage      `json:"ext_key_usage"`
// 	UnknownExtKeyUsage []asn1.ObjectIdentifier `json:"unknown_ext_key_usage"`

// 	BasicConstraintsValid bool `json:"basic_constraints_valid"`
// 	IsCA                  bool `json:"is_ca"`
// 	MaxPathLen            int  `json:"max_path_len"`
// 	MaxPathLenZero        bool `json:"max_path_len_zero"`

// 	SubjectKeyId   Hex `json:"subject_key_id"`
// 	AuthorityKeyId Hex `json:"authority_key_id"`

// 	OCSPServer            []string   `json:"ocsp_server"`
// 	IssuingCertificateURL []string   `json:"issuing_certificate_url"`
// 	DNSNames              []string   `json:"dns_names"`
// 	EmailAddresses        []string   `json:"email_addresses"`
// 	IPAddresses           []net.IP   `json:"ip_addresses"`
// 	URIs                  []*url.URL `json:"uris"`

// 	PermittedDNSDomainsCritical bool         `json:"permitted_dns_domains_critical"`
// 	PermittedDNSDomains         []string     `json:"permitted_dns_domains"`
// 	ExcludedDNSDomains          []string     `json:"excluded_dns_domains"`
// 	PermittedIPRanges           []*net.IPNet `json:"permitted_ip_ranges"`
// 	ExcludedIPRanges            []*net.IPNet `json:"excluded_ip_ranges"`
// 	PermittedEmailAddresses     []string     `json:"permitted_email_addresses"`
// 	ExcludedEmailAddresses      []string     `json:"excluded_email_addresses"`
// 	PermittedURIDomains         []string     `json:"permitted_uri_domains"`
// 	ExcludedURIDomains          []string     `json:"excluded_uri_domains"`

// 	CRLDistributionPoints []string `json:"crl_distribution_points"`

// 	PolicyIdentifiers []asn1.ObjectIdentifier `json:"policy_identifiers"`
// 	Policies          []x509.OID              `json:"policies"`
// }

type CertInfo struct {
	Issuer            string    `json:"issuer"`
	Subject           string    `json:"subject"`
	SerialNumber      string    `json:"serial_number"`
	ValidFrom         time.Time `json:"valid_from"`
	ValidTo           time.Time `json:"valid_to"`
	SHA1Fingerprint   string    `json:"sha1_fingerprint"`
	SHA256Fingerprint string    `json:"sha256_fingerprint"`
}

func MagicDecodeCertString(certParam string) ([]byte, error) {
	{
		// Unindent
		lines := strings.Split(certParam, "\n")
		var trimmed []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				trimmed = append(trimmed, line)
			}
		}
		certParam = strings.Join(trimmed, "\n")
	}

	if strings.Contains(certParam, "-----BEGIN CERTIFICATE-----") {
		block, _ := pem.Decode([]byte(certParam))
		if block == nil {
			return nil, fmt.Errorf("failed to parse PEM block")
		}
		return block.Bytes, nil
	}

	{
		// Remove all whitespace from hex, base64
		fields := strings.Fields(certParam)
		certParam = strings.Join(fields, "")
	}

	if certBytes, err := hex.DecodeString(certParam); err == nil {
		return certBytes, nil
	}

	{
		// URL Base64 to RFC, sans padding
		certParam = strings.ReplaceAll(certParam, "-", "+")
		certParam = strings.ReplaceAll(certParam, "_", "/")
		certParam = strings.TrimRight(certParam, "=")
	}

	if certBytes, err := base64.RawStdEncoding.DecodeString(certParam); err == nil {
		return certBytes, nil
	}

	msg := "failed to parse certificate string as PEM, Hex (DER), RFC Base64, or URL Base64 (whitespace removed)"
	return nil, fmt.Errorf(msg)
}

func Summarize(cert *x509.Certificate) *CertInfo {
	sha1Fingerprint := FingerprintSHA1(cert.Raw)
	sha256Fingerprint := FingerprintSHA256(cert.Raw)

	return &CertInfo{
		Issuer:            cert.Issuer.String(),
		Subject:           cert.Subject.String(),
		SerialNumber:      cert.SerialNumber.String(),
		ValidFrom:         cert.NotBefore,
		ValidTo:           cert.NotAfter,
		SHA1Fingerprint:   sha1Fingerprint,
		SHA256Fingerprint: sha256Fingerprint,
	}
}

func FingerprintSHA1(certBytes []byte) string {
	hasher := sha1.New()
	hasher.Write(certBytes)
	fingerprint := hasher.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(fingerprint))
}

func FingerprintSHA256(certBytes []byte) string {
	hasher := sha256.New()
	hasher.Write(certBytes)
	fingerprint := hasher.Sum(nil)

	return strings.ToUpper(hex.EncodeToString(fingerprint))
}
