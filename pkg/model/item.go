package model

import (
	"fmt"
	"time"

	"google.golang.org/api/forms/v1"
	"gorm.io/gorm"
)

type Item struct {
	Description      string         `json:"description" gorm:"column:description"`
	ItemId           string         `json:"item_id" gorm:"column:item_id"`
	Title            string         `json:"title" gorm:"column:title"`
	IsQuestionItem   bool           `json:"is_question_item" gorm:"column:is_question_item"`
	QuestionId       string         `json:"question_id" gorm:"column:question_id"`
	IsChoiceQuestion bool           `json:"is_choice_question" gorm:"column:is_choice_question"`
	IsDateQuestion   bool           `json:"is_date_question" gorm:"column:is_date_question"`
	IsRequired       bool           `json:"is_required" gorm:"column:is_required"`
	IsRowQuestion    bool           `json:"is_row_question" gorm:"column:is_row_question"`
	IsScaleQuestion  bool           `json:"is_scale_question" gorm:"column:is_scale_question"`
	IsTextQuestion   bool           `json:"is_text_question" gorm:"column:is_text_question"`
	IsTimeQuestion   bool           `json:"is_time_question" gorm:"column:is_time_question"`
	FormId           string         `json:"form_id" gorm:"column:form_id"`
	CreatedAt        *time.Time     `json:"created_at" gorm:"column:created_at"`
	UpdatedAt        *time.Time     `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (f *Item) TableName() string {
	return "google_workspace_forms.item"
}

func ItemToLocalSchema(fi forms.Item, form_id string) []Item {
	result := []Item{}
	out := Item{
		Description: fi.Description,
		ItemId:      fi.ItemId,
		Title:       fi.Title,
		FormId:      form_id,
	}
	isQuestionItem := fi.QuestionItem != nil && fi.QuestionItem.Question != nil
	if isQuestionItem {
		out.QuestionId = fi.QuestionItem.Question.QuestionId
		out.IsChoiceQuestion = fi.QuestionItem.Question.ChoiceQuestion != nil
		out.IsDateQuestion = fi.QuestionItem.Question.DateQuestion != nil
		out.IsRequired = fi.QuestionItem.Question.Required
		out.IsScaleQuestion = fi.QuestionItem.Question.ScaleQuestion != nil
		out.IsTextQuestion = fi.QuestionItem.Question.TextQuestion != nil
		out.IsTimeQuestion = fi.QuestionItem.Question.TimeQuestion != nil
		out.IsQuestionItem = isQuestionItem
		result = append(result, out)
	}
	isQuestionGroupItem := fi.QuestionGroupItem != nil && len(fi.QuestionGroupItem.Questions) > 0
	if isQuestionGroupItem {
		for _, question := range fi.QuestionGroupItem.Questions {
			clone := out
			clone.QuestionId = question.QuestionId
			if question.RowQuestion != nil {
				clone.Title = fmt.Sprint(out.Title, "-", question.RowQuestion.Title)
			}
			result = append(result, clone)
		}
	}
	return result
}
