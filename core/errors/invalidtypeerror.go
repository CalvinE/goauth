package errors

import (
	"fmt"

	"github.com/calvine/goauth/core/errors/codes"
)

func NewInvalidTypeError(actual string, includeStack bool) RichError {
	msg := fmt.Sprintf("invalid type encountered: %s", actual)
	err := NewRichError(codes.ErrCodeInvalidType, msg, includeStack).AddMetaData("actual", actual)
	return err
}
