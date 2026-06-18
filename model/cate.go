package model

import (
	"time"

	"github.com/shemic/dever/orm"
)

type Cate struct {
	ID        uint64    `dorm:"primaryKey;autoIncrement;comment:分类ID"`
	ChannelID uint64    `dorm:"type:bigint;not null;default:1;comment:所属频道"`
	Name      string    `dorm:"type:varchar(128);not null;comment:分类名称"`
	Status    int16     `dorm:"type:smallint;not null;default:1;comment:状态"`
	Sort      int       `dorm:"type:int;not null;default:100;comment:排序"`
	CreatedAt time.Time `dorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

type CateIndex struct {
	ChannelSort       struct{} `index:"channel_id,sort,id"`
	ChannelStatusSort struct{} `index:"channel_id,status,sort,id"`
}

var cateSeed = []map[string]any{
	{"id": DefaultCateID, "channel_id": DefaultChannelID, "name": "默认分类", "status": 1, "sort": 100},
}

var cateChannelRelation = orm.Relation{
	Field:      "channel_id",
	Option:     "source.NewChannelModel",
	OptionKeys: []string{"name"},
}

func NewCateModel() *orm.Model[Cate] {
	return orm.LoadModel[Cate]("分类", "cate", orm.ModelConfig{
		Index:    CateIndex{},
		Seeds:    cateSeed,
		Order:    "sort asc,id asc",
		Database: "default",
		Options: map[string]any{
			"status": enabledStatusOptions,
		},
		Relations: []orm.Relation{
			cateChannelRelation,
		},
	})
}
