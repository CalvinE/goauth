package jwt

import (
	"testing"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/internal/testutils"
)

func TestNewJWTValidator(t *testing.T) {
	type testCase struct {
		name                string
		jwtValidatorOptions JWTValidatorOptions
		expectedErrorCode   string
	}
	testCases := []testCase{
		{
			name: "GIVEN a valid set of jwt validator options EXPECT no errors to occur",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret: "test",
			},
		},
		{
			name: "GIVEN jwt validator options with HS algorithms allowed and no hmacSecret set EXPECT error code jwt validator no hmac secret provided",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
			},
			expectedErrorCode: coreerrors.ErrCodeJWTValidatorNoHMACSecretProvided,
		},
		{
			name:                "GIVEN jwt validator options with no allowed algorithms EXPECT error code jwt validator no algorithm specified",
			jwtValidatorOptions: JWTValidatorOptions{},
			expectedErrorCode:   coreerrors.ErrCodeJWTValidatorNoAlgorithmSpecified,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tc.jwtValidatorOptions)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			} else {
				if validator.GetID() == "" {
					t.Error("\tvalidator id should never be an empty string")
				}
			}
		})
	}
}

func TestValidateHeader(t *testing.T) {
	type testCase struct {
		name                string
		jwtValidatorOptions JWTValidatorOptions
		header              Header
		expectValid         bool
		expectedErrorCodes  []string
	}
	testCases := []testCase{
		{
			name: "GIVEN a header with an allowed algorithm EXPECT no errors to occurr",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
					Alg_HS384,
					Alg_HS512,
				},
				HMACSecret: "test secret",
			},
			header: Header{
				Algorithm: Alg_HS256,
			},
			expectValid: true,
		},
		{
			name: "GIVEN a header with an algorithm that is not in the allowed algorithms list EXPECT error code jwt algorithm not allowed",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS384,
					Alg_HS512,
				},
				HMACSecret: "test secret",
			},
			header: Header{
				Algorithm: Alg_HS256,
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTAlgorithmNotAllowed,
			},
			expectValid: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tc.jwtValidatorOptions)
			if err != nil {
				t.Errorf("\texpected an error to occurr while building JWT validator: %s", err.GetErrorCode())
			}
			errs, valid := validator.ValidateHeader(tc.header)
			numErrors := len(errs)
			numExpectedErrors := len(tc.expectedErrorCodes)
			if numErrors != numExpectedErrors {
				t.Errorf("\texpected number of errors incorrect: got - %d expected - %d", numErrors, numExpectedErrors)
			}
			for _, e := range errs {
				errorCode := e.GetErrorCode()
				found := false
				for _, eec := range tc.expectedErrorCodes {
					if errorCode == eec {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("\terror occurred that was not expected: %s", errorCode)
				}
			}
			if valid != tc.expectValid {
				t.Errorf("\texpected header valid state incorrect: got - %v expected - %v", valid, tc.expectValid)
			}
		})
	}
}

