// Code generated by "stringer -type=TokenType -output=token_type_string.go"; DO NOT EDIT.

package models

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TokenTypeInvalid-0]
	_ = x[TokenTypeCSRF-1]
	_ = x[TokenTypeConfirmContact-2]
	_ = x[TokenTypePasswordReset-3]
	_ = x[TokenTypeSession-4]
}

const _TokenType_name = "TokenTypeInvalidTokenTypeCSRFTokenTypeConfirmContactTokenTypePasswordResetTokenTypeSession"

var _TokenType_index = [...]uint8{0, 16, 29, 52, 74, 90}

func (i TokenType) String() string {
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
