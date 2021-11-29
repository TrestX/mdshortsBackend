package api

import (
	"MdShorts/pkg/entity"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/aekam27/trestCommon"
	"github.com/spf13/viper"
)

func GetHealthTopHeadlines(country, language, category string) (entity.TopNewsStruct, error) {
	queryString := "country=" + country + "&pageSize=100&category=" + category + "&language=" + language + "&apiKey=" + viper.GetString("newsapi.key")
	if country == "" {
		queryString = "language=" + language + "&pageSize=100&category=" + category + "&apiKey=" + viper.GetString("newsapi.key")
	}
	url := "https://newsapi.org/v2/top-headlines?" + queryString
	body, err := trestCommon.GetApi(" ", url)
	if err != nil {
		return entity.TopNewsStruct{}, err
	}
	var resp entity.TopNewsStruct
	err = json.Unmarshal(body, &resp)
	return resp, err
}
func GetNewslines(search, time, page string) (entity.TopNewsStruct, error) {
	url := "https://newsapi.org/v2/everything?q=" + search + "&pageSize=100&page=" + page + "&language=en&sortBy=publishedAt&apiKey=" + viper.GetString("newsapi.key")
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return entity.TopNewsStruct{}, err
	}
	res, err := client.Do(req)
	if err != nil {
		return entity.TopNewsStruct{}, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return entity.TopNewsStruct{}, err
	}
	var resp entity.TopNewsStruct
	err = json.Unmarshal(body, &resp)
	return resp, err
}

type CSS struct {
	To     string `json:"to"`
	Body   string `json:"body"`
	Source string `json:"source"`
}
type MSG struct {
	Messages []CSS `json:"messages"`
}
type CSSR struct {
	ResponseCode string `json:"response_code"`
}

func ClickSend(auth string, number string, otp int) (string, error) {
	strOtp := strconv.Itoa(otp)
	var cssl []CSS
	var css CSS
	var msg MSG
	css.To = number
	css.Body = strOtp + " " + viper.GetString("clicksend.mbody")
	css.Source = "php"
	cssl = append(cssl, css)
	msg.Messages = cssl
	url := viper.GetString("clicksend.postApi")
	method := "POST"
	basicAuth := "Basic " + auth
	client := &http.Client{}
	requestByte, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	requestReader := bytes.NewReader(requestByte)
	req, err := http.NewRequest(method, url, requestReader)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", basicAuth)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var resp CSSR
	err = json.Unmarshal(body, &resp)
	return resp.ResponseCode, err
}
