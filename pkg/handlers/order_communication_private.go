package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	alog "github.com/apex/log"
	internal_tools "github.com/coveredcreatives/thenolaconnect.com/pkg/internal"
	"github.com/coveredcreatives/thenolaconnect.com/pkg/model"
	"github.com/twilio/twilio-go"
	"gorm.io/gorm"
)

func listOrders(w io.Writer, r *http.Request, db *gorm.DB) (b []byte, err error) {
	defer alog.Trace("listOrders").Stop(&err)

	var orders []model.Order
	tx := db.
		Model(&model.Order{}).
		Where("is_pdf_generated = ?", true).
		Find(&orders)
	err = tx.Error
	if err != nil {
		return
	}
	alog.WithField("rows_affected", tx.RowsAffected).Info("find orders with pdf generated")
	b, err = json.Marshal(&orders)
	return
}

func generateOrder(w io.Writer, r *http.Request, db *gorm.DB, storagec *storage.Client, twilioc *twilio.RestClient) (b []byte, err error) {
	defer alog.Trace("generateOrder").Stop(&err)
	var valid_count int64
	tx := db.Model(&model.FormResponse{}).Where("response_id NOT IN (?)", db.Model(&model.Order{}).Select("form_response_id")).Count(&valid_count)
	err = tx.Error
	if err != nil {
		return
	}
	errchan := make(chan error, valid_count)
	defer close(errchan)
	responsechan := make(chan model.FormResponse, valid_count)
	orderchan := make(chan model.Order, valid_count)
	pgwchan := make(chan model.PDFGenerationWorker, valid_count)
	donechan := make(chan bool)
	defer close(donechan)

	go internal_tools.QueryFormResponsesWithoutAssociatedOrder(db, responsechan, errchan)

	go internal_tools.CreateOrder(db, storagec, twilioc, responsechan, orderchan, errchan)

	go internal_tools.CreatePDFGenerationWorkers(db, orderchan, pgwchan, errchan)

	go internal_tools.StreamPDFStoStorage(db, storagec, pgwchan, donechan)

	for {
		select {
		case err = <-errchan:
			return
		case <-donechan:
			_, _, err = internal_tools.CreateOrderConversation(db, twilioc)
			if err != gorm.ErrRecordNotFound && err != nil {
				return
			}
			if err == gorm.ErrRecordNotFound {
				err = nil
				return
			}
			err = internal_tools.CreateOrderMessage(db, twilioc)
			if err == gorm.ErrRecordNotFound {
				err = nil
			}
			return
		default:
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}

func deliverOrderToKitchen(w io.Writer, r *http.Request, db *gorm.DB, twilioc *twilio.RestClient, orderchan chan<- int) (b []byte, err error) {
	defer alog.Trace("deliverOrderToKitchen").Stop(&err)

	err = r.ParseForm()
	if err != nil {
		return
	}
	order_id, err := strconv.ParseInt(r.FormValue("order_id"), 10, 64)
	if err != nil {
		return
	}

	internal_tools.DeliverOrderToKitchen(db, twilioc, int(order_id), orderchan)
	return
}

func sms(w io.Writer, r *http.Request, db *gorm.DB, twilioc *twilio.RestClient, orderchan chan<- int) (b []byte, err error) {
	defer alog.Trace("sms").Stop(&err)

	err = r.ParseForm()
	if err != nil {
		return
	}

	to := r.FormValue("To")
	body := r.FormValue("Body")

	order, err := internal_tools.IdentifyOrderFromPhoneNumberLocal(db, to)
	if err != nil {
		return
	}

	err = internal_tools.ValidateOrder(db, twilioc, order.OrderId, body, orderchan)
	if err != nil {
		return
	}

	err = internal_tools.CreateOrderMessage(db, twilioc)

	return
}
