package jwk

type KeyOp string

const (
	KeyOpSign       KeyOp = "sign"       // (compute digital signature or MAC)
	KeyOpVerify     KeyOp = "verify"     // (verify digital signature or MAC)
	KeyOpEncrypt    KeyOp = "encrypt"    // (encrypt content)
	KeyOpDecrypt    KeyOp = "decrypt"    // (decrypt content and validate decryption, if applicable)
	KeyOpWrapKey    KeyOp = "wrapKey"    // (encrypt key)
	KeyOpUnwrapKey  KeyOp = "unwrapKey"  // (decrypt key and validate decryption, if applicable)
	KeyOpDeriveKey  KeyOp = "deriveKey"  // (derive key)
	KeyOpDeriveBits KeyOp = "deriveBits" // (derive bits not to be used as a key)
)
