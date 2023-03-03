package mocks

import (
	"github.com/googleapis/google-cloud-go-testing/storage/stiface"
)

type MockStorage struct {
	stiface.Client
}
