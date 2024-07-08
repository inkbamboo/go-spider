package model

import "time"

type BaseModel struct {
	ID        int64      `json:"id" structs:"id" gorm:"primaryKey;column:id;not null;unsigned"`
	CreatedAt *time.Time `json:"created_at,omitempty"  structs:"created_at" gorm:"autoCreateTime;column:created_at;type:TIMESTAMP;<-:create;comment:创建时间"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"  structs:"updated_at"  gorm:"autoUpdateTime;column:updated_at;type:TIMESTAMP;comment:更新时间"`
}
