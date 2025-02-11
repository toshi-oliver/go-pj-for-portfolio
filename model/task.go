package model

import "time"

type Task struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	MenuTitle  string    `json:"menuTitle" gorm:"not null"`
	MenuDetail string    `json:"menuDetail"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	User       User      `json:"user" gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`
	UserId     uint      `json:"userId" gorm:"not null"`
}

type TaskResponse struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	MenuTitle  string    `json:"menuTitle" gorm:"not null"`
	MenuDetail string    `json:"menuDetail" gorm:"not null"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type TaskResponsePaginated struct {
	Tasks       []TaskResponse `json:"tasks"`
	CurrentPage uint           `json:"currentPage"`
	LastPage    uint           `json:"lastPage"`
	TotalCount  int64          `json:"totalCount"`
}
