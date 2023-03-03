package pgw

import (
	"fmt"
	"io"
	"net/http"

	alog "github.com/apex/log"
	"gitlab.com/the-new-orleans-connection/qr-code/devtools"
	"gitlab.com/the-new-orleans-connection/qr-code/model"
)

func TriggerPDFConversionApi(pgw_env_config devtools.PDFGenerationWorkerServerConfig, pgw model.PDFGenerationWorker) (b []byte, err error) {
	alog.WithField("pdf_generation_worker_id", pgw.PDFGenerationWorkerId).Trace("TriggerPDFConversionApi").Stop(&err)
	client := http.Client{}
	url := fmt.Sprint(pgw_env_config.EnvBaseURL, "/pdf_generation_worker/", pgw.PDFGenerationWorkerId)
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	alog.
		WithField("url", url).
		WithField("status", resp.StatusCode).
		Info("response recieved")

	b, err = io.ReadAll(resp.Body)
	return
}