func TestValidateClaims(t *testing.T) {
	type testCase struct {
		name                string
		jwtValidatorOptions JWTValidatorOptions
		body                StandardClaims
		expectValid         bool
		expectedErrorCodes  []string
	}
	testCases := []testCase{
		// issuer tests
		{
			name: "GIVEN claims with an expected issuer when no issuer is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:     "test",
				ExpectedIssuer: "goauth",
				IssuerRequired: false,
			},
			body: StandardClaims{
				Issuer: "goauth",
			},
			expectValid: true,
		},
		{
			name: "GIVEN claims with an expected issuer when issuer is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:     "test",
				ExpectedIssuer: "goauth",
				IssuerRequired: true,
			},
			body: StandardClaims{
				Issuer: "goauth",
			},
			expectValid: true,
		},
		{
			name: "GIVEN claims with no issuer when issuer is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:     "test",
				IssuerRequired: false,
			},
			expectValid: true,
			body:        StandardClaims{},
		},
		{
			name: "GIVEN claims with an unexpected issuer when issuer is required EXPECT failure",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:     "test",
				ExpectedIssuer: "goauth",
				IssuerRequired: false,
			},
			body: StandardClaims{
				Issuer: "other",
			},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuerInvalid,
			},
		},
		{
			name: "GIVEN claims with an unexpected issuer when issuer is not requiredEXPECT failure",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:     "test",
				ExpectedIssuer: "goauth",
				IssuerRequired: true,
			},
			body: StandardClaims{
				Issuer: "other",
			},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuerInvalid,
			},
		},
		{
			name: "GIVEN claims with no issuer when issuer is required EXPECT failure",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:     "test",
				IssuerRequired: true,
			},
			body:        StandardClaims{},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuerMissing,
			},
		},
		// expire tests
		{
			name: "GIVEN a jwt with a valid expire EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret: "test",
			},
			body: StandardClaims{
				ExpirationTime: Time(time.Now().Add(time.Second * 10)),
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with a valid expire and exipre is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				ExpireRequired: true,
				HMACSecret:     "test",
			},
			body: StandardClaims{
				ExpirationTime: Time(time.Now().Add(time.Second * 10)),
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with no exipre and expire is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				ExpireRequired: false,
				HMACSecret:     "test",
			},
			body:        StandardClaims{},
			expectValid: true,
		},
		{
			name:        "GIVEN a jwt that has expired EXPECT error code jwt expired",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret: "test",
			},
			body: StandardClaims{
				ExpirationTime: Time(time.Now().Add(time.Second * -1)),
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTExpired,
			},
		},
		{
			name:        "GIVEN a jwt that has expired and exp is required EXPECT error code jwt expired",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				ExpireRequired: true,
				HMACSecret:     "test",
			},
			body: StandardClaims{
				ExpirationTime: Time(time.Now().Add(time.Second * -1)),
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTExpired,
			},
		},
		{
			name:        "GIVEN a jwt that has no expireation but expire is required EXPECT error code jwt expire missing",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				ExpireRequired: true,
				HMACSecret:     "test",
			},
			body: StandardClaims{},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTExpireMissing,
			},
		},
		// audience tests
		// {
		// 	name: "GIVEN EXPECT ",
		// },
		// issued at tests
		{
			name: "GIVEN a jwt with a valid issued at EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret: "test",
			},
			body: StandardClaims{
				IssuedAt: Time(time.Now().Add(time.Second * -10)),
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with a valid issued at and issued at is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				IssuedAtRequired: true,
				HMACSecret:       "test",
			},
			body: StandardClaims{
				IssuedAt: Time(time.Now().Add(time.Second * -10)),
			},
			expectValid: true,
		},
		{
			name: "GIVEN a jwt with no issued at and issued at is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				IssuedAtRequired: false,
				HMACSecret:       "test",
			},
			body:        StandardClaims{},
			expectValid: true,
		},
		{
			name:        "GIVEN a jwt that has an invalid issued at EXPECT error code jwt issued at invalid",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret: "test",
			},
			body: StandardClaims{
				IssuedAt: Time(time.Now().Add(time.Second * 10)), // in the future
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuedAtInvalid,
			},
		},
		{
			name:        "GIVEN a jwt that has an invalid issued at and issued at is required EXPECT error code jwt issued at invalid",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				IssuedAtRequired: true,
				HMACSecret:       "test",
			},
			body: StandardClaims{
				IssuedAt: Time(time.Now().Add(time.Second * 10)),
			},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuedAtInvalid,
			},
		},
		{
			name:        "GIVEN a jwt that has no issued at but expire is required EXPECT error code jwt issued at missing",
			expectValid: false,
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				IssuedAtRequired: true,
				HMACSecret:       "test",
			},
			body: StandardClaims{},
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIssuedAtMissing,
			},
		},
		// not before tests
		// {
		// 	name: "GIVEN EXPECT ",
		// },
		// subject tests
		{
			name: "GIVEN jwt with subject when the subject is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:      "test",
				SubjectRequired: true,
			},
			body: StandardClaims{
				Subject: "user id",
			},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with subject when the subject is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:      "test",
				SubjectRequired: false,
			},
			body: StandardClaims{
				Subject: "user id",
			},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with no subject when the subject is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:      "test",
				SubjectRequired: false,
			},
			body:        StandardClaims{},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with no id when the jwt id is required EXPECT error code jwt id missing",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:      "test",
				SubjectRequired: true,
			},
			body:        StandardClaims{},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTSubjectMissing,
			},
		},
		// jwt id tests
		{
			name: "GIVEN jwt with an id when the jwt id is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:  "test",
				JTIRequired: true,
			},
			body: StandardClaims{
				JWTID: "unique id",
			},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with an id when the jwt id is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:  "test",
				JTIRequired: false,
			},
			body: StandardClaims{
				JWTID: "unique id",
			},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with an id when the jwt id is required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:  "test",
				JTIRequired: true,
			},
			body: StandardClaims{
				JWTID: "unique id",
			},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with no id when the jwt id is not required EXPECT success",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:  "test",
				JTIRequired: false,
			},
			body:        StandardClaims{},
			expectValid: true,
		},
		{
			name: "GIVEN jwt with no id when the jwt id is required EXPECT error code jwt id missing",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret:  "test",
				JTIRequired: true,
			},
			body:        StandardClaims{},
			expectValid: false,
			expectedErrorCodes: []string{
				coreerrors.ErrCodeJWTIDMissing,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tc.jwtValidatorOptions)
			if err != nil {
				t.Errorf("\texpected an error to occurr while building JWT validator: %s", err.GetErrorCode())
			}
			errs, valid := validator.ValidateClaims(tc.body)
			numErrors := len(errs)
			numExpectedErrors := len(tc.expectedErrorCodes)
			if numErrors != numExpectedErrors {
				t.Errorf("\texpected number of errors incorrect: got - %d expected - %d", numErrors, numExpectedErrors)
			}
			for _, e := range errs {
				errorCode := e.GetErrorCode()
				found := false
				for _, eec := range tc.expectedErrorCodes {
					if errorCode == eec {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("\terror occurred that was not expected: %s", errorCode)
				}
			}
			if valid != tc.expectValid {
				t.Errorf("\texpected claims valid state incorrect: got - %v expected - %v", valid, tc.expectValid)
			}
		})
	}
	t.Error("\ttest not implemented yet...")
}

