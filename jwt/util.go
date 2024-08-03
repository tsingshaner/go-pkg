package jwt

import (
	stdErrs "errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/tsingshaner/go-pkg/errors"
)

var (
	RESTErrTokenInvalid = errors.NewREST(http.StatusUnauthorized, "j_b5t0", "token invalid")
	RESTErrTokenExpired = errors.NewREST(http.StatusUnauthorized, "j_b5t1", "token expired")
	RESTErrUnknown      = errors.NewREST(http.StatusInternalServerError, "j_dwt1", "unknown error")
)

func BuildRESTError(e error) (ok bool, restErr error) {
	switch {
	case stdErrs.Is(e, jwt.ErrTokenExpired):
		return true, stdErrs.Join(RESTErrTokenExpired, e)
	case
		stdErrs.Is(e, jwt.ErrTokenMalformed),
		stdErrs.Is(e, jwt.ErrTokenUnverifiable),
		stdErrs.Is(e, jwt.ErrTokenSignatureInvalid),
		stdErrs.Is(e, jwt.ErrTokenInvalidClaims),
		stdErrs.Is(e, ErrAlgNotDefined),
		stdErrs.Is(e, ErrTokenInferFailed):
		return true, stdErrs.Join(RESTErrTokenInvalid, e)
	}

	return false, stdErrs.Join(RESTErrUnknown, e)
}
