package model

import (
	"time"

	"github.com/shemic/dever/orm"
)

type Origin struct {
	ID        uint64    `dorm:"primaryKey;autoIncrement;comment:来源ID"`
	Name      string    `dorm:"type:varchar(128);not null;comment:来源名称"`
	Status    int16     `dorm:"type:smallint;not null;default:1;comment:状态"`
	Sort      int       `dorm:"type:int;not null;default:100;comment:排序"`
	CreatedAt time.Time `dorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

type OriginIndex struct {
	StatusSort struct{} `index:"status,sort,id"`
}

var originSeed = []map[string]any{
	{"id": DefaultOriginID, "name": "默认来源", "status": 1, "sort": 100},
}

func NewOriginModel() *orm.Model[Origin] {
	return orm.LoadModel[Origin]("来源", "source_origin", orm.ModelConfig{
		Index:    OriginIndex{},
		Seeds:    originSeed,
		Order:    "sort asc,id asc",
		Database: "default",
		Options: map[string]any{
			"status": enabledStatusOptions,
		},
	})
}
