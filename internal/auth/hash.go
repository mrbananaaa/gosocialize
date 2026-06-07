package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidHash         = errors.New("invalid hash format")
	ErrIncompatibleVersion = errors.New("incompatible version")
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(password, encoded string) (bool, error)
	NeedsUpgrade(encoded string) (bool, error)
}

type Argon2Config struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
	Version     uint32
}

type Argon2Hasher struct {
	config Argon2Config
}

func NewArgon2Hasher() *Argon2Hasher {
	return &Argon2Hasher{
		config: Argon2Config{
			Memory:      64 * 1024,
			Iterations:  4,
			Parallelism: 4,
			SaltLength:  16,
			KeyLength:   64,
			Version:     1,
		},
	}
}

func (a *Argon2Hasher) Hash(password string) (string, error) {
	salt, err := generateRandomBytes(a.config.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		a.config.Iterations,
		a.config.Memory,
		a.config.Parallelism,
		a.config.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		a.config.Version,
		a.config.Memory,
		a.config.Iterations,
		a.config.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encoded, nil
}

func (a *Argon2Hasher) Compare(password, encoded string) (bool, error) {
	cfg, salt, hash, err := decodeHash(encoded)
	if err != nil {
		return false, err
	}

	computed := argon2.IDKey(
		[]byte(password),
		salt,
		cfg.Iterations,
		cfg.Memory,
		cfg.Parallelism,
		cfg.KeyLength,
	)

	if subtle.ConstantTimeCompare(hash, computed) == 1 {
		return true, nil
	}

	return false, nil
}

func (a *Argon2Hasher) NeedsUpgrade(encoded string) (bool, error) {
	cfg, _, _, err := decodeHash(encoded)
	if err != nil {
		return false, err
	}

	if cfg.Version != a.config.Version {
		return true, nil
	}

	if cfg.Memory != a.config.Memory ||
		cfg.Iterations != a.config.Iterations ||
		cfg.Parallelism != a.config.Parallelism ||
		cfg.KeyLength != a.config.KeyLength {
		return true, nil
	}

	return false, nil
}

func decodeHash(encoded string) (*Argon2Config, []byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version uint32
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}

	if version != 1 {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	cfg := &Argon2Config{}
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d",
		&cfg.Memory,
		&cfg.Iterations,
		&cfg.Parallelism,
	)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}

	cfg.KeyLength = uint32(len(hash))
	cfg.SaltLength = uint32(len(salt))
	cfg.Version = version

	return cfg, salt, hash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
