package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/malfanmh/wabotapi/model"
	"io"
	"net/http"
)

type Finpay struct {
	client  *http.Client
	baseURL string
}

func NewFinpay(client *http.Client, baseURL string) *Finpay {
	return &Finpay{client, baseURL}
}

func (f *Finpay) GetPaymentLink(ctx context.Context, client model.Client, payload model.PaymentLink) (result model.PaymentLinkResponse, err error) {
	url := fmt.Sprintf("%s/pg/payment/card/initiate", f.baseURL)
	payload.URL.CallbackURL = client.FinpayCallbackURL
	b, err := json.Marshal(payload)
	if err != nil {
		err = fmt.Errorf("error marshalling payload: %w", err)
		return
	}
	fmt.Println("req URL :", url, string(b))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		err = fmt.Errorf("error creating request: %w", err)
		return
	}
	req.SetBasicAuth(client.FinpayMerchantID, client.FinpaySecret)
	req.Header.Set("Content-Type", "application/json")

	response, err := f.client.Do(req)
	if err != nil {
		err = fmt.Errorf("error executing request: %w", err)
		return
	}
	defer response.Body.Close()

	b, err = io.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("error reading response body: %w", err)
		return
	}
	fmt.Println("GetPaymentLink response", string(b))
	if err = json.Unmarshal(b, &result); err != nil {
		err = fmt.Errorf("error unmarshalling response body: %w", err)
		return
	}
	//if err = json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
	//	err = fmt.Errorf("error parsing finpy response: %w", err)
	//	return
	//}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("finpy service responded with status [code %d, msg :%v]", response.StatusCode, result.ResponseMessage)
		return
	}
	return
}
