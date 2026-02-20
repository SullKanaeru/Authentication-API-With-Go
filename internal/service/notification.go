package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
)

type NotificationService struct {
	SMTPHost    string
	SMTPPort    string
	SMTPUser    string
	SMTPPass    string
	SMTPSender  string
	FonnteToken string
}

func NewNotificationService(host, port, user, pass, sender, fonnteToken string) *NotificationService {
	return &NotificationService{
		SMTPHost:    host,
		SMTPPort:    port,
		SMTPUser:    user,
		SMTPPass:    pass,
		SMTPSender:  sender,
		FonnteToken: fonnteToken,
	}
}

func (n *NotificationService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", n.SMTPUser, n.SMTPPass, n.SMTPHost)
	
	msg := []byte(fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", n.SMTPSender, n.SMTPUser, to, subject, body))
	
	addr := fmt.Sprintf("%s:%s", n.SMTPHost, n.SMTPPort)
	return smtp.SendMail(addr, auth, n.SMTPUser, []string{to}, msg)
}

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