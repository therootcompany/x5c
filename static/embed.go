package static

import "embed"

// curl -L -O https://andybrewer.github.io/mvp/mvp.css

//go:embed index.html mvp.css fileviewer.js
var FS embed.FS
var Prefix = ""
