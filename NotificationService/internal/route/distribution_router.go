package route

import (
	"NotificationService/internal/dto"
	"NotificationService/internal/repo"
	"NotificationService/internal/utils"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	authToken           string
	DistributionService *RouteDistributionService
)

var (
	ErrInvalidDistributionFormat = errors.New("invalid disctribution info format")
	ErrNonExistenseDistribution  = errors.New("distribution does not exist")
)

func NewDistributionRouter(r *mux.Router) {
	MessageSender = &MessageSenderImpl{repo.MessageRepo}
	DistributionService = &RouteDistributionService{repo.DistributionRepo, MessageSender}
	r = r.PathPrefix("/distribution").Subrouter()
	r.
		HandleFunc("/create", DistributionService.createDistributionHandler).
		Headers("Content-Type", "application/json").
		Methods(http.MethodPost)
	r.
		HandleFunc("/full-info", DistributionService.fullInfoDistributionHandler).
		Methods(http.MethodGet)
	r.
		HandleFunc("/info/{id}", DistributionService.singleInfoDistributionHandler).
		Methods(http.MethodGet)
	r.
		HandleFunc("/modify/{id}", DistributionService.modifyDistributionHandler).
		Headers("Content-Type", "application/json").
		Methods(http.MethodPut)
	r.
		HandleFunc("/delete/{id}", DistributionService.deleteDistributionHandler).
		Methods(http.MethodDelete)
	r.
		HandleFunc("/handle", DistributionService.handleDistributionHandler).
		Methods(http.MethodPatch)
	authToken = utils.PMan.Get("fbrq_token").(string)
}

type RouteDistributionService struct {
	DistributionRepo repo.DistributionRepoInterface
	MessageSender    MessageSenderInterface
}

// @Tags			distribution
// @Description	Create distribution
// @Router			/distribution/create [post]
// @ID				create-distribution
// @Accept			json
// @Param			input	body		dto.Distribution	true	"distribution info"
// @Success		200		{string}	string				"error message if failure"
// @failure		500
func (s *RouteDistributionService) createDistributionHandler(w http.ResponseWriter, r *http.Request) {
	var distribution dto.Distribution
	if !parseBody(w, r, &distribution) {
		return
	}
	w.WriteHeader(http.StatusOK)
	if !distribution.IsValid() {
		w.Write([]byte(ErrInvalidDistributionFormat.Error()))
		return
	}
	id := s.DistributionRepo.CreateDistribution(&distribution)
	if distribution.SholdSend() {
		go s.MessageSender.SendDistribution(distribution.WithId(id))
	}
	w.Write([]byte(id))
}

// @Tags			distribution
// @Description	Get info about all distributions
// @Router			/distribution/full-info [get]
// @ID				get-all-distributions
// @Success		200	{array}	dto.DistributionWithId
// @failure		502
func (s *RouteDistributionService) fullInfoDistributionHandler(w http.ResponseWriter, r *http.Request) {
	distributions := s.DistributionRepo.FindAllDistributions()
	if distributions == nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	body, err := json.Marshal(distributions)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Println(err)
		return
	}
	w.Header().Add("Content-type", "application/json")
	numBytes, err := w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Println(err)
		return
	}
	if numBytes != len(body) {
		w.WriteHeader(http.StatusBadGateway)
		log.Println("not full response sended")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Tags			distribution
// @Description	Get info about distribution
// @Router			/distribution/info/{id} [get]
// @ID				get-single-distribution
// @Param			id	path		int	true	"distribution id"
// @Success		200	{object}	dto.DistributionWithId
// @failure		502
func (s *RouteDistributionService) singleInfoDistributionHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	distribution := s.DistributionRepo.FindDistributionById(id)
	if distribution == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ErrNonExistenseDistribution.Error()))
		return
	}
	body, err := json.Marshal(*distribution)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Println(err)
		return
	}
	w.Header().Add("Content-type", "application/json")
	numBytes, err := w.Write(body)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		log.Println(err)
		return
	}
	if numBytes != len(body) {
		w.WriteHeader(http.StatusBadGateway)
		log.Println("not full response sended")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Tags			distribution
// @Description	Modify distribution
// @Router			/distribution/modify/{id} [put]
// @ID				modify-distribution
// @Param			id		path	string				true	"distribution id"
// @Param			input	body	dto.Distribution	true	"distribution info"
// @Success		200
// @failure		502
func (s *RouteDistributionService) modifyDistributionHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	w.WriteHeader(http.StatusOK)
	if !s.DistributionRepo.IsDistributionExist(id) {
		w.Write([]byte(ErrNonExistenseDistribution.Error()))
		return
	}
	var distribution dto.Distribution
	if !parseBody(w, r, &distribution) {
		return
	}
	s.DistributionRepo.UpdateDistribution(id, &distribution)
}

// @Tags			distribution
// @Description	Delete distribution
// @Router			/distribution/delete/{id} [delete]
// @ID				delete-distribution
// @Param			id	path	string	true	"distribution id"
// @Success		200
func (s *RouteDistributionService) deleteDistributionHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	s.DistributionRepo.DeleteDistribution(id)
	w.WriteHeader(http.StatusOK)
}

// @Tags			distribution
// @Description	handle distribution
// @Router			/distribution/handle [patch]
// @ID				handle-distributions
// @Success		200
func (s *RouteDistributionService) handleDistributionHandler(w http.ResponseWriter, r *http.Request) {
	s.MessageSender.HandleDistributions()
	w.WriteHeader(http.StatusOK)
}
