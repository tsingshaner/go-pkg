package jwt

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tsingshaner/go-pkg/log/console"
)

var (
	ErrJWAKeyNotFound  = errors.New("jwa key not found")
	ErrMethodNotFound  = errors.New("jwt method not found")
	ErrKeyRequired     = errors.New("key required")
	ErrPEMPathRequired = errors.New("required a public and private pem path")
	ErrPEMLoadFailed   = errors.New("pem file load failed")
	ErrPEMParseFailed  = errors.New("pem file parse failed")
)

type ALG string

const (
	HS256 ALG = "HS256"
	HS384 ALG = "HS384"
	HS512 ALG = "HS512"

	EdDSA ALG = "EdDSA"

	ES256 ALG = "ES256"
	ES384 ALG = "ES384"
	ES512 ALG = "ES512"

	RS256 ALG = "RS256"
	RS384 ALG = "RS384"
	RS512 ALG = "RS512"
)

type KeyItem struct {
	Method     jwt.SigningMethod
	PrivateKey any
	PublicKey  any
}

type KeyMap map[string]*KeyItem

func (opts *Options) LoadKeys() (KeyMap, error) {
	keys := make(KeyMap, len(opts.Methods))

	for _, method := range opts.Methods {
		if item, err := method.build(); err != nil {
			return nil, fmt.Errorf("(%s) %w", method.Alg, err)
		} else {
			console.Info("(%s) load key success", method.Alg)
			keys[method.Alg] = item
		}
	}

	return keys, nil
}

func (j JWA) build() (*KeyItem, error) {
	item := &KeyItem{}

	if item.Method = jwt.GetSigningMethod(j.Alg); item.Method == nil {
		return nil, ErrMethodNotFound
	}

	if strings.HasPrefix(j.Alg, "HS") {
		if j.Key == "" {
			return nil, ErrKeyRequired
		}

		item.PrivateKey = []byte(j.Key)
		item.PublicKey = []byte(j.Key)
		return item, nil
	}

	if j.Pem.PrivatePath == "" || j.Pem.PublicPath == "" {
		return nil, ErrPEMPathRequired
	}

	var err error
	switch ALG(j.Alg) {
	case ES256, ES384, ES512:
		item.PrivateKey, item.PublicKey, err = loadPEM(j.Pem,
			jwt.ParseECPrivateKeyFromPEM, jwt.ParseECPublicKeyFromPEM,
		)
	case RS256, RS384, RS512:
		item.PrivateKey, item.PublicKey, err = loadPEM(j.Pem,
			jwt.ParseRSAPublicKeyFromPEM, jwt.ParseRSAPrivateKeyFromPEM,
		)
	case EdDSA:
		item.PrivateKey, item.PublicKey, err = loadPEM(j.Pem,
			jwt.ParseEdPublicKeyFromPEM, jwt.ParseEdPrivateKeyFromPEM,
		)
	default:
		err = ErrMethodNotFound
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

type parser[T any] func([]byte) (T, error)

func loadPEM[Pub, Pri any](pem PEM, pubParser parser[Pub], priParser parser[Pri]) (Pub, Pri, error) {
	pubKey, err := parsePEM(pem.PublicPath, pubParser)
	if err != nil {
		return *new(Pub), *new(Pri), fmt.Errorf("public %w", err)
	}

	priKey, err := parsePEM(pem.PrivatePath, priParser)
	if err != nil {
		return *new(Pub), *new(Pri), fmt.Errorf("private %w", err)
	}

	return pubKey, priKey, nil
}

func parsePEM[T any](pemPath string, parser parser[T]) (T, error) {
	if pem, err := os.ReadFile(pemPath); err != nil {
		return *new(T), errors.Join(ErrPEMLoadFailed, err)
	} else if key, err := parser(pem); err != nil {
		return *new(T), errors.Join(ErrPEMParseFailed, err)
	} else {
		return key, nil
	}
}
