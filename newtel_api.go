package newtel_api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type NewTelConfig struct {
	ApiKey  string `key:"api_key" validate:"nonzero"`
	ApiSign string `key:"api_sign" validate:"nonzero"`
}

type NewTel struct {
	Cfg        NewTelConfig
	BaseUrl    string
	HTTPClient *http.Client
}

const (
	baseURL = "https://api.new-tel.net/"
)

func NewTelClient(cfg NewTelConfig) *NewTel {
	return &NewTel{
		Cfg:     cfg,
		BaseUrl: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

type CallPasswordData struct {
	DstNumber string `json:"dstNumber"`
	Async     string `json:"async"`
	Pin       string `json:"pin"`
	Timeout   string `json:"timeout"`
}

type CallPasswordResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result      string `json:"result"`
		CallDetails struct {
			CallId      string `json:"callId"`
			Code        string `json:"pin"`
			Status      string `json:"status"`
			Oper        string `json:"oper"`
			Region      string `json:"region"`
			ReasonCode  string `json:"reasonCode"`
			PhoneNumber string `json:"phoneNumber"`
		} `json:"callDetails"`
	} `json:"data"`
}

func (n *NewTel) CallPassword(to, code string) (*CallPasswordResponse, error) {
	if to[0] == '+' {
		to = to[1:]
	} //TODO: add other validations

	if len(code) != 4 {
		return nil, fmt.Errorf("code should be a number with len equal to 4 : %v", code)
	} //TODO: check on numbers

	data := CallPasswordData{
		DstNumber: to,
		Async:     "0",
		Pin:       "1" + code,
		Timeout:   "30",
	}
	resp, err := n.MakeRequest("call-password/start-password-call", data)
	if err != nil {
		return nil, err
	}
	println(string(resp))
	callResponse := new(CallPasswordResponse)
	err = json.Unmarshal(resp, callResponse)
	return callResponse, err
}

func getKey(method string, apiKey string, data []byte, apiSign string) string {
	timeU := strconv.Itoa(int(time.Now().Unix()))
	hash := sha256.Sum256([]byte(method + "\n" + timeU + "\n" + apiKey + "\n" + string(data) + "\n" + apiSign))
	return apiKey + timeU + fmt.Sprintf("%x", hash)
}

func (n *NewTel) MakeRequest(method string, data interface{}) ([]byte, error) {
	binaryData, _ := json.Marshal(data)

	req, err := http.NewRequest("POST", n.BaseUrl+method, bytes.NewBuffer(binaryData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+getKey(method, n.Cfg.ApiKey, binaryData, n.Cfg.ApiSign))

	res, err := n.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("req failed: %v", err)
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read resp failed: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		/*exception := new(ErrorResponse)
		err = json.Unmarshal(responseBody, exception)
		if err != nil {
			return fmt.Errorf("unmarshal err failed: %v", err)
		}*/
		return nil, fmt.Errorf("call not ok: %s", responseBody)
	}

	return responseBody, err
}

type ErrorResponse struct {
	Status string `json:"status"`
	Msg    string `json:"message"`
}
