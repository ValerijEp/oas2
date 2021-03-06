package oas

import (
	"fmt"
	"net/http"

	"github.com/hypnoglow/oas2/validate"
)

// NewQueryValidator returns new Middleware that validates request query
// parameters against OpenAPI 2.0 spec.
func NewQueryValidator(errHandler RequestErrorHandler) Middleware {
	return queryValidatorMiddleware{
		errHandler: errHandler,
	}
}

type queryValidatorMiddleware struct {
	errHandler RequestErrorHandler
}

func (m queryValidatorMiddleware) Apply(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		op := GetOperation(req)
		if op == nil {
			next.ServeHTTP(w, req)
			return
		}

		if errs := validate.Query(op.Parameters, req.URL.Query()); len(errs) > 0 {
			err := ValidationError{error: fmt.Errorf("validation error"), errs: errs}
			if !m.errHandler(w, req, err) {
				return
			}
		}

		next.ServeHTTP(w, req)
	})
}
