package studiofrontendassets

import "embed"

//go:embed dist/* public/* index.html
var Assets embed.FS

const Dir = "studio/frontend"
