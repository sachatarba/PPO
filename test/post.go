package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const (
	url          = "http://localhost:8080/api/v1/applications"
	token        = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50SWQiOiJhNjUyYTgzYi01ODdlLTQ2ZDMtYmYwMS00ZDJjNWIzYjUyY2YiLCJwZXJtaXNzaW9ucyI6IltcIk5vbmVcIixcIkRlZmF1bHRcIixcIkNhblJlYWRBbmRFZGl0U2VsZkFsbEFwcGxpY2F0aW9uc1wiXSIsImV4cCI6MTcyOTg3MDc5MCwiaXNzIjoiTXVuQ29udGVzdEJhY2tlbmQiLCJhdWQiOiJNdW5Db250ZXN0Q2xpZW50cyJ9.mwCY7RsD_VgIH9LiT2H0SKXSNmFSf7EEJnjFS-3FBPY"
	nominationID = "4034fa39-c06e-42b3-af5b-5b7993dd85d9"
)

type RequestBody struct {
	ApplicationID string `json:"applicationId"`
	NominationID  string `json:"nominationId"`
}

func sendRequest() error {
	// Генерация случайного UUID для applicationId
	appID := uuid.New().String()

	// Формирование тела запроса
	requestBody := RequestBody{
		ApplicationID: appID,
		NominationID:  nominationID,
	}
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	// Создание HTTP-запроса
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Установка заголовков
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Чтение ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Логирование ответа
	fmt.Printf("Response Status: %s\nResponse Body: %s\n", resp.Status, body)
	return nil
}

func main() {
	numRequests := 100
	// fmt.Print("Enter number of requests to send: ")
	// fmt.Scan(&numRequests)

	for i := 0; i < numRequests; i++ {
		// fmt.Printf("Sending request %d...\n", i+1)
		if err := sendRequest(); err != nil {
			log.Fatalf("Request %d failed: %v", i+1, err)
		}
		time.Sleep(500 * time.Millisecond) // Пауза между запросами
	}

	// fmt.Println("All requests sent successfully.")
}
