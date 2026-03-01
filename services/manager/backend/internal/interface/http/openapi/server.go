package openapi

import (
	"context"
	"errors"

	"github.com/sky0621/techcv/services/manager/backend/internal/domain/model"
	openapigen "github.com/sky0621/techcv/services/manager/backend/internal/interface/http/openapi/gen"
	profileregister "github.com/sky0621/techcv/services/manager/backend/internal/usecase/profile/register"
)

type profileRegistrar interface {
	Register(ctx context.Context, name, nickname string) (model.Profile, error)
}

type StrictServer struct {
	profileService profileRegistrar
}

func NewStrictServer(profileService profileRegistrar) *StrictServer {
	return &StrictServer{
		profileService: profileService,
	}
}

func (s *StrictServer) CreateProfile(ctx context.Context, request openapigen.CreateProfileRequestObject) (openapigen.CreateProfileResponseObject, error) {
	if request.Body == nil {
		return openapigen.CreateProfile400JSONResponse{
			Error: "invalid request body",
		}, nil
	}

	profile, err := s.profileService.Register(ctx, request.Body.Name, request.Body.Nickname)
	if err != nil {
		if errors.Is(err, profileregister.ErrInvalidInput) {
			return openapigen.CreateProfile400JSONResponse{
				Error: "name and nickname are required",
			}, nil
		}
		return openapigen.CreateProfile500JSONResponse{
			Error: "failed to register profile",
		}, nil
	}

	return openapigen.CreateProfile201JSONResponse{
		Id:       profile.ID,
		Name:     profile.Name,
		Nickname: profile.Nickname,
	}, nil
}

var _ openapigen.StrictServerInterface = (*StrictServer)(nil)
