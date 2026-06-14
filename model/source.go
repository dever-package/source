package model

import (
	"time"

	"github.com/shemic/dever/orm"
)

type Source struct {
	ID        uint64    `dorm:"primaryKey;autoIncrement;comment:资源ID"`
	Title     string    `dorm:"type:varchar(180);not null;comment:标题"`
	ChannelID uint64    `dorm:"type:bigint;not null;default:1;comment:频道"`
	CateID    uint64    `dorm:"type:bigint;not null;default:1;comment:分类"`
	OriginID  uint64    `dorm:"type:bigint;not null;default:1;comment:来源"`
	StatusID  uint64    `dorm:"type:bigint;not null;default:1;comment:资源状态"`
	Content   string    `dorm:"type:text;not null;default:'';comment:内容"`
	Images    string    `dorm:"type:text;not null;default:'';comment:图片列表"`
	Status    int16     `dorm:"type:smallint;not null;default:1;comment:是否展示"`
	Sort      int       `dorm:"type:int;not null;default:100;comment:排序"`
	CreatedAt time.Time `dorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}

type SourceIndex struct {
	ChannelCateStatusSort struct{} `index:"channel_id,cate_id,status,sort,id"`
	OriginStatusSort      struct{} `index:"origin_id,status,sort,id"`
	StatusIDSort          struct{} `index:"status_id,sort,id"`
	StatusSort            struct{} `index:"status,sort,id"`
	CreatedAt             struct{} `index:"created_at"`
}

var (
	sourceChannelRelation = orm.Relation{
		Field:      "channel_id",
		Name:       "channel",
		Option:     "source.NewChannelModel",
		OptionKeys: []string{"name"},
	}
	sourceCateRelation = orm.Relation{
		Field:      "cate_id",
		Name:       "cate",
		Option:     "source.NewCateModel",
		OptionKeys: []string{"name", "channel_id"},
	}
	sourceOriginRelation = orm.Relation{
		Field:      "origin_id",
		Name:       "origin",
		Option:     "source.NewOriginModel",
		OptionKeys: []string{"name"},
	}
	sourceStatusRelation = orm.Relation{
		Field:      "status_id",
		Name:       "source_status",
		Option:     "source.NewStatusModel",
		OptionKeys: []string{"name", "color", "online", "status"},
	}
)

func NewSourceModel() *orm.Model[Source] {
	return orm.LoadModel[Source]("资源", "source", orm.ModelConfig{
		Index:    SourceIndex{},
		Order:    "sort asc,id desc",
		Database: "default",
		Options: map[string]any{
			"status": enabledStatusOptions,
		},
		Relations: []orm.Relation{
			sourceChannelRelation,
			sourceCateRelation,
			sourceOriginRelation,
			sourceStatusRelation,
		},
	})
}
