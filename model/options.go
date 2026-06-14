package model

var enabledStatusOptions = []map[string]any{
	{"id": 1, "value": "启用", "label": "启用", "color": "#0f766e"},
	{"id": 2, "value": "停用", "label": "停用", "color": "#737373"},
}

var yesNoOptions = []map[string]any{
	{"id": 1, "value": "是", "label": "是", "color": "#0f766e"},
	{"id": 0, "value": "否", "label": "否", "color": "#737373"},
}

const (
	DefaultChannelID uint64 = 1
	DefaultCateID    uint64 = 1
	DefaultOriginID  uint64 = 1
	DefaultStatusID  uint64 = 1
)
