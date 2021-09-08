package services

import (
	"context"

	"github.com/calvine/goauth/core/models"
	"github.com/calvine/richerror/errors"
)

// TODO: change all instances of error to RichError

// LoginService is a service used to facilitate logging in
type LoginService interface {
	// LoginWithContact attempts to confirm a users credentials and if they match it returns true and resets the users ConsecutiveFailedLoginAttempts, otherwise it returns false and increments the users ConsecutiveFailedLoginAttempts
	// The principal should only work when it has been confirmed
	LoginWithPrimaryContact(ctx context.Context, principal, principalType, password string, initiator string) (models.User, errors.RichError)
	// StartPasswordResetByContact sets a password reset token for the user with the corresponding principal and type that are confirmed.
	StartPasswordResetByContact(ctx context.Context, principal, principalType string, initiator string) (string, errors.RichError)
	// ResetPassword resets a users password given a userId and new password hash and salt.
	ResetPassword(ctx context.Context, passwordResetToken string, newPasswordHash string, initiator string) (bool, errors.RichError)
}

// UserService is a service that facilitates access to user related data.
type UserService interface {
	// GetUserByConfirmedContact gets a user record via a confirmed contact
	GetUserByConfirmedContact(ctx context.Context, contactPrincipal string, initiator string) (models.User, errors.RichError)
	// AddUser adds a user record to the database
	AddUser(ctx context.Context, user *models.User, initiator string) errors.RichError
	// UpdateUser updates a use record in the database
	UpdateUser(ctx context.Context, user *models.User, initiator string) errors.RichError
	// GetUserPrimaryContact gets a users primary contact
	GetUserPrimaryContact(ctx context.Context, userId string, initiator string) (models.Contact, errors.RichError)
	// GetUsersContacts gets all of a users contacts
	GetUsersContacts(ctx context.Context, userId string, initiator string) ([]models.Contact, errors.RichError)
	// GetUsersConfirmedContacts gets all of a users confirmed contacts
	GetUsersConfirmedContacts(ctx context.Context, userId string, initiator string) ([]models.Contact, errors.RichError)
	// AddContact adds a contact to a user
	AddContact(ctx context.Context, contact *models.Contact, initiator string) errors.RichError
	// UpdateContact updates a contact for a user
	UpdateContact(ctx context.Context, contact *models.Contact, initiator string) errors.RichError
	// ConfirmContact takes a confirmation code and updates the users contact record to be confirmed.
	ConfirmContact(ctx context.Context, confirmationCode string, initiator string) (bool, errors.RichError)
}

type EmailService interface {
	SendPlainTextEmail(to []string, subject, body string) errors.RichError
}

type TokenService interface {
	GetToken(tokenValue string, expectedTokenType models.TokenType) (models.Token, errors.RichError)
	PutToken(token models.Token) errors.RichError
	DeleteToken(tokenValue string) errors.RichError
}
