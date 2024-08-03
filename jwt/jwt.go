package jwt

import "time"

type (
	PEM struct {
		PrivatePath string `mapstructure:"privatePath"`
		PublicPath  string `mapstructure:"publicPath"`
	}

	JWA struct {
		Alg string `mapstructure:"alg"`
		Key string `mapstructure:"key"`
		Pem PEM    `mapstructure:"pem"`
	}

	ClaimsOption struct {
		Audience      []string      `mapstructure:"audience"`
		Issuer        string        `mapstructure:"issuer"`
		Subject       string        `mapstructure:"subject"`
		Expire        time.Duration `mapstructure:"expire"`
		RefreshExpire time.Duration `mapstructure:"refreshExpire"`
	}

	Options struct {
		Claims  ClaimsOption `mapstructure:"claims"`
		Methods []JWA        `mapstructure:"methods"`
	}
)

func New(opts *Options) (*TokenMeta, error) {
	tm := &TokenMeta{
		claims: opts.Claims,
	}

	if km, err := opts.LoadKeys(); err != nil {
		return nil, err
	} else {
		tm.KeyMap = km
	}

	return tm, nil
}
