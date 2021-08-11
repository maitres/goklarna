package main

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/maitres/goklarna"
)

func main() {

	u, _ := url.Parse(goklarna.BaseUrlEuroPlayground)

	c := goklarna.NewClient(goklarna.Config{
		BaseURL:     u,
		APIUsername: os.Getenv("API_USERNAME"),
		APIPassword: os.Getenv("API_PASSWORD"),
		Timeout:     10 * time.Second,
	})

	createNewSession(c)
	//generateCustomerToken(c)
	//placeOrderByCustomerToken(c)
}

func createNewSession(client goklarna.Client) {
	pmsrv := goklarna.NewPaymentSrv(client)

	s, err := pmsrv.CreateNewSession(&goklarna.PaymentOrder{
		PurchaseCountry:  goklarna.PurchaseCountrySE,
		PurchaseCurrency: goklarna.PurchaseCurrencySEK,
		Locale:           goklarna.LocaleSweden,
		OrderAmount:      1000,
		OrderLines: []*goklarna.OrderLine{
			{
				Type:        goklarna.DigitalLineType,
				Name:        "Authorization",
				Quantity:    1,
				UnitPrice:   1000,
				TotalAmount: 1000,
			},
		},
		Customer: &goklarna.CustomerInfo{
			DateOfBirth:                  "1985-09-06",
			NationalIdentificationNumber: "850906-4583",
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(s.ClientToken)
	fmt.Println(s.SessionID)
}

func generateCustomerToken(client goklarna.Client) {
	pmsrv := goklarna.NewPaymentSrv(client)

	var (
		authToken = "7da684f9-645e-1016-836f-50048840f740"
		request   = goklarna.CustomerTokenRequest{
			BillingAddress: nil,
			Customer: &goklarna.CustomerInfo{
				DateOfBirth:                  "1985-09-06",
				NationalIdentificationNumber: "850906-4583",
			},
			Description:      "some text",
			IntendedUse:      "SUBSCRIPTION",
			PurchaseCountry:  goklarna.PurchaseCountrySE,
			PurchaseCurrency: goklarna.PurchaseCurrencySEK,
			Locale:           goklarna.LocaleSweden,
		}
	)

	resp, err := pmsrv.GetCustomerToken(authToken, &request)
	fmt.Println(err)
	fmt.Println(resp)
}

func placeOrderByCustomerToken(client goklarna.Client) {
	toksrv := goklarna.NewTokenSrv(client)

	order, err := toksrv.CreateNewOrder("b3acb3a2-79c7-4ad9-880d-abc06e3968f0", &goklarna.PaymentOrder{
		PurchaseCountry:  goklarna.PurchaseCountrySE,
		PurchaseCurrency: goklarna.PurchaseCurrencySEK,
		Locale:           goklarna.LocaleSweden,
		OrderAmount:      10000,
		OrderLines: []*goklarna.OrderLine{
			{
				Type:        goklarna.DigitalLineType,
				Name:        "Authorization",
				Quantity:    1,
				UnitPrice:   10000,
				TotalAmount: 10000,
			},
		},
		Customer: &goklarna.CustomerInfo{
			DateOfBirth:                  "1985-09-06",
			NationalIdentificationNumber: "850906-4583",
		},
		AutoCapture:        goklarna.Bool(true),
		MerchantReference1: "124124",
	})

	fmt.Println(err)
	fmt.Println(order)
}
