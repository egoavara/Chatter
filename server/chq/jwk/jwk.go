package jwk

import (
	"crypto/x509"
	"net/url"
)

type Set struct {
	Keys []*Key
}
type Key struct {
	KeyType                KeyType             // https://tools.ietf.org/html/rfc7517#section-4.1
	KeyUse                 KeyUse              // https://tools.ietf.org/html/rfc7517#section-4.2
	KeyOperations          map[KeyOp]struct{}  // https://tools.ietf.org/html/rfc7517#section-4.3
	Algorithm              Algorithm           // https://tools.ietf.org/html/rfc7517#section-4.4
	KeyID                  string              // https://tools.ietf.org/html/rfc7515#section-4.5
	X509URL                *url.URL            // https://tools.ietf.org/html/rfc7515#section-4.6
	X509CertChain          []*x509.Certificate // https://tools.ietf.org/html/rfc7515#section-4.7
	X509CertThumbprint     string              // https://tools.ietf.org/html/rfc7515#section-4.8
	X509CertThumbprintS256 string              // https://tools.ietf.org/html/rfc7515#section-4.9
	// Go language stdlib crypto
	Raw interface{}
}

func (set *Set) First() *Key {
	if len(set.Keys) > 0 {
		return set.Keys[0]
	}
	return nil
}
