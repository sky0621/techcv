package domain

import "context"

type Profile struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
}

type ProfileRepository interface {
	Save(ctx context.Context, profile Profile) error
}
