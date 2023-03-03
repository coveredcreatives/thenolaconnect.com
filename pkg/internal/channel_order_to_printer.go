package internal

import (
	"time"

	alog "github.com/apex/log"
	"gorm.io/gorm"
)

func ChannelOrdersToPrinter(db *gorm.DB, orderchan <-chan int) (err error) {
	for {
		select {
		case order_id := <-orderchan:
			alog.WithField("order_id", order_id).Info("recieve from order channel")
			_, err = OrderToPrinterChannel(db, order_id)
			if err != nil {
				alog.WithField("order_id", order_id).WithError(err).Error("failed to pipe order to printer")
				return
			}
		default:
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}
