package model

import (
	"time"

	"gorm.io/gorm"
)

type FormResponse struct {
	CreateTime        time.Time      `json:"createTime" gorm:"column:create_time"`
	FormId            string         `json:"formId" gorm:"column:form_id"`
	LastSubmittedTime time.Time      `json:"lastSubmittedTime" gorm:"column:last_submitted_time"`
	RespondentEmail   string         `json:"respondentEmail" gorm:"column:respondent_email"`
	ResponseId        string         `json:"responseId" gorm:"column:response_id"`
	TotalScore        float64        `json:"totalScore" gorm:"column:total_score"`
	CreatedAt         time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (f *FormResponse) TableName() string {
	return "google_workspace_forms.form_response"
}
