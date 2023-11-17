package goklarna

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	disputeApiURL = "/disputes/v2/disputes"
)

type (
	DisputeSrv interface {
		GetLast250Disputes() ([]Dispute, error)
	}

	disputeSrv struct {
		client Client
	}

	DisputeResponse struct {
		Pagination struct {
			Limit    int      `json:"limit"`
			Count    int      `json:"count"`
			Total    int      `json:"total"`
			NextPage []string `json:"next_page"`
		} `json:"pagination"`
		Disputes []Dispute `json:"disputes"`
	}

	Dispute struct {
		DisputeKrn            string    `json:"dispute_krn"`
		Reason                string    `json:"reason"`
		OpenedAt              time.Time `json:"opened_at"`
		ClosedAt              time.Time `json:"closed_at"`
		ClosingReason         string    `json:"closing_reason"`
		ClosingReasonDetailed string    `json:"closing_reason_detailed,omitempty"`
		CaptureId             string    `json:"capture_id"`
		Region                string    `json:"region"`
		Order                 struct {
			OrderId            string    `json:"order_id"`
			CreatedAt          time.Time `json:"created_at"`
			OrderAmount        int       `json:"order_amount"`
			PurchaseCurrency   string    `json:"purchase_currency"`
			MerchantReference1 string    `json:"merchant_reference1"`
		} `json:"order"`
		Merchant struct {
			MerchantId string `json:"merchant_id"`
			Name       string `json:"name"`
		} `json:"merchant"`
		Requests []struct {
			RequestId               int           `json:"request_id"`
			CreatedAt               time.Time     `json:"created_at"`
			Comment                 string        `json:"comment"`
			OptionalRequestedFields []string      `json:"optional_requested_fields"`
			RequestedFields         []interface{} `json:"requested_fields"`
			Responses               []interface{} `json:"responses"`
			Attachments             []interface{} `json:"attachments"`
		} `json:"requests"`
		Status              string `json:"status"`
		InvestigationStatus string `json:"investigation_status,omitempty"`
		DisputedAmount      struct {
			Amount   int    `json:"amount"`
			Currency string `json:"currency"`
		} `json:"disputed_amount,omitempty"`
	}
)

func NewDisputeSrv(client Client) DisputeSrv {
	return &disputeSrv{client: client}
}

func (r *disputeSrv) GetLast250Disputes() ([]Dispute, error) {
	resp, err := r.client.Get(fmt.Sprintf("%s?limit=250&sort_by=-opened_at", disputeApiURL))
	if err != nil {
		return nil, err
	}

	var data DisputeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Disputes, nil
}
