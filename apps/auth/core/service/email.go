package service

import (
	"net/mail"
	"strings"

	"github.com/juju/errors"
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
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return types.Email{}, errors.NewBadRequest(err, "create email address")
	}

	res := types.Email{
		ID:     uuid.New(),
		UserID: userID,
		Email:  addr.Address,
	}
	err = s.adapter.CreateEmail(&res)
	if err != nil {
		return types.Email{}, err
	}
	return res, nil
}

func (s emailService) DeleteEmail(userID uuid.UUID, email string) error {
	email = strings.TrimSpace(email)
	return s.adapter.DeleteEmail(userID, email)
}
