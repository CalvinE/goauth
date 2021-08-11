package errors

/* WARNING: This is GENERATED CODE Please do not edit. */

import (
	"github.com/calvine/goauth/core/errors/codes"
)

// NewContactNotConfirmedError creates a new specific error
func NewContactNotConfirmedError(contactId string, principal string, principalType string, includeStack bool) RichError {
	msg := "contact is not confirmed"
	err := NewRichError(codes.ErrCodeContactNotConfirmed, msg).AddMetaData("contactId", contactId).AddMetaData("principal", principal).AddMetaData("principalType", principalType)
	if includeStack {
		err = err.WithStack(1)
	}
	return err
}
