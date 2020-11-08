package shotgun_consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
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

type PaginationParameter struct {
	Size int `json:"size,omitempty"`
	Number int `json:"number,omitempty"`
}

type SearchRequest struct {
	Filters [][]string            	`json:"filters"`
	Fields  []string              	`json:"fields"`
	Sort    string              	`json:"sort,omitempty"`
	Page  	*PaginationParameter 	`json:"page,omitempty"`
}

type SearchResponse struct {
	Data []Record `json:"data"`
	PageLinks struct {
		CurrentPage string `json:"self"`
		NextPage string `json:"next"`
		PreviousPage string `json:"prev"`
	} `json:"links"`
}

type Record struct {
	ID int `json:"id"`
}

type EventLogEntry struct {
	ID int `json:"id"`
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


func GetLatestEventLogEntry(shotgunURL string, authToken *ShotgunAuth) (*EventLogEntry, error) {
	searchURL := shotgunURL + "/entity/event_log_entry/_search"
	reqBody := SearchRequest{
		Filters: [][]string {},
		Fields:[]string {
			"id",
		},
		Sort: "-created_at",
		Page: &PaginationParameter {
			Size: 10,
			Number: 1,
		},
	}
	jsonBody, _ := json.Marshal(reqBody)
	bodyBytes := bytes.NewBuffer(jsonBody)
	req, _ := http.NewRequest("POST", searchURL, bodyBytes)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/vnd+shotgun.api3_array+json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", authToken.TokenType, authToken.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("failed to query Shotgun EventLogEntries")
		return nil, err
	}
	respBytes, _ := ioutil.ReadAll(resp.Body)
	searchResp := SearchResponse{}
	err = json.Unmarshal(respBytes, &searchResp)
	if err != nil {
		logrus.Error("failed to unmarshal Shotgun response")
		return nil, err
	}

	return &EventLogEntry{ searchResp.Data[0].ID }, nil
}

func GetNewEvents(shotgunURL string, authToken *ShotgunAuth, lastEventID int) (<-chan *EventLogEntry, error){
	searchURL := shotgunURL + "/entity/event_log_entry/_search"
	reqBody := SearchRequest{
		Filters: [][]string {
			{
				"id", "greater_than", strconv.Itoa(lastEventID),
			},
		},
		Fields:[]string {
			"id",
		},
		Sort: "id",
		Page: &PaginationParameter {
			Size: 10,
			Number: 1,
		},
	}
	jsonBody, _ := json.Marshal(reqBody)
	bodyBytes := bytes.NewBuffer(jsonBody)
	req, _ := http.NewRequest("POST", searchURL, bodyBytes)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/vnd+shotgun.api3_array+json")
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", authToken.TokenType, authToken.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("failed to query Shotgun EventLogEntries")
		return nil, err
	}
	respBytes, _ := ioutil.ReadAll(resp.Body)
	searchResp := SearchResponse{}
	err = json.Unmarshal(respBytes, &searchResp)
	if err != nil {
		logrus.Error("failed to unmarshal Shotgun response")
		return nil, err
	}

	chnl := make(chan *EventLogEntry)
	go func() {
		for _, entry := range searchResp.Data {
			chnl <- &EventLogEntry{ entry.ID }
		}
		close(chnl)
	}()

	return chnl, nil
}