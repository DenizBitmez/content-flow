package services

import (
	"bytes"
	"content-flow/internal/database"
	"content-flow/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

// TriggerWebhooks runs asynchronously to avoid blocking the API response
func TriggerWebhooks(event string, payload interface{}) {
	go func() {
		var webhooks []models.Webhook
		// content.create -> %content.create% logic or simple check
		if err := database.DB.Where("enabled = ?", true).Find(&webhooks).Error; err != nil {
			log.Println("Error fetching webhooks:", err)
			return
		}

		for _, wh := range webhooks {
			if shouldTrigger(wh.Events, event) {
				sendWebhook(wh.URL, event, payload)
			}
		}
	}()
}

func shouldTrigger(registeredEvents, currentEvent string) bool {
	// If registered events is empty or "*", trigger all
	if registeredEvents == "" || registeredEvents == "*" {
		return true
	}
	events := strings.Split(registeredEvents, ",")
	for _, e := range events {
		if strings.TrimSpace(e) == currentEvent {
			return true
		}
	}
	return false
}

func sendWebhook(url, event string, payload interface{}) {
	body := map[string]interface{}{
		"event":     event,
		"timestamp": time.Now().Unix(),
		"data":      payload,
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println("Failed to create webhook request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "ContentFlow-CMS-Webhook")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to deliver webhook to %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("Webhook failed for %s with status %d\n", url, resp.StatusCode)
	} else {
		log.Printf("Webhook delivered to %s for event %s\n", url, event)
	}
}

// CRUD for Webhooks

func CreateWebhook(url, events string) (*models.Webhook, error) {
	webhook := &models.Webhook{
		URL:     url,
		Events:  events,
		Enabled: true,
	}
	err := database.DB.Create(webhook).Error
	return webhook, err
}

func GetAllWebhooks() ([]models.Webhook, error) {
	var webhooks []models.Webhook
	err := database.DB.Find(&webhooks).Error
	return webhooks, err
}
