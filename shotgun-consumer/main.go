package shotgun_consumer

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ShotgunAuth struct {
	ExpiresIn		int		`json:"expires_in"`
	RefreshToken	string	`json:"refresh_token"`
	AccessToken		string	`json:"access_token"`
	TokenType		string	`json:"token_type"`
}

type ShotgunErrorResponse struct {
	Errors []ShotgunErrorBody `json:"errors"`
}

type ShotgunErrorBody struct {
	ID string `json:"id"`
	Status int `json:"status"`
	ErrorCode int `json:"code"`
	Title string `json:"title"`
	Detail string `json:"detail"`
	Source map[string]interface{} `json:"source"`
	Meta string `json:"meta"`
}


func AuthenticateAsScript(shotgunURL, scriptName, scriptKey string) *ShotgunAuth {
	data := url.Values{}
	data.Set("client_id", scriptName)
	data.Set("client_secret", scriptKey)
	data.Set("grant_type", "client_credentials")
	requestBody := strings.NewReader(data.Encode())

	authURL := shotgunURL + "/auth/access_token"
	req, _ := http.NewRequest("POST", authURL, requestBody)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("failed to authenticate user, unhandled exception")
		return nil
	}

	if resp.StatusCode > 400 {
		logrus.Error("failed to authenticate user")
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		errResp := ShotgunErrorResponse{}
		err = json.Unmarshal(bodyBytes, &errResp)
		if err != nil {
			logrus.Error("failed to unmarshal error response from Shotgun")
		}
		for _, v := range errResp.Errors {
			logrus.WithFields(logrus.Fields{
				"statusCode": v.Status,
				"title": v.Title,
				"description": v.Detail,
			}).Error("error info from Shotgun")
		}
		return nil
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	a := ShotgunAuth{}
	err = json.Unmarshal(bodyBytes, &a)
	if err != nil {
		logrus.WithError(err).Error("failed to unmarshal auth response body")
		return nil
	}

	return &a
}