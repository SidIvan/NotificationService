package route

import (
	"NotificationService/internal/dto"
	"NotificationService/internal/repo"
	"NotificationService/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
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
	startAt := "1980-09-01T15:00:00+03:00"
	message := "test message"
	opCode := "967"
	tag := "Dr. Livesey"
	endAt := "3000-09-01T15:00:00+03:00"
	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)
	w.EXPECT().Write([]byte(savedId))

	distributionRepo := repo.NewMockDistributionRepoInterface(ctrl)
	distributionRepo.EXPECT().CreateDistribution(gomock.Any()).Return(savedId)

	var wg sync.WaitGroup
	wg.Add(1)
	messageSender := NewMockMessageSenderInterface(ctrl)
	messageSender.EXPECT().SendDistribution(gomock.Any()).Do(
		func(distributionWithId *dto.DistributionWithId) {
			expectedDisributionInfo := dto.Distribution{
				StartAt: startAt,
				Message: message,
				Filter: dto.Filter{
					OpCode: opCode,
					Tag:    tag,
				},
				EndAt: endAt,
			}
			if *distributionWithId != *(&expectedDisributionInfo).WithId(savedId) {
				t.Errorf("wrong distribution payload")
			}
			wg.Done()
		},
	)

	r, err := http.NewRequest(http.MethodPost, "", strings.NewReader(fmt.Sprintf(
		`{
			"StartAt":"%s",
			"Message":"%s",
			"Filter":{
				"OpCode":"%s",
				"Tag":"%s"
			},
			"EndAt":"%s"
		}`, startAt, message, opCode, tag, endAt,
	)))
	if err != nil {
		t.Errorf(err.Error())
	}

	distributionService := &RouteDistributionService{distributionRepo, messageSender}

	distributionService.createDistributionHandler(w, r)
	wg.Wait()
}

func TestFailedCreateDistributionHandlerWrongFormat(t *testing.T) {
	ctrl := gomock.NewController(t)
	startAt := "1980-09-01T15:00:00+03:00"
	message := "test message"
	opCode := "invalidOpCode"
	tag := "Dr. Livesey"
	endAt := "3000-09-01T15:00:00+03:00"

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)
	w.EXPECT().Write([]byte(ErrInvalidDistributionFormat.Error()))

	distributionRepo := repo.NewMockDistributionRepoInterface(ctrl)

	r, err := http.NewRequest(http.MethodPost, "", strings.NewReader(fmt.Sprintf(
		`{
			"StartAt":"%s",
			"Message":"%s",
			"Filter":{
				"OpCode":"%s",
				"Tag":"%s"
			},
			"EndAt":"%s"
		}`, startAt, message, opCode, tag, endAt,
	)))
	if err != nil {
		t.Errorf(err.Error())
	}

	distributionService := &RouteDistributionService{distributionRepo, nil}

	distributionService.createDistributionHandler(w, r)
}

var (
	distrId                     = "DISTR_ID"
	startAt                     = "1980-09-01T15:00:00+03:00"
	message                     = "test message"
	opCode                      = "967"
	tag                         = "Dr. Livesey"
	endAt                       = "3000-09-01T15:00:00+03:00"
	expectedSuccessDistribution = dto.Distribution{
		StartAt: startAt,
		Message: message,
		Filter: dto.Filter{
			OpCode: opCode,
			Tag:    tag,
		},
		EndAt: endAt,
	}
)

func TestSuccessModifyDistributionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)

	distributionRepo := repo.NewMockDistributionRepoInterface(ctrl)
	distributionRepo.EXPECT().IsDistributionExist(distrId).Return(true)
	distributionRepo.EXPECT().UpdateDistribution(distrId, &expectedSuccessDistribution)

	r, err := http.NewRequest(http.MethodPut, "", strings.NewReader(fmt.Sprintf(
		`{
			"StartAt":"%s",
			"Message":"%s",
			"Filter":{
				"OpCode":"%s",
				"Tag":"%s"
			},
			"EndAt":"%s"
		}`, startAt, message, opCode, tag, endAt,
	)))
	if err != nil {
		t.Errorf(err.Error())
	}
	vars := make(map[string]string)
	vars["id"] = distrId
	r = mux.SetURLVars(r, vars)

	distributionService := &RouteDistributionService{distributionRepo, nil}

	distributionService.modifyDistributionHandler(w, r)
}

