package route

import (
	"NotificationService/internal/dto"
	"NotificationService/internal/repo"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	ErrClientAlreadyExists = errors.New("client already exists")
	ErrInvalidClientInfo   = errors.New("invalid client info")
)

var ClientService *RouteClientService

func NewClientRouter(r *mux.Router) {
	ClientService = &RouteClientService{ClientRepo: repo.ClientRepo}
	r = r.PathPrefix("/client").Subrouter()
	r.
		HandleFunc("/create", ClientService.createClientHandler).
		Headers("Content-Type", "application/json").
		Methods(http.MethodPost)
	r.
		HandleFunc("/delete/{phoneNumber}", ClientService.deleteClientHandler).
		Methods(http.MethodDelete)
	r.
		HandleFunc("/modify/{id}", ClientService.modifyClientHandler).
		Headers("Content-type", "application/json").
		Methods(http.MethodPut)
}

type RouteClientService struct {
	ClientRepo repo.ClientRepoInterface
}

//	@Tags			client
//	@Description	Modify client info
//	@Router			/client/modify/{id} [put]
//	@ID				modify-client
//	@Accept			json
//	@Param			input	body		dto.Client	true	"client info"
//	@Param			id		path		string		true	"client id"
//	@Success		200		{string}	string		"error message if failure"
//	@failure		500
func (s *RouteClientService) modifyClientHandler(w http.ResponseWriter, r *http.Request) {
	var client dto.Client
	if !parseBody(w, r, &client) {
		return
	}
	w.WriteHeader(http.StatusOK)
	if !client.IsValid() {
		w.Write([]byte("invalid client info"))
		return
	}
	if !s.ClientRepo.IsClientExist(client.PhoneNumber) {
		w.Write([]byte("client does not exist"))
		return
	}
	id := mux.Vars(r)["id"]
	s.ClientRepo.UpdateClient(id, &client)
}

//	@Tags			client
//	@Description	Create new client
//	@Router			/client/create [post]
//	@ID				create-client
//	@Param			input	body		dto.Client	true	"client info"
//	@Success		200		{string}	string		"error message if failure"
func (s *RouteClientService) createClientHandler(w http.ResponseWriter, r *http.Request) {
	var newClient dto.Client
	if !parseBody(w, r, &newClient) {
		return
	}
	w.WriteHeader(http.StatusOK)
	if !newClient.IsValid() {
		w.Write([]byte(ErrInvalidClientInfo.Error()))
		return
	}
	if s.ClientRepo.IsClientExist(newClient.PhoneNumber) {
		w.Write([]byte(ErrClientAlreadyExists.Error()))
		return
	}
	w.Write([]byte(s.ClientRepo.CreateClient(&newClient)))
}

//	@Tags			client
//	@Description	Delete client
//	@Router			/client/delete/{phoneNumber} [delete]
//	@ID				delete-client
//	@Param			id	path	int	true	"phoneNumber"
//	@Success		200
func (s *RouteClientService) deleteClientHandler(w http.ResponseWriter, r *http.Request) {
	phoneNumber := mux.Vars(r)["phoneNumber"]
	s.ClientRepo.DeleteClient(phoneNumber)
	w.WriteHeader(http.StatusOK)
}
