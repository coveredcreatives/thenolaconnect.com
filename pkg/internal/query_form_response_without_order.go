package internal

import (
	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/model"
	"gorm.io/gorm"
)

func QueryFormResponsesWithoutAssociatedOrder(db *gorm.DB, responsechan chan<- model.FormResponse, errchan chan<- error) {
	responses_without_associated_order := []model.FormResponse{}
	tx := db.Where("response_id NOT IN (?)", db.Model(&model.Order{}).Select("form_response_id")).Find(&responses_without_associated_order)
	if tx.Error != nil {
		errchan <- tx.Error
		return
	}
	alog.WithField("rows_affected", tx.RowsAffected).Info("find responses without associated order from database")
	for _, r := range responses_without_associated_order {
		response := r
		responsechan <- response
		alog.WithField("response_id", response.ResponseId).Info("deliver to channel")
	}
	close(responsechan)
}
