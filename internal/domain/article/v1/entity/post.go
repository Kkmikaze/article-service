package entity

import (
	articlev1 "article-service/stubs/article/v1"
	_ "gorm.io/gorm"
	"time"
)

type Post struct {
	CreatedDate time.Time            `gorm:"type:timestamptz;default:now()" json:"created_date"`
	UpdatedDate time.Time            `gorm:"type:timestamptz;default:now()" json:"updated_date"`
	ID          string               `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title       string               `gorm:"type:varchar(200);not null" json:"title"`
	Content     string               `gorm:"type:text;not null" json:"content"`
	Category    string               `gorm:"type:varchar(100);not null" json:"category"`
	Status      articlev1.PostStatus `gorm:"type:int;default:0;not null" json:"status"`
}

func (e *Post) TableName() string {
	return "posts"
}
