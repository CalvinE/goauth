package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/richerror/errors"
)

// ErrCodeNoJWTSigningMaterialFound no jwt signing material found for given query
const ErrCodeNoJWTSigningMaterialFound = "NoJWTSigningMaterialFound"

// NewNoJWTSigningMaterialFoundError creates a new specific error
func NewNoJWTSigningMaterialFoundError(keyID string, includeStack bool) errors.RichError {
	msg := "no jwt signing material found for given query"
	err := errors.NewRichError(ErrCodeNoJWTSigningMaterialFound, msg).AddMetaData("keyID", keyID)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}

func IsNoJWTSigningMaterialFoundError(err errors.ReadOnlyRichError) bool {
	return err.GetErrorCode() == ErrCodeNoJWTSigningMaterialFound
}
