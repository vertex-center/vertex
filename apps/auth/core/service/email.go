package service

import (
	"strings"

	"github.com/vertex-center/vertex/apps/auth/core/port"
	"github.com/vertex-center/vertex/apps/auth/core/types"
)

type EmailService struct {
	adapter port.EmailAdapter
}

func NewEmailService(adapter port.EmailAdapter) port.EmailService {
	return &EmailService{
		adapter: adapter,
	}
}

func (s EmailService) GetEmails(userID uint) ([]types.Email, error) {
	return s.adapter.GetEmails(userID)
}

func (s EmailService) CreateEmail(userID uint, email string) (types.Email, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return types.Email{}, types.ErrEmailEmpty
	}

	res := types.Email{
		UserID: userID,
		Email:  email,
	}
	err := s.adapter.CreateEmail(&res)
	return res, err
}

func (s EmailService) DeleteEmail(userID uint, email string) error {
	email = strings.TrimSpace(email)
	return s.adapter.DeleteEmail(userID, email)
}
