package service

import (
	"strings"

	"github.com/vertex-center/uuid"
	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
)

type emailService struct {
	adapter port.EmailAdapter
}

func NewEmailService(adapter port.EmailAdapter) port.EmailService {
	return &emailService{
		adapter: adapter,
	}
}

func (s emailService) GetEmails(userID uuid.UUID) ([]types.Email, error) {
	return s.adapter.GetEmails(userID)
}

func (s emailService) CreateEmail(userID uuid.UUID, email string) (types.Email, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return types.Email{}, types.ErrEmailEmpty
	}

	res := types.Email{
		ID:     uuid.New(),
		UserID: userID,
		Email:  email,
	}
	err := s.adapter.CreateEmail(&res)
	return res, err
}

func (s emailService) DeleteEmail(userID uuid.UUID, email string) error {
	email = strings.TrimSpace(email)
	return s.adapter.DeleteEmail(userID, email)
}
