package route

import (
	"NotificationService/internal/repo"
	"NotificationService/internal/utils"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestSuccessCreateClientHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	savedId := "CLIENT_ID"
	rWriter := utils.NewMockResponseWriter(ctrl)
	rWriter.EXPECT().WriteHeader(http.StatusOK)
	rWriter.EXPECT().Write([]byte(savedId))

	clientRepo := repo.NewMockClientRepoInterface(ctrl)
	clientRepo.EXPECT().IsClientExist(gomock.Any()).Return(false)
	clientRepo.EXPECT().CreateClient(gomock.Any()).Return(savedId)

	ClientService = &RouteClientService{clientRepo}
	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(
		`{
			"phoneNumber":"71234567890",
			"opCode":"123",
			"tag":"Dr. Livesey",
			"utc":"+03:00"
		}`,
	))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	ClientService.createClientHandler(rWriter, req)
}

func TestFailedCreateClientHandlerAccountExist(t *testing.T) {
	ctrl := gomock.NewController(t)

	rWriter := utils.NewMockResponseWriter(ctrl)
	rWriter.EXPECT().WriteHeader(http.StatusOK)
	rWriter.EXPECT().Write([]byte(ErrClientAlreadyExists.Error()))

	clientRepo := repo.NewMockClientRepoInterface(ctrl)
	clientRepo.EXPECT().IsClientExist(gomock.Any()).Return(true)

	ClientService = &RouteClientService{clientRepo}
	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(
		`{
			"phoneNumber":"71234567890",
			"opCode":"123",
			"tag":"Dr. Livesey",
			"utc":"+03:00"
		}`,
	))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	ClientService.createClientHandler(rWriter, req)
}

func TestFailedCreateClientHandlerInvalidClientInfo(t *testing.T) {
	ctrl := gomock.NewController(t)

	rWriter := utils.NewMockResponseWriter(ctrl)
	rWriter.EXPECT().WriteHeader(http.StatusOK)
	rWriter.EXPECT().Write([]byte(ErrInvalidClientInfo.Error()))

	clientRepo := repo.NewMockClientRepoInterface(ctrl)

	ClientService = &RouteClientService{clientRepo}
	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(
		`{
			"phoneNumber":"71234567890",
			"opCode":"",
			"tag":"Dr. Livesey",
			"utc":"+03:00"
		}`,
	))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	ClientService.createClientHandler(rWriter, req)
}

func TestDeleteClientHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)

	vars := make(map[string]string)
	vars["phoneNumber"] = "71234567890"
	r, err := http.NewRequest(http.MethodDelete, "", nil)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	r = mux.SetURLVars(r, vars)

	clientRepo := repo.NewMockClientRepoInterface(ctrl)
	clientRepo.EXPECT().DeleteClient("71234567890")

	clientService := &RouteClientService{clientRepo}

	clientService.deleteClientHandler(w, r)
}

func TestSuccessModifyClientHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	w := utils.NewMockResponseWriter(ctrl)
	w.EXPECT().WriteHeader(http.StatusOK)

	vars := make(map[string]string)
	vars["id"] = "123"
	r, err := http.NewRequest(http.MethodPut, "", strings.NewReader(
		`{
			"phoneNumber":"71234567890",
			"opCode":"123",
			"tag":"Dr. Livesey",
			"utc":"+03:00"
		}`,
	))
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	r = mux.SetURLVars(r, vars)

	clientRepo := repo.NewMockClientRepoInterface(ctrl)
	clientRepo.EXPECT().UpdateClient(gomock.Any(), gomock.Any())
	clientRepo.EXPECT().IsClientExist(gomock.Any()).Return(true)

	clientService := &RouteClientService{clientRepo}

	clientService.modifyClientHandler(w, r)
}
