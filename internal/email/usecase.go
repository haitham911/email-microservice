//go:generate mockgen -source usecase.go -destination mock/usecase.go -package mock
package email

import (
	"context"

	"github.com/email-microservice/internal/models"
	"github.com/email-microservice/pkg/utils"
	"github.com/google/uuid"
)

// Image useCase interface
type EmailsUseCase interface {
	SendEmail(ctx context.Context, deliveryBody []byte) error
	FindEmailById(ctx context.Context, emailID uuid.UUID) (*models.Email, error)
	FindEmailsByReceiver(ctx context.Context, to string, query *utils.PaginationQuery) (*models.EmailsList, error)
}
