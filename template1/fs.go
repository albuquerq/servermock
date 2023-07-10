package template1

import "embed"

//go:embed root/*.gotmpl
//go:embed parts/*.gotmpl
var FS embed.FS
