package database

import "embed"

//go:embed migrations
var migrationFiles embed.FS
