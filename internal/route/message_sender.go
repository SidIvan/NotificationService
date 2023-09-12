package route

import (
	"NotificationService/internal/dto"
	"NotificationService/internal/repo"
	"NotificationService/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var MessageSender *MessageSenderImpl

type MessageSenderImpl struct {
}

type MessageSenderInterface interface {
	SendDistribution(*dto.DistributionWithId)
	HandleDistributions()
}

func (s MessageSenderImpl) sendAndSaveStatistics(d *dto.DistributionWithId, c *dto.ClientWithId) {
	messageId := utils.GenerateId()
	phoneNumber, _ := strconv.Atoi(c.PhoneNumber)
	message := dto.Message{
		SendedAt:        dto.GetCurTime(),
		DisctributionId: d.GetId(),
		ClientId:        c.GetId(),
	}
	data, _ := json.Marshal(dto.SendInfo{
		Id:    messageId,
		Phone: phoneNumber,
		Text:  d.Message,
	})
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://probe.fbrq.cloud/v1/send/%d", messageId), bytes.NewReader(data))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Authorization", authToken)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		repo.CreateFailedMessage(&message)
		return
	}
	if resp.StatusCode != http.StatusOK {
		log.Printf("fbqr send request ended with code %d\n", resp.StatusCode)
		repo.CreateFailedMessage(&message)
	}
	repo.CreateSuccessMessage(&message)
}

func (s MessageSenderImpl) SendDistribution(d *dto.DistributionWithId) {
	clients := repo.ClientRepo.FindClientsByFilter(&d.Filter)
	if clients == nil {
		return
	}
	for _, client := range *clients {
		st := repo.GetStatus(d.Id, client.Id)
		if st == dto.OkStatus {
			return
		}
		go s.sendAndSaveStatistics(d, &client)
	}
}

func (s MessageSenderImpl) HandleDistributions() {
	distributions := repo.DistributionRepo.FindAllDistributions()
	for _, distribution := range *distributions {
		if distribution.SholdSend() {
			go s.SendDistribution(&distribution)
		}
	}
}

func HandleLoop() {
	for {
		time.Sleep(time.Second * 60)
		go MessageSender.HandleDistributions()
	}
}
