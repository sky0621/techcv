package user

import "context"

// UserRepository defines persistence operations for user aggregates.
type UserRepository interface {
	Create(ctx context.Context, user User) error
	FindByID(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email Email) (User, error)
	FindByFirebaseUID(ctx context.Context, firebaseUID FirebaseUID) (User, error)
	UpdateFromFirebase(ctx context.Context, user User) error
}
