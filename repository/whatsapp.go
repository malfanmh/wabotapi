package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/malfanmh/wabotapi/model"
	"io"
	"net/http"
)

type whatsappAPI struct {
	baseURL           string
	businessAccountID string
	token             string
	c                 *http.Client
}

func NewWhatsappAPI(url, id, token string, client *http.Client) *whatsappAPI {
	return &whatsappAPI{
		baseURL:           url,
		businessAccountID: id,
		token:             token,
		c:                 client,
	}
}

func (wa *whatsappAPI) SendTemplate(ctx context.Context, from, to string, template model.WATemplate, params map[string]interface{}) (string, error) {
	url := fmt.Sprintf("%s/%s/messages", wa.baseURL, from)

	data := fmt.Sprintf(`{
        "messaging_product": "whatsapp",
        "to": "%s",
        "type": "template",
        "template": {
            "name": "%s",
            "language": {
                "code": "%s"
            }
        }
    }`, to, template.Name, template.Language)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+wa.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := wa.c.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println(string(b))
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status code %d", resp.StatusCode)
	}
	return string(b), nil
}

func (wa *whatsappAPI) SendText(ctx context.Context, from, to string, text string, params map[string]interface{}) (string, error) {
	url := fmt.Sprintf("%s/%s/messages", wa.baseURL, from)

	data := fmt.Sprintf(`{
		"messaging_product": "whatsapp",
		"recipient_type": "individual",
		"to": "%s",
		"type": "text",
		"text": {
			"body": "%s"
		}
	}`, to, text)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+wa.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := wa.c.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status code %d", resp.StatusCode)
	}
	return string(b), nil
}
