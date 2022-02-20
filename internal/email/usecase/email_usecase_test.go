package usecase

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/email-microservice/config"
	"github.com/email-microservice/internal/email/mock"
	"github.com/email-microservice/internal/models"
	"github.com/email-microservice/pkg/logger"
	"github.com/email-microservice/pkg/mime_types"
	"github.com/email-microservice/pkg/utils"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestEmailUseCase_SendEmail(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Logger: config.Logger{},
		Smtp: config.Smtp{
			User: "mailservice@mail.ru",
		},
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	emailsRepository := mock.NewMockEmailsRepository(ctrl)
	mailer := mock.NewMockMailer(ctrl)

	emailUC := NewEmailUseCase(emailsRepository, appLogger, mailer, cfg)

	deliveryBody := []byte(`{"to": ["alex@gmail.com"], "from": "mailservice@mail.ru",
  "subject": "registration confirmation", "body": "registration confirmation body",
  "contentType": "text/plain"}`)

	ctx := context.Background()

	mail := &models.Email{}
	err := json.Unmarshal(deliveryBody, &mail)
	require.NoError(t, err)

	mailUUID := uuid.New()
	//mail.EmailID = mailUUID
	mail.From = "alex@gmail.com"

	mailer.EXPECT().Send(gomock.Any(), gomock.Any()).Return(nil)
	emailsRepository.EXPECT().CreateEmail(gomock.Any(), gomock.Any()).Return(&models.Email{
		EmailID:     mailUUID,
		To:          mail.To,
		From:        mail.From,
		Body:        mail.Body,
		Subject:     mail.Subject,
		ContentType: mail.ContentType,
		CreatedAt:   mail.CreatedAt,
	}, nil)

	err = emailUC.SendEmail(ctx, deliveryBody)
	require.NoError(t, err)
}

func TestEmailUseCase_FindEmailsByReceiver(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Logger: config.Logger{},
		Smtp: config.Smtp{
			User: "mailservice@mail.ru",
		},
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()

	emailsRepository := mock.NewMockEmailsRepository(ctrl)
	mailer := mock.NewMockMailer(ctrl)

	emailUC := NewEmailUseCase(emailsRepository, appLogger, mailer, cfg)

	emailUUID := uuid.New()
	mockEmail := &models.Email{
		EmailID:     emailUUID,
		To:          []string{"mail@gmail.com"},
		From:        "alex@gmail.com",
		Body:        "<span>some text content</span>",
		Subject:     "Confirm your email",
		ContentType: mime_types.MIMEApplicationJSON,
	}

	pq := &utils.PaginationQuery{
		Size: 10,
		Page: 1,
	}

	mockEmailsList := &models.EmailsList{
		TotalCount: 1,
		TotalPages: 1,
		Page:       1,
		Size:       1,
		HasMore:    false,
		Emails:     []*models.Email{mockEmail},
	}

	ctx := context.Background()
	emailsRepository.EXPECT().FindEmailsByReceiver(gomock.Any(), mockEmail.To[0], pq).Return(mockEmailsList, nil)

	result, err := emailUC.FindEmailsByReceiver(ctx, mockEmail.To[0], pq)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, uint64(len(result.Emails)), mockEmailsList.TotalCount)
}
