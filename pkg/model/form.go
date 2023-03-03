package model

import (
	"time"

	"google.golang.org/api/forms/v1"
	"gorm.io/gorm"
)

type Form struct {
	FormId        string         `json:"form_id" gorm:"column:form_id"`
	Description   string         `json:"description" gorm:"column:description"`
	DocumentTitle string         `json:"document_title" gorm:"column:document_title"`
	Title         string         `json:"title" gorm:"column:title"`
	LinkedSheetId string         `json:"linked_sheet_id" gorm:"column:linked_sheet_id"`
	ResponderUri  string         `json:"responder_uri" gorm:"column:responder_uri"`
	RevisionId    string         `json:"revision_id" gorm:"column:revision_id"`
	Items         []Item         `gorm:"-"`
	CreatedAt     *time.Time     `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     *time.Time     `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (f *Form) TableName() string {
	return "google_workspace_forms.form"
}

func FormToLocalSchema(form forms.Form) Form {
	items := []Item{}
	for _, i := range form.Items {
		fi := ItemToLocalSchema(*i, form.FormId)
		items = append(items, fi...)
	}
	return Form{
		FormId:        form.FormId,
		Description:   form.Info.Description,
		DocumentTitle: form.Info.DocumentTitle,
		Title:         form.Info.Title,
		LinkedSheetId: form.LinkedSheetId,
		ResponderUri:  form.ResponderUri,
		RevisionId:    form.RevisionId,
		Items:         items,
	}
}
