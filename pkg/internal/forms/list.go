package forms

import (
	"encoding/json"
	"time"

	alog "github.com/apex/log"
	"gitlab.com/the-new-orleans-connection/qr-code/model"
	"google.golang.org/api/forms/v1"
	"gorm.io/gorm"
)

func ListFormResponsesSynchronizeDB(db *gorm.DB, forms_service *forms.Service, form_id string, next_page_token string) (err error) {
	defer alog.WithField("form_id", form_id).WithField("next_page_token", next_page_token).Trace("ListFormResponsesSynchronizeDB").Stop(&err)

	list_form_responses_response, err := listFormResponses(forms_service, form_id, next_page_token)
	if err != nil {
		return
	}
	alog.
		WithField("num_rows", len(list_form_responses_response.Responses)).
		WithField("http_status_code", list_form_responses_response.HTTPStatusCode).
		WithField("next_page_token", list_form_responses_response.NextPageToken).
		Info("responses from google_forms")

	for _, form_response := range list_form_responses_response.Responses {
		err = createInDBIfNotExists(db, *form_response, form_id)
		if err != nil {
			return
		}
	}

	if list_form_responses_response.NextPageToken != "" {
		return ListFormResponsesSynchronizeDB(db, forms_service, form_id, list_form_responses_response.NextPageToken)
	} else {
		return nil
	}
}

func listFormResponses(forms_service *forms.Service, form_id string, next_page_token string) (list_form_responses_response *forms.ListFormResponsesResponse, err error) {
	defer alog.Trace("listFormResponses").Stop(&err)

	form_response_list_call := forms_service.Forms.Responses.List(form_id)
	if next_page_token != "" {
		form_response_list_call.PageToken(next_page_token)
	}
	return form_response_list_call.Do()
}

func createInDBIfNotExists(db *gorm.DB, response forms.FormResponse, form_id string) (err error) {
	defer alog.WithField("response_id", response.ResponseId).WithField("form_id", form_id).Trace("createInDbIfNotExists").Stop(&err)

	response_to_bytes, err := json.Marshal(response)
	if err != nil {
		return
	}

	alog.WithField("response", string(response_to_bytes)).Info("response as json")

	answers := map[string]model.Answer{}
	fr := model.FormResponse{
		FormId:          form_id,
		RespondentEmail: response.RespondentEmail,
		ResponseId:      response.ResponseId,
	}
	fr_create_time, err := time.Parse(time.RFC3339, response.CreateTime)
	if err != nil {
		return
	}
	fr.CreateTime = fr_create_time

	fr_last_submitted_time, err := time.Parse(time.RFC3339, response.LastSubmittedTime)
	if err != nil {
		return
	}
	fr.LastSubmittedTime = fr_last_submitted_time

	return db.Transaction(func(db_tx *gorm.DB) (err error) {
		defer alog.Trace("begin transaction").Stop(&err)
		tx := db_tx.Model(&model.FormResponse{}).FirstOrCreate(&fr, model.FormResponse{ResponseId: fr.ResponseId, FormId: fr.FormId})
		err = tx.Error
		if err != nil {
			return
		}
		alog.WithField("rows_affected", tx.RowsAffected).WithField("is_created", tx.RowsAffected > 0).Info("is response created")
		for k, v := range response.Answers {
			var textanswers []byte
			if v.TextAnswers != nil {
				textanswers, err = v.TextAnswers.MarshalJSON()
				if err != nil {
					return
				}
			} else {
				textanswers = []byte("[]")
			}
			a := model.Answer{
				QuestionId:     v.QuestionId,
				FormId:         response.FormId,
				FormResponseId: response.ResponseId,
				TextAnswers:    string(textanswers),
			}
			answers[k] = a
			tx := db_tx.Model(&model.Answer{}).FirstOrCreate(&a, &model.Answer{
				QuestionId:     v.QuestionId,
				FormId:         response.FormId,
				FormResponseId: response.ResponseId,
				TextAnswers:    string(textanswers),
			})
			err = tx.Error
			if err != nil {
				return
			}
			alog.
				WithField("question_id", a.QuestionId).
				WithField("rows_affected", tx.RowsAffected).
				WithField("is_created", tx.RowsAffected > 0).
				Info("is answer created")
		}
		return
	})
}

func ListFormSynchronizeDB(db *gorm.DB, forms_service *forms.Service, form_id string) (err error) {
	defer alog.WithField("form_id", form_id).Trace("ListFormSynchronizeDB").Stop(&err)

	form_get_call := forms_service.Forms.Get(form_id)

	form_get_call_response, err := form_get_call.Do()
	if err != nil {
		return
	}

	f := model.FormToLocalSchema(*form_get_call_response)

	tx := db.Omit("Items").FirstOrCreate(&f, model.Form{FormId: f.FormId})
	err = tx.Error
	if err != nil {
		return
	}
	alog.WithField("rows_affected", tx.RowsAffected).WithField("is_created", tx.RowsAffected > 0).Info("is form created in db")

	for _, item := range f.Items {
		tx := db.FirstOrCreate(&item, model.Item{QuestionId: item.QuestionId})
		err = tx.Error
		if err != nil {
			return
		}
		alog.WithField("rows_affected", tx.RowsAffected).WithField("is_created", tx.RowsAffected > 0).Info("is form items created")
	}
	return
}
