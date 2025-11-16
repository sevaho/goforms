package assets

import (
	"embed"
)

//go:embed templates
var Templates embed.FS

//go:embed migrations
var Migrations embed.FS

//go:embed static
var Static embed.FS
