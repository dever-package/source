package source

import "embed"

// ManifestFS 内嵌 source 组件声明。
//
//go:embed dever.json
var ManifestFS embed.FS

// PageFS 内嵌 source 后台页面配置。
//
//go:embed front/page/*/*/*.json
var PageFS embed.FS