func TestFailedModifyDistributionHandlerNonExistingId(t *testing.T) {
	ctrl := gomock.NewController(t)

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)
	w.EXPECT().Write([]byte(ErrNonExistenseDistribution.Error()))

	distributionRepo := repo.NewMockDistributionRepoInterface(ctrl)
	distributionRepo.EXPECT().IsDistributionExist(distrId).Return(false)

	r, err := http.NewRequest(http.MethodPost, "", strings.NewReader(fmt.Sprintf(
		`{
			"StartAt":"%s",
			"Message":"%s",
			"Filter":{
				"OpCode":"%s",
				"Tag":"%s"
			},
			"EndAt":"%s"
		}`, startAt, message, opCode, tag, endAt,
	)))
	if err != nil {
		t.Errorf(err.Error())
	}
	vars := make(map[string]string)
	vars["id"] = distrId
	r = mux.SetURLVars(r, vars)

	distributionService := &RouteDistributionService{distributionRepo, nil}

	distributionService.modifyDistributionHandler(w, r)
}

func TestDeleteDistributionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)

	distributionRepo := repo.NewMockDistributionRepoInterface(ctrl)
	distributionRepo.EXPECT().DeleteDistribution(distrId)

	r, err := http.NewRequest(http.MethodPost, "", nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	vars := make(map[string]string)
	vars["id"] = distrId
	r = mux.SetURLVars(r, vars)

	distributionService := &RouteDistributionService{distributionRepo, nil}

	distributionService.deleteDistributionHandler(w, r)
}

func TestHandleDistributionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)

	messageSender := NewMockMessageSenderInterface(ctrl)
	messageSender.EXPECT().HandleDistributions()

	distributionService := &RouteDistributionService{nil, messageSender}
	r, err := http.NewRequest(http.MethodPatch, "", nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	distributionService.handleDistributionHandler(w, r)
}

func TestSuccessSingleInfoDistributionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)
	expectedResponseBody, _ := json.Marshal(expectedSuccessDistribution.WithId(distrId))
	w.EXPECT().Write(expectedResponseBody).Return(len(expectedResponseBody), nil)
	w.EXPECT().Header().Return(make(http.Header))

	distributionRepo := repo.NewMockDistributionRepoInterface(ctrl)
	distributionRepo.EXPECT().FindDistributionById(distrId).Return(expectedSuccessDistribution.WithId(distrId))

	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	vars := make(map[string]string)
	vars["id"] = distrId
	r = mux.SetURLVars(r, vars)

	distributionService := &RouteDistributionService{distributionRepo, nil}

	distributionService.singleInfoDistributionHandler(w, r)
}

func TestSuccessSingleInfoDistributionHandlerNotExistenceDistribution(t *testing.T) {
	ctrl := gomock.NewController(t)

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)

	distributionRepo := repo.NewMockDistributionRepoInterface(ctrl)
	distributionRepo.EXPECT().FindDistributionById(distrId).Return(nil)
	w.EXPECT().Write([]byte(ErrNonExistenseDistribution.Error()))

	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Errorf(err.Error())
	}
	vars := make(map[string]string)
	vars["id"] = distrId
	r = mux.SetURLVars(r, vars)

	distributionService := &RouteDistributionService{distributionRepo, nil}

	distributionService.singleInfoDistributionHandler(w, r)
}

func TestFullInfoDistributionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	id1, id2 := "507f1f77bcf86cd799439011", "507f191e810c19729de860ea"
	startAt2 := "1981-09-01T15:00:00+03:00"
	message2 := "test message2"
	opCode2 := "903"
	tag2 := "Cpt. Smollet"
	endAt2 := "3001-09-01T15:00:00+03:00"
	expectedDistribution2 := dto.Distribution{
		StartAt: startAt2,
		Message: message2,
		Filter: dto.Filter{
			OpCode: opCode2,
			Tag:    tag2,
		},
		EndAt: endAt2,
	}
	expectedResponse := []dto.DistributionWithId{
		*expectedSuccessDistribution.WithId(id1),
		*expectedDistribution2.WithId(id2),
	}
	expectedResponseBody, _ := json.Marshal(expectedResponse)

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)
	w.EXPECT().Write([]byte(expectedResponseBody)).Return(len(expectedResponseBody), nil)
	w.EXPECT().Header().Return(make(http.Header))

	distributionRepo := repo.NewMockDistributionRepoInterface(ctrl)
	distributionRepo.EXPECT().FindAllDistributions().Return(&expectedResponse)

	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Errorf(err.Error())
	}

	distributionService := &RouteDistributionService{distributionRepo, nil}

	distributionService.fullInfoDistributionHandler(w, r)
}
