package goklarna

import (
	"encoding/json"
	"fmt"
)

const (
	paymentSessionApiURL       = "/payments/v1/sessions"
	paymentAuthorizationApiURL = "/payments/v1/authorizations"
)

type (
	// PaymentSrv type describe the payment api client methods
	PaymentSrv interface {
		CreateNewSession(*PaymentOrder) (*PaymentSession, error)
		UpdateExistingSession(string, *PaymentOrder) error
		CreateNewOrder(string, *PaymentOrder) (*PaymentOrderInfo, error)
		CancelExistingAuthorization(string) error
		GetCustomerToken(authorizationToken string, r *CustomerTokenRequest) (*CustomerTokenResponse, error)
	}

	paymentSrv struct {
		client Client
	}

	// PaymentOrderInfo type is the response coming back from creating an order in the Payment API
	PaymentOrderInfo struct {
		OrderID                 string      `json:"order_id,omitempty"`
		RedirectURL             string      `json:"redirect_url,omitempty"`
		FraudStatus             string      `json:"fraud_status,omitempty"`
		AuthorizedPaymentMethod interface{} `json:"authorized_payment_method"`
	}

	// SessionResponse type encapsulate the two fields that the API response with when creating a new session
	PaymentSession struct {
		// SessionID Id of the created session
		SessionID string `json:"session_id"`
		// ClientToken Token to be passed to the JS client
		ClientToken string `json:"client_token"`
	}

	// PaymentOrder type is the request payload to create an order from the Payment API by providing the order
	// structure and the authorization token
	PaymentOrder struct {
		Design                 string               `json:"design,omitempty"`
		PurchaseCountry        string               `json:"purchase_country"`
		PurchaseCurrency       string               `json:"purchase_currency"`
		Locale                 string               `json:"locale"`
		BillingAddress         *Address             `json:"billing_address"`
		ShippingAddress        *Address             `json:"shipping_address,omitempty"`
		OrderAmount            int                  `json:"order_amount"`
		OrderTaxAmount         int                  `json:"order_tax_amount"`
		OrderLines             []*OrderLine         `json:"order_lines"`
		Customer               *CustomerInfo        `json:"customer,omitempty"`
		MerchantData           *string              `json:"merchant_data,omitempty"`
		MerchantURLS           *PaymentMerchantURLS `json:"merchant_urls,omitempty"`
		MerchantReference1     string               `json:"merchant_reference1,omitempty"`
		MerchantReference2     string               `json:"merchant_reference2,omitempty"`
		Options                *PaymentOptions      `json:"options,omitempty"`
		Attachment             *Attachment          `json:"attachment,omitempty"`
		AutoCapture            *bool                `json:"auto_capture"`
		PaymentMethodReference *string              `json:"payment_method_reference"`
		Intent                 *string              `json:"intent"`
	}

	// PaymentOptions type Options for this purchase
	PaymentOptions struct {
		ColorButton              string `json:"color_button,omitempty"`
		ColorButtonText          string `json:"color_button_text,omitempty"`
		ColorCheckbox            string `json:"color_checkbox,omitempty"`
		ColorCheckboxCheckmark   string `json:"color_checkbox_checkmark,omitempty"`
		ColorHeader              string `json:"color_header,omitempty"`
		ColorLink                string `json:"color_link,omitempty"`
		ColorBorder              string `json:"color_border,omitempty"`
		ColorBorderSelected      string `json:"color_border_selected,omitempty"`
		ColorText                string `json:"color_text,omitempty"`
		ColorDetails             string `json:"color_details,omitempty"`
		ColorTextSecondary       string `json:"color_text_secondary,omitempty"`
		RadiusBorder             string `json:"radius_border,omitempty"`
		AllowPartialAddress      *bool  `json:"allow_partial_address,omitempty"`
		DisableClientSideUpdates string `json:"disable_client_side_updates,omitempty"`
	}

	// The merchant urls structure
	PaymentMerchantURLS struct {
		Confirmation string `json:"confirmation"`
		Notification string `json:"notification,omitempty"`
	}

	// CustomerInfo type is Information about the liable customer of the order
	CustomerInfo struct {
		DateOfBirth                  string `json:"date_of_birth,omitempty"`                  // yyyy-mm-dd
		Gender                       string `json:"gender,omitempty"`                         // 'male' or 'female'
		LastFourSSN                  string `json:"last_four_ssn,omitempty"`                  // for US customers
		NationalIdentificationNumber string `json:"national_identification_number,omitempty"` // for EU customers
		KlarnaAccessToken            string `json:"klarna_access_token,omitempty"`
	}

	CustomerTokenRequest struct {
		BillingAddress   *Address      `json:"billing_address"`
		Customer         *CustomerInfo `json:"customer,omitempty"`
		Description      string        `json:"description"`
		IntendedUse      string        `json:"intended_use"`
		PurchaseCountry  string        `json:"purchase_country"`
		PurchaseCurrency string        `json:"purchase_currency"`
		Locale           string        `json:"locale"`
		MerchantUrls     struct {
			Confirmation string `json:"confirmation"`
		} `json:"merchant_urls"`
	}

	CustomerTokenResponse struct {
		BillingAddress         *Address      `json:"billing_address"`
		Customer               *CustomerInfo `json:"customer,omitempty"`
		PaymentMethodReference *string       `json:"payment_method_reference"`
		RedirectUrl            string        `json:"redirect_url"`
		TokenId                string        `json:"token_id"`
	}
)

