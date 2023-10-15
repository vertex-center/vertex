package types

type PublicKey struct {
	Type              string `json:"type"`
	FingerprintSHA256 string `json:"fingerprint_sha_256"`
}
