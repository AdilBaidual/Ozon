package paseto

import (
	"encoding/hex"
	"time"

	"github.com/o1egl/paseto"
)

type Paseto struct {
	symmetricKey []byte
	paseto       *paseto.V2
}

func NewPaseto(symmetricKey string) (*Paseto, error) {
	key, err := hex.DecodeString(symmetricKey)
	if err != nil {
		return nil, err
	}

	return &Paseto{
		symmetricKey: key,
		paseto:       paseto.NewV2(),
	}, nil
}

func (p *Paseto) GenerateAccessToken(uuid string) (string, error) {
	now := time.Now()
	exp := now.Add(AccessTTL)
	nbt := now

	jsonToken := paseto.JSONToken{
		Jti:        "access",
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	jsonToken.Set("uuid", uuid)

	token, err := p.paseto.Encrypt(p.symmetricKey, jsonToken, nil)
	return token, err
}

func (p *Paseto) GenerateRefreshToken(uuid string) (string, error) {
	now := time.Now()
	exp := now.Add(RefreshTTL)
	nbt := now

	jsonToken := paseto.JSONToken{
		Jti:        "refresh",
		IssuedAt:   now,
		Expiration: exp,
		NotBefore:  nbt,
	}
	jsonToken.Set("uuid", uuid)

	token, err := p.paseto.Encrypt(p.symmetricKey, jsonToken, nil)
	return token, err
}

func (p *Paseto) ValidateToken(token string) (string, error) {
	payload := &paseto.JSONToken{}

	err := p.paseto.Decrypt(token, p.symmetricKey, payload, nil)
	if err != nil {
		return "", err
	}

	uuid := payload.Get("uuid")
	return uuid, err
}
