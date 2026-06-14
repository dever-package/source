package source

import "embed"

// PageFS 内嵌 source 后台页面配置。
//
//go:embed front/page/*/*/*.json
var PageFS embed.FS
