package gen

import (
	_ "embed"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"

	"github.com/tsingshaner/go-pkg/conf"
	"github.com/tsingshaner/go-pkg/log/console"
)

//go:embed errs.go.template
var ErrorPkgTemplateStr string
var ErrorPkgTemplate *template.Template

type Error struct {
	Key string `mapstructure:"key"`
	Msg string `mapstructure:"msg"`
}

type ErrorMap map[string]Error

type ErrorConfig struct {
	File     string              `mapstructure:"file"`
	Package  string              `mapstructure:"pkg"`
	ModCode  string              `mapstructure:"modCode"`
	BasicErr ErrorMap            `mapstructure:"basic"`
	RestErr  map[string]ErrorMap `mapstructure:"rest"`
}

func init() {
	var err error
	if ErrorPkgTemplate, err = template.New("ErrorPkg").Parse(ErrorPkgTemplateStr); err != nil {
		console.Fatal("%+v", err)
	}
}

func GeneratePkg(c *ErrorConfig) {
	file, err := os.OpenFile(c.File, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		console.Fatal("open file (%s) err %+v", c.File, err)
	}

	if err := ErrorPkgTemplate.Execute(file, GenErrorPkgData(c)); err != nil {
		println()
		console.Fatal("%+v", err)
	}

	console.Info("generate error pkg %s success", c.File)
}

func ReadErrors() *ErrorConfig {
	r := conf.New(&ErrorConfig{}, conf.ParseArgs())
	if err := r.Load(); err != nil {
		console.Fatal("load error config err\n%+v", err)
	}

	return r.Value
}

type ErrorObj struct {
	Error
	Code string
}

type RestErrorObj struct {
	ErrorObj
	Status string
}

type ErrorPkgTemplateData struct {
	Package   string
	BasicErrs []ErrorObj
	RestErrs  map[string][]RestErrorObj
}

func GenErrorPkgData(conf *ErrorConfig) *ErrorPkgTemplateData {
	data := &ErrorPkgTemplateData{
		Package:  conf.Package,
		RestErrs: make(map[string][]RestErrorObj),
	}

	for code, basic := range conf.BasicErr {
		data.BasicErrs = append(data.BasicErrs, ErrorObj{
			Code:  conf.ModCode + code,
			Error: basic,
		})
	}

	for status, errs := range conf.RestErr {
		statusStr := TranslateStatus(status)

		for code, err := range errs {
			data.RestErrs[statusStr] = append(data.RestErrs[statusStr], RestErrorObj{
				Status: TranslateStatus(status),
				ErrorObj: ErrorObj{
					Code:  conf.ModCode + strconv.Itoa(HttpStatus[statusStr]) + code,
					Error: err,
				},
			})
		}
	}

	return data
}

var HttpStatus = map[string]int{
	"Continue":                      http.StatusContinue,
	"SwitchingProtocols":            http.StatusSwitchingProtocols,
	"Processing":                    http.StatusProcessing,
	"EarlyHints":                    http.StatusEarlyHints,
	"OK":                            http.StatusOK,
	"Created":                       http.StatusCreated,
	"Accepted":                      http.StatusAccepted,
	"NonAuthoritativeInfo":          http.StatusNonAuthoritativeInfo,
	"NoContent":                     http.StatusNoContent,
	"ResetContent":                  http.StatusResetContent,
	"PartialContent":                http.StatusPartialContent,
	"MultiStatus":                   http.StatusMultiStatus,
	"AlreadyReported":               http.StatusAlreadyReported,
	"IMUsed":                        http.StatusIMUsed,
	"MultipleChoices":               http.StatusMultipleChoices,
	"MovedPermanently":              http.StatusMovedPermanently,
	"Found":                         http.StatusFound,
	"SeeOther":                      http.StatusSeeOther,
	"NotModified":                   http.StatusNotModified,
	"UseProxy":                      http.StatusUseProxy,
	"TemporaryRedirect":             http.StatusTemporaryRedirect,
	"PermanentRedirect":             http.StatusPermanentRedirect,
	"BadRequest":                    http.StatusBadRequest,
	"Unauthorized":                  http.StatusUnauthorized,
	"PaymentRequired":               http.StatusPaymentRequired,
	"Forbidden":                     http.StatusForbidden,
	"NotFound":                      http.StatusNotFound,
	"MethodNotAllowed":              http.StatusMethodNotAllowed,
	"NotAcceptable":                 http.StatusNotAcceptable,
	"ProxyAuthRequired":             http.StatusProxyAuthRequired,
	"RequestTimeout":                http.StatusRequestTimeout,
	"Conflict":                      http.StatusConflict,
	"Gone":                          http.StatusGone,
	"LengthRequired":                http.StatusLengthRequired,
	"PreconditionFailed":            http.StatusPreconditionFailed,
	"RequestEntityTooLarge":         http.StatusRequestEntityTooLarge,
	"RequestURITooLong":             http.StatusRequestURITooLong,
	"UnsupportedMediaType":          http.StatusUnsupportedMediaType,
	"RequestedRangeNotSatisfiable":  http.StatusRequestedRangeNotSatisfiable,
	"ExpectationFailed":             http.StatusExpectationFailed,
	"Teapot":                        http.StatusTeapot,
	"MisdirectedRequest":            http.StatusMisdirectedRequest,
	"UnprocessableEntity":           http.StatusUnprocessableEntity,
	"Locked":                        http.StatusLocked,
	"FailedDependency":              http.StatusFailedDependency,
	"TooEarly":                      http.StatusTooEarly,
	"UpgradeRequired":               http.StatusUpgradeRequired,
	"PreconditionRequired":          http.StatusPreconditionRequired,
	"TooManyRequests":               http.StatusTooManyRequests,
	"RequestHeaderFieldsTooLarge":   http.StatusRequestHeaderFieldsTooLarge,
	"UnavailableForLegalReasons":    http.StatusUnavailableForLegalReasons,
	"InternalServerError":           http.StatusInternalServerError,
	"NotImplemented":                http.StatusNotImplemented,
	"BadGateway":                    http.StatusBadGateway,
	"ServiceUnavailable":            http.StatusServiceUnavailable,
	"GatewayTimeout":                http.StatusGatewayTimeout,
	"HTTPVersionNotSupported":       http.StatusHTTPVersionNotSupported,
	"VariantAlsoNegotiates":         http.StatusVariantAlsoNegotiates,
	"InsufficientStorage":           http.StatusInsufficientStorage,
	"LoopDetected":                  http.StatusLoopDetected,
	"NotExtended":                   http.StatusNotExtended,
	"NetworkAuthenticationRequired": http.StatusNetworkAuthenticationRequired,
}

func TranslateStatus(status string) string {
	status = strings.ToLower(status)
	for k := range HttpStatus {
		if strings.ToLower(k) == status {
			return k
		}
	}

	console.Fatal("status %s not found in net/http.Status", status)
	return ""
}