const (
	PurchaseCountrySE = "SE"

	PurchaseCurrencySEK = "SEK"

	LocaleSweden = "sv-SE"
)

// CreateNewSession method calls payment session api and return an error if there is any, PaymentSession struct
// is returned on success
func (srv *paymentSrv) CreateNewSession(po *PaymentOrder) (*PaymentSession, error) {
	res, err := srv.client.Post(paymentSessionApiURL, po)
	if nil != err {
		return nil, err
	}
	ps := new(PaymentSession)
	err = json.NewDecoder(res.Body).Decode(ps)

	return ps, err
}

// GetCustomerToken makes query to generate customer token
func (srv *paymentSrv) GetCustomerToken(authorizationToken string, r *CustomerTokenRequest) (*CustomerTokenResponse, error) {
	uri := fmt.Sprintf("%s/%s/customer-token", paymentAuthorizationApiURL, authorizationToken)
	res, err := srv.client.Post(uri, r)
	if err != nil {
		return nil, err
	}

	tokenResp := new(CustomerTokenResponse)
	err = json.NewDecoder(res.Body).Decode(tokenResp)
	return tokenResp, err
}

// UpdateExistingSession method calls update payment session api and return an error if there is any
func (srv *paymentSrv) UpdateExistingSession(id string, po *PaymentOrder) error {
	uri := fmt.Sprintf("%s/%s", paymentSessionApiURL, id)
	_, err := srv.client.Post(uri, po)

	return err
}

// CreateNewOrder method creates a new payment order with the given token and order
func (srv *paymentSrv) CreateNewOrder(token string, po *PaymentOrder) (*PaymentOrderInfo, error) {
	path := fmt.Sprintf("%s/%s/order", paymentAuthorizationApiURL, token)
	res, err := srv.client.Post(path, po)
	if nil != err {
		return nil, err
	}

	pof := new(PaymentOrderInfo)
	err = json.NewDecoder(res.Body).Decode(pof)

	return pof, err
}

// CancelExistingAuthorization method calls the API end-point
func (srv *paymentSrv) CancelExistingAuthorization(token string) error {
	path := fmt.Sprintf("%s/%s", paymentAuthorizationApiURL, token)
	_, err := srv.client.Delete(path)

	return err
}

// NewPaymentSrv Return a new payment instance while providing
func NewPaymentSrv(c Client) PaymentSrv {
	return &paymentSrv{c}
}
