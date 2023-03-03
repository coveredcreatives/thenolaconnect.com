package internal

import (
	"net/http"
	"os/exec"

	alog "github.com/apex/log"
	"gitlab.com/the-new-orleans-connection/qr-code/devtools"
	"gitlab.com/the-new-orleans-connection/qr-code/model"
	"gorm.io/gorm"
)

func OrderToPrinterChannel(db *gorm.DB, order_id int) (order model.Order, err error) {
	alog.WithField("order_id", order_id).Trace("OrderToPrinterChannel").Stop(&err)
	order_communication_config, err := devtools.OrderCommunicationConfigFromEnv()
	if err != nil {
		return
	}

	tx := db.Model(&model.Order{OrderId: order_id}).First(&order)
	err = tx.Error
	if err != nil {
		return
	}
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, order.OrderFileStorageURL, nil)
	if err != nil {
		return
	}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	alog.WithField("status", res.Status).WithField("order_file_storage_url", order.OrderFileStorageURL).Info("request for download")

	cmd := exec.Command("lpr", "-H", order_communication_config.EnvPrinterIPv4Address)
	cmd.Stdin = res.Body

	err = cmd.Run()
	if err != nil {
		return
	}

	tx = db.Model(&model.Order{OrderId: order_id}).Update("is_delivered_to_kitchen", true)
	err = tx.Error
	if err != nil {
		return
	}
	return
}
