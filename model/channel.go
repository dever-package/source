package model

import (
	"time"

	"github.com/shemic/dever/orm"
)

type Channel struct {
	ID        uint64    `dorm:"primaryKey;autoIncrement;comment:频道ID"`
	Name      string    `dorm:"type:varchar(128);not null;comment:频道名称"`
	Status    int16     `dorm:"type:smallint;not null;default:1;comment:状态"`
	Sort      int       `dorm:"type:int;not null;default:100;comment:排序"`
	CreatedAt time.Time `dorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

type ChannelIndex struct {
	StatusSort struct{} `index:"status,sort,id"`
}

var channelSeed = []map[string]any{
	{"id": DefaultChannelID, "name": "默认频道", "status": 1, "sort": 100},
}

func NewChannelModel() *orm.Model[Channel] {
	return orm.LoadModel[Channel]("频道", "source_channel", orm.ModelConfig{
		Index:    ChannelIndex{},
		Seeds:    channelSeed,
		Order:    "sort asc,id asc",
		Database: "default",
		Options: map[string]any{
			"status": enabledStatusOptions,
		},
	})
}
