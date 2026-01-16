package templates

import "embed"

// FS contains the contents of the templates directory.
// The //go:embed directive is crucial for embedding these files into the compiled binary.
// Since this file is located in the 'templates' directory,
// the paths 'common' and 'rest' correctly refer to the template folders.
//

//go:embed common/**
//go:embed rest/**
//go:embed db/**
var FS embed.FS
