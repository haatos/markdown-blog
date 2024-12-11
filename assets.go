package assets

import "embed"

//go:embed public/**
var StaticFS embed.FS

//go:embed migrations/*.sql
var MigrationsFS embed.FS
