/*
 * Copyright (C) 2013, Jordon Smith <jrd@mockra.net>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright
 *    notice, this list of conditions and the following disclaimer in the
 *    documentation and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY AUTHOR AND CONTRIBUTORS ``AS IS'' AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED.  IN NO EVENT SHALL AUTHOR OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS
 * OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
 * LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY
 * OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF
 * SUCH DAMAGE.
 */
package tesgo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	//API_URL = "https://secure.techfortesco.com/groceryapi_b1/restservice.aspx" // stable beta1
	//API_URL = "https://secure.techfortesco.com/groceryapi/restservice.aspx" // stable nightly
	API_URL = "https://secure.techfortesco.com/groceryapi_ops/restservice.aspx" // unstable ops-test
)

type Session struct {
	devKey, appKey string
	sessionKey     string
}

type LoginResponse struct {
	StatusCode             int
	StatusInfo             string
	BranchNumber           string
	CustomerId             string
	CustomerName           string
	SessionKey             string
	InAmmendOrderMode      string
	ChosenDeliverySlotInfo string
	CustomerForename       string
}

type SearchResponse struct {
	StatusCode        int
	StatusInfo        string
	PageNumber        int
	TotalPageCount    int
	TotalProductCount int
	PageProductCount  int
	Products          []Product
}

type ChangeBasketResponse struct {
	StatusCode int
	StatusInfo string
}

type ListBasketResponse struct {
	StatusCode                 int
	StatusInfo                 string
	BasketID                   string
	InAmendOrderMode           string
	BasketGuideMultiBuySavings string
	BasketGuidePrice           string
	BasketQuantity             string
	BasketTotalClubcardPoints  string
	BasketLines                []BasketLine
}

type Product struct {
	BaseProductId                 string
	EANBarcode                    string
	CheaperAlternativeProductId   string
	CookingAndUsage               string
	ExtendedDescription           string
	HealthierAlternativeProductId string
	ImagePath                     string
	MaximumPurchaseQuantity       int
	Name                          string
	OfferPromotion                string
	OfferValidity                 string
	OfferLabelImagePath           string
	Price                         float32
	PriceDescription              string
	ProductId                     string
	ProductType                   string
	Rating                        int
	StorageInfo                   string
	UnitPrice                     float32
	UnitType                      string
	RDA_Calories_Count            string
	RDA_Calories_Percent          string
	RDA_Sugar_Grammes             string
	RDA_Sugar_Percent             string
	RDA_Fat_Grammes               string
	RDA_Fat_Percent               string
	RDA_Saturates_Grammes         string
	RDA_Saturates_Percent         string
	RDA_Salt_Grammes              string
	RDA_Salt_Percent              string
	NutrientsCount                int
	Nutrients                     []Nutrient
	IngredientsCount              int
	Ingredients                   []Ingredient
}

type BasketLine struct {
	BasketLineErrorMessage  string
	BasketLineGuidePrice    string
	BasketLinePromoMessage  string
	BasketLineQuantity      string
	BaseProductId           string
	EANBarcode              string
	ImagePath               string
	MaximumPurchaseQuantity int
	Name                    string
	OfferPromotion          string
	OfferValidity           string
	Price                   float32
	PriceDescription        string
	ProductId               string
	ProductType             string
	StorageInfo             string
	UnitPrice               float32
	UnitType                string
	NoteForPersonalShopper  string
	SubstitutionNote        string
}

type Nutrient struct {
	NutrientName       string
	SampleDescription  string
	SampleSize         string
	ServingDescription string
	ServingSize        string
}

type Ingredient struct {
	Name string
}

func tescoError(code int, info string) error {
	return errors.New(strconv.Itoa(code) + ": " + info)
}

// Returns a new session for using the Tesco API.
// You are required to Login() with a Tesco account to
// obtain a valid session key before calling API functions.
//
// Sign up for account here: https://secure.techfortesco.com/tescoapiweb/secretlogin.aspx
func New(devKey, appKey string) *Session {
	return &Session{devKey: devKey, appKey: appKey}
}

// Login with Tesco account to gain a valid session key, required to call
// the other API functions.
//
// You can sign up for an account here: https://secure.tesco.com/register/
func (s *Session) Login(email, password string) error {
	u, _ := url.Parse(API_URL)
	v := url.Values{}
	v.Add("COMMAND", "LOGIN")
	v.Add("DEVELOPERKEY", s.devKey)
	v.Add("APPLICATIONKEY", s.appKey)
	v.Add("EMAIL", email)
	v.Add("PASSWORD", password)
	u.RawQuery = v.Encode()

	b, err := httpGetBody(u.String())
	if err != nil {
		return err
	}

	// Error 150 is returned as a string, reported
	b = []byte(strings.Replace(string(b), "\"StatusCode\": \"150\"", "\"StatusCode\": 150", 1))

	var resp LoginResponse
	if err = json.Unmarshal(b, &resp); err != nil {
		return err
	}

	if resp.StatusCode != 0 {
		return tescoError(resp.StatusCode, resp.StatusInfo)
	}
	s.sessionKey = resp.SessionKey
	return nil
}

