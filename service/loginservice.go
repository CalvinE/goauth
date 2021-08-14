package service

import (
	"context"
	"time"

	"github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/goauth/core/utilities"
)

type loginService struct {
	userRepo     repo.UserRepo
	contactRepo  repo.ContactRepo
	auditLogRepo repo.AuditLogRepo
}

func NewLoginService(userRepo repo.UserRepo, contactRepo repo.ContactRepo, auditLogRepo repo.AuditLogRepo) loginService {
	return loginService{
		userRepo:     userRepo,
		contactRepo:  contactRepo,
		auditLogRepo: auditLogRepo,
	}
}

func (ls loginService) LoginWithPrimaryContact(ctx context.Context, principal, principalType, password string, initiator string) (models.User, error) {
	user, err := ls.userRepo.GetUserByPrimaryContact(ctx, principalType, principal)
	if err != nil {
		return models.User{}, err
	}
	now := time.Now().UTC()
	// is user locked out?
	if user.LockedOutUntil.HasValue && user.LockedOutUntil.Value.After(now) {
		return models.User{}, errors.NewUserLockedOutError(user.ID, true)
	}
	// is contact confirmed?
	contact, err := ls.contactRepo.GetPrimaryContactByUserId(ctx, user.ID)
	if err != nil {
		return models.User{}, err
	}
	if !contact.IsPrimary {
		return models.User{}, errors.NewLoginContactNotPrimaryError(contact.ID, contact.Principal, contact.Type, true)
	}
	if !contact.ConfirmedDate.HasValue { // || contact.ConfirmedDate.Value.After(now)
		// TODO: return error that primary contact is not confirmed.
		return models.User{}, errors.NewContactNotConfirmedError(contact.ID, contact.Principal, contact.Type, true)
	}
	// check password
	saltedString := utilities.InterleaveStrings(password, user.Salt)
	// TODO: use bcrypt...
	computedHash, err := utilities.SHA512(saltedString)
	if err != nil {
		return models.User{}, errors.NewComputeHashFailedError("SHA512", err, true)
	}
	if computedHash != user.PasswordHash {
		user.ConsecutiveFailedLoginAttempts += 1
		// TODO: make max ConsecutiveFailedLoginAttempts configurable
		if user.ConsecutiveFailedLoginAttempts >= 10 {
			user.ConsecutiveFailedLoginAttempts = 0
			// TODO: make lockout time configurable
			user.LockedOutUntil.Set(now.Add(time.Minute * 15))
		}
		err = ls.userRepo.UpdateUser(ctx, &user, user.ID)
		return models.User{}, errors.NewLoginFailedWrongPasswordError(user.ID, true)
	}
	if user.ConsecutiveFailedLoginAttempts > 0 {
		// reset consecutive failed login attempts because we have a successful login
		user.ConsecutiveFailedLoginAttempts = 0
		err = ls.userRepo.UpdateUser(ctx, &user, user.ID)
		if err != nil {
			return models.User{}, err
		}
	}
	return user, nil
}

// func (ls loginService) StartPasswordResetByContact(ctx context.Context, principal, principalType string, initiator string) (string, error) {

// }

// func (ls loginService) ConfirmContact(ctx context.Context, confirmationCode string, initiator string) (bool, error) {

// }

// func (ls loginService) ResetPassword(ctx context.Context, userId string, newPasswordHash string, newSalt string, initiator string) (bool, error) {

// }
