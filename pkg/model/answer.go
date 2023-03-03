package model

import (
	"time"

	"gorm.io/gorm"
)

type Answer struct {
	QuestionId     string         `json:"questionId" gorm:"column:question_id"`
	FormId         string         `json:"formId" gorm:"form_id"`
	FormResponseId string         `json:"formResponseId" gorm:"form_response_id"`
	TextAnswers    string         `json:"textAnswers" gorm:"text_answers"`
	CreatedAt      time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (a *Answer) TableName() string {
	return "google_workspace_forms.answer"
}