// Searches for products using text or 13-digit numeric barcode.
//
// Warning: Using extended search on a generic search will take a very long time, avoid using
// for anything but ProductID searches or barcode searches.
func (s *Session) ProductSearch(search string, page int, extended bool) (*SearchResponse, error) {
	if s.sessionKey == "" {
		return nil, errors.New("No session key, you must Login()")
	}

	u, _ := url.Parse(API_URL)
	v := url.Values{}
	v.Add("COMMAND", "PRODUCTSEARCH")
	v.Add("SESSIONKEY", s.sessionKey)
	v.Add("SEARCHTEXT", search)
	v.Add("PAGE", strconv.Itoa(page))
	v.Add("EXTENDEDINFO", yesNo(extended))
	u.RawQuery = v.Encode()

	b, err := httpGetBody(u.String())
	if err != nil {
		return nil, err
	}

	if extended { // commas missing in extended arrays (reported)
		b = []byte(strings.Replace(string(b), "}\r\n{", "},\r\n{", -1))
	}

	var resp SearchResponse
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, errors.New(err.Error() + "\n" + string(b))
	}

	if resp.StatusCode != 0 {
		return nil, errors.New(strconv.Itoa(resp.StatusCode) + ": " + resp.StatusInfo)
	}
	return &resp, nil
}

// Adds products to the basket, removes them from the basket, and updates the basket.
//
// quantity:
// A positive or negative value that changes the products in the
// basket by that quantity, according to these rules:
// 1)
// If the product was absent from the basket before that
// product was added, it is inserted into the basket at the
// requested quantity.
// 2)
// If the product was already in the basket, the quantity is
// increased by requested quantity if positive, or reduced by
// the requested quantity if the requested quantity is negative.
// 3)
// If a negative requested quantity is equal to or larger than
// the existing quantity, the product is removed from the
// basket.
// 4)
// For products that sell by weight, quantities added or
// removed are still “each”. For example, if you are adding
// apples that are priced per Kg, selecting “2” for this
// parameter will add 2 individual apples to the basket, not 2
// Kg of apples.
func (s *Session) ChangeBasket(productId string, quantity int, substitute bool) (*ChangeBasketResponse, error) {
	if s.sessionKey == "" {
		return nil, errors.New("No session key, you must Login()")
	}

	u, _ := url.Parse(API_URL)
	v := url.Values{}
	v.Add("COMMAND", "CHANGEBASKET")
	v.Add("SESSIONKEY", s.sessionKey)
	v.Add("PRODUCTID", productId)
	v.Add("CHANGEQUANTITY", strconv.Itoa(quantity))
	v.Add("SUBSTITUTION", yesNo(substitute))
	u.RawQuery = v.Encode()

	b, err := httpGetBody(u.String())
	if err != nil {
		return nil, err
	}

	var resp ChangeBasketResponse
	if err = json.Unmarshal(b, &resp); err != nil {
		return nil, errors.New(err.Error() + "\n" + string(b))
	}

	if resp.StatusCode != 0 {
		return nil, tescoError(resp.StatusCode, resp.StatusInfo)
	}
	return &resp, nil
}

// Lists the contents of the basket.
//
// fast:
// massively speeds up retrieval of the basket at the
// cost of not being able to find all of the core attributes required
// for a product, such as EANBarcode. Use FAST=Y to get the
// API to abandon further searching to retrieve all the core
// attributes when retrieving the basket. Test this mode to see if
// your application can cope without the missing attributes.
func (s *Session) ListBasket(fast bool) (*ListBasketResponse, error) {
	u, _ := url.Parse(API_URL)
	v := url.Values{}
	v.Add("COMMAND", "LISTBASKET")
	v.Add("SESSIONKEY", s.sessionKey)
	v.Add("FAST", yesNo(fast))
	u.RawQuery = v.Encode()

	b, err := httpGetBody(u.String())
	if err != nil {
		return nil, err
	}

	var resp ListBasketResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, errors.New(err.Error() + "\n" + string(b))
	}

	if resp.StatusCode != 0 {
		return nil, tescoError(resp.StatusCode, resp.StatusInfo)
	}
	return &resp, nil
}

func httpGetBody(url string) ([]byte, error) {
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func yesNo(cond bool) string {
	// ternary, where art thou?
	if cond {
		return "Y"
	}
	return "N"
}
