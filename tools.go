//go:build tools
// +build tools

// Place any runtime dependencies as imports in this file.
// Go modules will be forced to download and install them.
package tools

import (
	_ "github.com/deepmap/oapi-codegen/pkg/codegen"
)
