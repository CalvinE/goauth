package service

import (
	"fmt"

	coreerrors "github.com/calvine/goauth/core/errors"
	"github.com/calvine/richerror/errors"
)

const (
	MockEmailService = "mock"
	SMTPEmailService = "smtp"
)

type EmailRecipientInfo struct {
	Bcc []string `json:"bcc"`
	Cc  []string `json:"cc"`
	To  []string `json:"to"`
}

type EmailContentInfo struct {
	Body       string             `json:"body"`
	IsHTMLBody bool               `json:"isHtmlBody"`
	Subject    string             `json:"subject"`
	Recipients EmailRecipientInfo `json:",inline"`
}

type EmailService interface {
	SendPlainTextEmail(to []string, subject, body string) errors.RichError
}

type mockEmailService struct{}

func (mse mockEmailService) SendPlainTextEmail(to []string, subject, body string) errors.RichError {
	fmt.Println("********** BEGIN EMAIL  **********")

	fmt.Printf("TO:\t%v\n\n", to)

	fmt.Printf("SUBJECT:\t%s\n\n", subject)

	fmt.Printf("BODY:\t%s\n\n", body)

	fmt.Println("********** END EMAIL  **********")
	return nil
}

func NewEmailService(serviceType string, options interface{}) (EmailService, errors.RichError) {
	switch serviceType {
	case MockEmailService:
		return mockEmailService{}, nil
	// case SMTPEmailService: // TODO: implement this...
	default:
		return nil, coreerrors.NewComponentNotImplementedError("email service", serviceType, true)
	}
}
