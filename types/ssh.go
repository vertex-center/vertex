package types

type PublicKey struct {
	Type              string `json:"type"`
	FingerprintSHA256 string `json:"fingerprint_sha_256"`
}

type SshAdapterPort interface {
	GetAll() ([]PublicKey, error)
	Add(key string) error
	Remove(fingerprint string) error
}
