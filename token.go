package goklarna

import (
	"encoding/json"
	"fmt"
)

const (
	customerTokenApiURL = "/customer-token/v1/tokens"
)

type (
	// TokenSrv type describe the token api client methods
	TokenSrv interface {
		CreateNewOrder(token string, po *PaymentOrder) (*PaymentOrderInfo, error)
	}

	tokenSrv struct {
		client Client
	}
)

// CreateNewOrder method creates a new payment order with the given token and order
func (srv *tokenSrv) CreateNewOrder(token string, po *PaymentOrder) (*PaymentOrderInfo, error) {
	path := fmt.Sprintf("%s/%s/order", customerTokenApiURL, token)
	res, err := srv.client.Post(path, po)
	if nil != err {
		return nil, err
	}

	pof := new(PaymentOrderInfo)
	err = json.NewDecoder(res.Body).Decode(pof)

	return pof, err
}

// NewTokenSrv Return a new token instance while providing
func NewTokenSrv(c Client) TokenSrv {
	return &tokenSrv{c}
}
