package dbtools

import (
	"context"

	"github.com/gosimple/slug"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/metal-toolbox/fleetdb/internal/models"
)

// RegisterHooks adds any hooks that are configured to the models library
func RegisterHooks() {
	models.AddServerComponentTypeHook(boil.BeforeInsertHook, setServerComponentTypeSlug)
	models.AddServerCredentialTypeHook(boil.BeforeInsertHook, setServerCredentialTypeSlug)
}

func setServerComponentTypeSlug(_ context.Context, _ boil.ContextExecutor, t *models.ServerComponentType) error {
	if t.Slug == "" {
		t.Slug = slug.Make(t.Name)
	}

	return nil
}

func setServerCredentialTypeSlug(_ context.Context, _ boil.ContextExecutor, t *models.ServerCredentialType) error {
	if t.Slug == "" {
		t.Slug = slug.Make(t.Name)
	}

	return nil
}
