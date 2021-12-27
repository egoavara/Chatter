package jwk

type Algorithm string

const (
	AlgorithmNo    Algorithm = ""
	AlgorithmHS256 Algorithm = "HS256"
	AlgorithmHS384 Algorithm = "HS384"
	AlgorithmHS512 Algorithm = "HS512"
	AlgorithmRS256 Algorithm = "RS256"
	AlgorithmRS384 Algorithm = "RS384"
	AlgorithmRS512 Algorithm = "RS512"
	AlgorithmES256 Algorithm = "ES256"
	AlgorithmES384 Algorithm = "ES384"
	AlgorithmES512 Algorithm = "ES512"
	AlgorithmPS256 Algorithm = "PS256"
	AlgorithmPS384 Algorithm = "PS384"
	AlgorithmPS512 Algorithm = "PS512"
)
