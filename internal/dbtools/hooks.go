package dbtools

import (
	"context"

	"github.com/gosimple/slug"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"go.hollow.sh/serverservice/internal/models"
)

// RegisterHooks adds any hooks that are configured to the models library
func RegisterHooks() {
	models.AddServerComponentTypeHook(boil.BeforeInsertHook, setServerComponentTypeSlug)
	models.AddServerSecretTypeHook(boil.BeforeInsertHook, setServerSecretTypeSlug)
}

func setServerComponentTypeSlug(ctx context.Context, exec boil.ContextExecutor, t *models.ServerComponentType) error {
	if t.Slug == "" {
		t.Slug = slug.Make(t.Name)
	}

	return nil
}

func setServerSecretTypeSlug(ctx context.Context, exec boil.ContextExecutor, t *models.ServerSecretType) error {
	if t.Slug == "" {
		t.Slug = slug.Make(t.Name)
	}

	return nil
}
