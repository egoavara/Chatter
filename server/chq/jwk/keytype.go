package jwk

type KeyType string

const (
	KeyTypeEC    KeyType = "EC"  // Elliptic Curve
	KeyTypeOctet KeyType = "oct" // Octet sequence
	KeyTypeRSA   KeyType = "RSA" // RSA
)