func TestValidateSignature(t *testing.T) {
	type testCase struct {
		name                 string
		jwtValidatorOptions  JWTValidatorOptions
		alg                  string
		encodedHeaderAndBody string
		signature            string
		expectValid          bool
		expectedErrorCode    string
	}
	testCases := []testCase{
		{
			name: "GIVEN a valid HS256 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret: "test",
			},
			alg:                  Alg_HS256,
			encodedHeaderAndBody: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "5mhBHqs5_DTLdINd9p5m7ZJ6XD0Xc55kIaCRY5r6HRA",
			expectValid:          true,
		},
		{
			name: "GIVEN a valid HS384 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS384,
				},
				HMACSecret: "test",
			},
			alg:                  Alg_HS384,
			encodedHeaderAndBody: "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "KOZqnJ-wEzC-JvqqIHGKBIGgbYHH2Fej71TpBctnIguBkf3EdSYiwuRMSz35uY8E",
			expectValid:          true,
		},
		{
			name: "GIVEN a valid HS512 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS512,
				},
				HMACSecret: "test",
			},
			alg:                  Alg_HS512,
			encodedHeaderAndBody: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "VXfjNdZn9mDxRYhiaCi8rYYtcuNe3KCfK3LvggWSaHwjZsag9ugMOuDPOeeBD3oNhK-cOkTvRLy_ERbgnEyxYA",
			expectValid:          true,
		},
		{
			name: "GIVEN an invalid HS256 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS256,
				},
				HMACSecret: "test",
			},
			alg:                  Alg_HS256,
			encodedHeaderAndBody: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "5mhhHqs5_DTLdINd9p5m7ZJ6XD0Xc55kIaCRY5r6HRA",
			expectValid:          false,
		},
		{
			name: "GIVEN an invalid HS384 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS384,
				},
				HMACSecret: "test",
			},
			alg:                  Alg_HS384,
			encodedHeaderAndBody: "eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "KOZqqJ-wEzC-JvqqIHGKBIGgbYHH2Fej71TpBctnIguBkf3EdSYiwuRMSz35uY8E",
			expectValid:          false,
		},
		{
			name: "GIVEN an invalid HS512 jwt signature for the given encodedHeadAndBody EXPECT the signature to be valid",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_HS512,
				},
				HMACSecret: "test",
			},
			alg:                  Alg_HS512,
			encodedHeaderAndBody: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "VXfjNddn9mDxRYhiaCi8rYYtcuNe3KCfK3LvggWSaHwjZsag9ugMOuDPOeeBD3oNhK-cOkTvRLy_ERbgnEyxYA",
			expectValid:          false,
		},
		{
			name: "GIVEN a non supported jwt signature algorithm for the given encodedHeadAndBody EXPECT error code jwt algorithm not implemented",
			jwtValidatorOptions: JWTValidatorOptions{
				AllowedAlgorithms: []string{
					Alg_NONE,
				},
			},
			alg:                  Alg_NONE,
			encodedHeaderAndBody: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ",
			signature:            "VXfjNddn9mDxRYhiaCi8rYYtcuNe3KCfK3LvggWSaHwjZsag9ugMOuDPOeeBD3oNhK-cOkTvRLy_ERbgnEyxYA",
			expectValid:          false,
			expectedErrorCode:    coreerrors.ErrCodeJWTAlgorithmNotImplemented,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			validator, err := NewJWTValidator(tc.jwtValidatorOptions)
			if err != nil {
				t.Errorf("\texpected an error to occurr while building JWT validator: %s", err.GetErrorCode())
			}
			valid, err := validator.ValidateSignature(tc.alg, tc.encodedHeaderAndBody, tc.signature)
			if err != nil {
				testutils.HandleTestError(t, err, tc.expectedErrorCode)
			} else if tc.expectedErrorCode != "" {
				t.Errorf("\texpected an error to occurr: %s", tc.expectedErrorCode)
			}
			if valid != tc.expectValid {
				t.Errorf("\texpected signature valid state incorrect: got - %v expected - %v", valid, tc.expectValid)
			}
		})
	}
}
