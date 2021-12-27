package jwk

type KeyUse string

const (
	KeyUseNo  KeyUse = ""
	KeyUseSig KeyUse = "sig"
	KeyUseEnc KeyUse = "enc"
)
