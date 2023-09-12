package route

import (
	"NotificationService/internal/repo"
	"NotificationService/internal/utils"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSuccessCreateDistributionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	savedId := "CLIENT_ID"
	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)
	w.EXPECT().Write([]byte(savedId))

	distributionRepo := repo.NewMockDistributionRepoInterface(ctrl)
	distributionRepo.EXPECT().CreateDistribution(gomock.Any()).Return(savedId)

	r, err := http.NewRequest(http.MethodPost, "", strings.NewReader(
		`{
			"StartAt":"3000-09-01T15:00:00+03:00",
			"Message":"AXAXAXAXA",
			"Filter":{
				"OpCode":"967",
				"Tag":"Dr. Livesey"
			},
			"EndAt":"3001-09-01T15:00:00+03:00"
		}`,
	))
	if err != nil {
		t.Errorf(err.Error())
	}

	distributionService := &RouteDistributionService{distributionRepo, nil}

	distributionService.createDistributionHandler(w, r)
}

func TestSuccesCreateDistributionHandlerRunSendDistribution(t *testing.T) {
	ctrl := gomock.NewController(t)
	savedId := "CLIENT_ID"
	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)
	w.EXPECT().Write([]byte(savedId))

	distributionRepo := repo.NewMockDistributionRepoInterface(ctrl)
	distributionRepo.EXPECT().CreateDistribution(gomock.Any()).Return(savedId)

	messageSender := NewMockMessageSenderInterface(ctrl)
	messageSender.EXPECT().SendDistribution(gomock.Any())

	r, err := http.NewRequest(http.MethodPost, "", strings.NewReader(
		`{
			"StartAt":"1980-09-01T15:00:00+03:00",
			"Message":"test message",
			"Filter":{
				"OpCode":"967",
				"Tag":"Dr. Livesey"
			},
			"EndAt":"3000-09-01T15:00:00+03:00"
		}`,
	))
	if err != nil {
		t.Errorf(err.Error())
	}

	distributionService := &RouteDistributionService{distributionRepo, messageSender}

	distributionService.createDistributionHandler(w, r)
}
