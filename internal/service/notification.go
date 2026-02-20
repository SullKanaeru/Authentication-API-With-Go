package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type NotificationService struct {
	BrevoAPIKey string
	SenderEmail string
	SenderName  string
	FonnteToken string
}

func NewNotificationService(brevoKey, senderEmail, senderName, fonnteToken string) *NotificationService {
	return &NotificationService{
		BrevoAPIKey: brevoKey,
		SenderEmail: senderEmail,
		SenderName:  senderName,
		FonnteToken: fonnteToken,
	}
}

// Logika Kirim Email via Brevo API (Anti-Blokir)
func (n *NotificationService) SendEmail(to, subject, body string) error {
	url := "https://api.brevo.com/v3/smtp/email"

	// Format payload sesuai standar Brevo
	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  n.SenderName,
			"email": n.SenderEmail,
		},
		"to": []map[string]string{
			{
				"email": to,
			},
		},
		"subject":     subject,
		"textContent": body,
	}

	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("accept", "application/json")
	req.Header.Set("api-key", n.BrevoAPIKey)
	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Brevo mengembalikan status 201 Created jika sukses
	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error dari brevo: %s", string(bodyBytes))
	}

	return nil
}

// Logika Kirim WA via Fonnte (Tetap Sama)
func (n *NotificationService) SendWhatsApp(target, message string) error {
	url := "https://api.fonnte.com/send"

	payload := map[string]string{
		"target":  target,
		"message": message,
	}
	payloadBytes, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", n.FonnteToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error dari fonnte: %s", string(bodyBytes))
	}

	return nil
}