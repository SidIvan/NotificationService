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
	DistributionService RouteDistributionService
)

var (
	ErrInvalidDistributionFormat = errors.New("invalid disctribution info format")
)

func NewDistributionRouter(r *mux.Router) {
	MessageSender := &MessageSenderImpl{}
	DistributionService := RouteDistributionService{repo.DistributionRepo, MessageSender}
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
}

func (s *RouteDistributionService) singleInfoDistributionHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	distribution := s.DistributionRepo.FindDistributionById(id)
	if distribution == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Document not found"))
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
}

func (s *RouteDistributionService) modifyDistributionHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if !s.DistributionRepo.IsDistributionExist(id) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("distribution does not exist"))
		return
	}
	var distribution dto.Distribution
	if !parseBody(w, r, &distribution) {
		return
	}
	repo.DistributionRepo.UpdateDistribution(id, &distribution)
}

func (s *RouteDistributionService) deleteDistributionHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	s.DistributionRepo.DeleteDistribution(id)
	w.WriteHeader(http.StatusOK)
}

func (s *RouteDistributionService) handleDistributionHandler(w http.ResponseWriter, r *http.Request) {
	s.MessageSender.HandleDistributions()
	w.WriteHeader(http.StatusOK)
}
