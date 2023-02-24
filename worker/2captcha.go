package worker

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TwoCaptcha struct {
	Host string
	Key  string
}

type twoCaptchaInResponse struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
	Error   string `json:"error_text"`
}

type twoCaptchaResResponse struct {
	Status  int    `json:"status"`
	Request string `json:"request"`
}

func (tc *TwoCaptcha) recaptcha(pageUrl string, siteKey string, dataS string) (gRecaptchaResponse string, err error) {
	apiUrl := fmt.Sprintf(
		"http://%s/in.php?key=%s&method=userrecaptcha&googlekey=%s&data-s=%s&pageurl=%s&json=1",
		tc.Host,
		url.QueryEscape(tc.Key),
		url.QueryEscape(siteKey),
		url.QueryEscape(dataS),
		url.QueryEscape(pageUrl),
	)

	resp, err := http.Get(apiUrl)
	if err != nil {
		err = errors.New("twoCaptcha: fetch failed " + err.Error())
		return
	}

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("twoCaptcha: failed to read response body " + err.Error())
		return
	}

	var inResponse twoCaptchaInResponse
	json.Unmarshal(respBodyBytes, &inResponse)

	if inResponse.Status != 1 || len(inResponse.Request) == 0 {
		err = errors.New("twoCaptcha: invalide inResponse " + inResponse.Error)
		return
	}

	for i := 0; i < 15; i++ {
		time.Sleep(20 * time.Second)

		apiUrl = fmt.Sprintf(
			"http://%s/res.php?key=%s&action=get&id=%s&json=1",
			tc.Host,
			url.QueryEscape(tc.Key),
			url.QueryEscape(inResponse.Request),
		)

		resp, err = http.Get(apiUrl)
		if err != nil {
			err = errors.New("twoCaptcha: fetch failed " + err.Error())
			return
		}

		respBodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			err = errors.New("twoCaptcha: failed to read response body " + err.Error())
			return
		}

		var resResponse twoCaptchaResResponse
		json.Unmarshal(respBodyBytes, &resResponse)

		if resResponse.Status > 0 {
			gRecaptchaResponse = resResponse.Request
			return
		} else if strings.HasPrefix(resResponse.Request, "ERROR") {
			err = errors.New("twoCaptcha: " + resResponse.Request)
			return
		}
	}

	err = errors.New("twoCaptcha: solver timeout")

	return
}
