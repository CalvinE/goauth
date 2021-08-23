package service

import (
	"context"
	"time"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/goauth/core/models"
	repo "github.com/calvine/goauth/core/repositories"
	"github.com/calvine/richerror/errors"
	"golang.org/x/crypto/bcrypt"
)

type loginService struct {
	auditLogRepo repo.AuditLogRepo
	contactRepo  repo.ContactRepo
	emailService EmailService
	userRepo     repo.UserRepo
}

func NewLoginService(auditLogRepo repo.AuditLogRepo, contactRepo repo.ContactRepo, emailService EmailService, userRepo repo.UserRepo) loginService {
	return loginService{
		auditLogRepo: auditLogRepo,
		contactRepo:  contactRepo,
		emailService: emailService,
		userRepo:     userRepo,
	}
}

func (ls loginService) LoginWithPrimaryContact(ctx context.Context, principal, principalType, password string, initiator string) (models.User, errors.RichError) {
	user, contact, err := ls.userRepo.GetUserAndContactByPrimaryContact(ctx, principalType, principal)
	if err != nil {
		return models.User{}, err
	}
	if !contact.IsPrimary {
		return models.User{}, coreerrors.NewLoginContactNotPrimaryError(contact.ID, contact.Principal, contact.Type, true)
	}
	now := time.Now().UTC()
	// is user locked out?
	if user.LockedOutUntil.HasValue && user.LockedOutUntil.Value.After(now) {
		return models.User{}, coreerrors.NewUserLockedOutError(user.ID, true)
	}
	if !contact.ConfirmedDate.HasValue { // || contact.ConfirmedDate.Value.After(now)
		return models.User{}, coreerrors.NewContactNotConfirmedError(contact.ID, contact.Principal, contact.Type, true)
	}
	// check password
	// hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(password), 10)
	bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if bcryptErr == bcrypt.ErrMismatchedHashAndPassword {
		user.ConsecutiveFailedLoginAttempts += 1
		// TODO: make max ConsecutiveFailedLoginAttempts configurable
		if user.ConsecutiveFailedLoginAttempts >= 10 {
			user.ConsecutiveFailedLoginAttempts = 0
			// TODO: make lockout time configurable
			user.LockedOutUntil.Set(now.Add(time.Minute * 15))
		}
		_ = ls.userRepo.UpdateUser(ctx, &user, user.ID)
		return models.User{}, coreerrors.NewLoginFailedWrongPasswordError(user.ID, true)
	} else if bcryptErr != nil {
		return models.User{}, coreerrors.NewBcryptPasswordHashErrorError(user.ID, bcryptErr, true)
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
