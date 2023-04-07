package internal

import (
	"time"

	alog "github.com/apex/log"
	"github.com/coveredcreatives/thenolaconnect.com/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreatePDFGenerationWorkers(db *gorm.DB, orderchan chan model.Order, pgwchan chan<- model.PDFGenerationWorker, errchan chan<- error) {
	defer alog.Trace("generateOrdersPDFGenerationWorkers").Stop(nil)
	now := time.Now()
	for o := range orderchan {
		order := o
		alog.WithField("order.OrderId", order.OrderId).Info("receive order from order channel")
		pgw := model.PDFGenerationWorker{
			FormId:  order.FormId,
			OrderId: order.OrderId,
			StartAt: now,
		}

		tx := db.
			Clauses(clause.Returning{Columns: []clause.Column{{Table: "order_communication.pdf_generation_worker", Name: "pdf_generation_worker_id"}}}).
			Omit("pdf_generation_worker_id").
			FirstOrCreate(&pgw, &model.PDFGenerationWorker{OrderId: order.OrderId})
		if tx.Error != nil {
			errchan <- tx.Error
			return
		}
		alog.WithField("rows_affected", tx.RowsAffected).
			WithField("pgw.PDFGenerationWorkerId", pgw.PDFGenerationWorkerId).
			Info("create pdf generation worker for appropriate orders in database")

		tx = db.
			Model(&model.Order{}).
			Where(&model.Order{OrderId: order.OrderId}).
			Update("pdf_generation_worker_id", pgw.PDFGenerationWorkerId)
		if tx.Error != nil {
			errchan <- tx.Error
			return
		}

		alog.WithField("rows_affected", tx.RowsAffected).
			WithField("order.OrderId", order.OrderId).
			WithField("order.PDFGenerationWorkerId", order.PDFGenerationWorkerId).
			Info("update order with PDFGenerationWorkerId in database")

		pgwchan <- pgw
	}
	close(pgwchan)
}
