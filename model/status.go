package model

import (
	"time"

	"github.com/shemic/dever/orm"
)

type Status struct {
	ID        uint64    `dorm:"primaryKey;autoIncrement;comment:状态ID"`
	Name      string    `dorm:"type:varchar(64);not null;comment:状态名称"`
	Color     string    `dorm:"type:varchar(32);not null;default:'';comment:展示颜色"`
	Online    int16     `dorm:"type:smallint;not null;default:1;comment:可展示"`
	Status    int16     `dorm:"type:smallint;not null;default:1;comment:状态"`
	Sort      int       `dorm:"type:int;not null;default:100;comment:排序"`
	CreatedAt time.Time `dorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

type StatusIndex struct {
	Name       struct{} `unique:"name"`
	StatusSort struct{} `index:"status,sort,id"`
	Online     struct{} `index:"online,status,sort,id"`
}

var statusSeed = []map[string]any{
	{"id": 1, "name": "上架", "color": "#0f766e", "online": 1, "status": 1, "sort": 100},
	{"id": 2, "name": "下架", "color": "#737373", "online": 0, "status": 1, "sort": 200},
}

func NewStatusModel() *orm.Model[Status] {
	return orm.LoadModel[Status]("资源状态", "source_status", orm.ModelConfig{
		Index:    StatusIndex{},
		Seeds:    statusSeed,
		Order:    "sort asc,id asc",
		Database: "default",
		Options: map[string]any{
			"online": yesNoOptions,
			"status": enabledStatusOptions,
		},
	})
}
