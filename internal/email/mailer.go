//go:generate mockgen -source mailer.go -destination mock/mailer.go -package mock
package email

import (
	"context"
	"github.com/email-microservice/internal/models"
)

// Mailer interface
type Mailer interface {
	Send(ctx context.Context, email *models.Email) error
}
