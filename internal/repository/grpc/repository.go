package grpc

import "github.com/mephistolie/chefbook-backend-encryption/internal/config"

type Repository struct {
	Auth    *Auth
	Profile *Profile
	Recipe  *Recipe
}

func NewRepository(cfg *config.Config) (*Repository, error) {
	authService, err := NewAuth(*cfg.AuthService.Addr)
	if err != nil {
		return nil, err
	}
	profileService, err := NewProfile(*cfg.ProfileService.Addr)
	if err != nil {
		return nil, err
	}
	recipeService, err := NewRecipe(*cfg.RecipeService.Addr)
	if err != nil {
		return nil, err
	}

	return &Repository{
		Auth:    authService,
		Profile: profileService,
		Recipe:  recipeService,
	}, nil
}

func (r *Repository) Stop() error {
	_ = r.Recipe.Conn.Close()
	return nil
}
