package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tsingshaner/go-pkg/util"
)

var (
	ErrAlgNotDefined    = errors.New("alg not defined")
	ErrTokenInferFailed = errors.New("token infer failed")
	ErrTokenInvalid     = errors.New("token if invalid")
)

type (
	TokenMeta struct {
		KeyMap KeyMap
		claims ClaimsOption
	}

	Claims           = jwt.Claims
	RegisteredClaims = jwt.RegisteredClaims
	Token            = jwt.Token
	TokenOption      = jwt.TokenOption
)

func (tm *TokenMeta) GetKeyByAlg(alg string) (any, error) {
	if key, ok := tm.KeyMap[alg]; ok {
		return key, nil
	}

	return nil, ErrJWAKeyNotFound
}

func (tm *TokenMeta) SignedWithClaims(alg ALG, claims Claims, fns ...TokenOption) (string, error) {
	key, ok := tm.KeyMap[string(alg)]
	if !ok {
		return "", fmt.Errorf("(%s): %w", alg, ErrMethodNotFound)
	}

	return jwt.NewWithClaims(key.Method, claims, fns...).SignedString(key.PrivateKey)
}

func (tm *TokenMeta) NewRegisteredClaims(isRefresh bool) *RegisteredClaims {
	now := time.Now()

	return &RegisteredClaims{
		Issuer:    tm.claims.Issuer,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Audience:  jwt.ClaimStrings(tm.claims.Audience),
		ExpiresAt: jwt.NewNumericDate(now.Add(
			util.Pick(isRefresh, tm.claims.RefreshExpire, tm.claims.Expire)),
		),
	}
}

func (tm *TokenMeta) ParseWithClaims(token string, claims Claims) (*Token, error) {
	return jwt.ParseWithClaims(token, claims, func(token *Token) (any, error) {
		if key, ok := tm.KeyMap[token.Method.Alg()]; ok {
			return key.PublicKey, nil
		} else {
			return nil, fmt.Errorf("(%s): %w", token.Method.Alg(), ErrMethodNotFound)
		}
	})
}

func ParseWithClaims[T Claims](tm *TokenMeta, token string, claims Claims) (T, error) {
	if t, err := tm.ParseWithClaims(token, claims); err != nil {
		return *new(T), err
	} else if c, ok := t.Claims.(T); ok {
		if !t.Valid {
			return *new(T), ErrTokenInvalid
		}
		return c, nil
	} else {
		return *new(T), ErrTokenInferFailed
	}
}
