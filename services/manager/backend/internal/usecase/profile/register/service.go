package register

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain"
)

var ErrInvalidInput = errors.New("invalid input")

type IDGenerator interface {
	NewID() string
}

type Service struct {
	profiles domain.ProfileRepository
	idGen    IDGenerator
}

func NewService(profiles domain.ProfileRepository, idGen IDGenerator) Service {
	return Service{
		profiles: profiles,
		idGen:    idGen,
	}
}

func (s Service) Register(ctx context.Context, name, nickname string) (domain.Profile, error) {
	trimmedName := strings.TrimSpace(name)
	trimmedNickname := strings.TrimSpace(nickname)

	if trimmedName == "" || trimmedNickname == "" {
		return domain.Profile{}, ErrInvalidInput
	}

	profile := domain.Profile{
		ID:       s.idGen.NewID(),
		Name:     trimmedName,
		Nickname: trimmedNickname,
	}

	if err := s.profiles.Save(ctx, profile); err != nil {
		return domain.Profile{}, fmt.Errorf("save profile: %w", err)
	}

	return profile, nil
}
