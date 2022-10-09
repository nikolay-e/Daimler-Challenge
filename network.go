package main

import (
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type ApiKeyData struct {
	Key                         string `json:"Key"`
	KeyAauthVersion             string `json:"Key AauthVersion"`
	ReadAccessToken             string `json:"Read Access Token"`
	ReadAccessTokenAauthVersion string `json:"Read Access Token AauthVersion"`
}

type Network struct {
	ApiKey  ApiKeyData
	baseUrl string
}

const (
	BASE_URL = "https://api.themoviedb.org/3/"
)

func (n *Network) MakeSeriesSearchUrl(page int, query string) string {
	request := fmt.Sprintf(BASE_URL+"search/multi?api_key=%s&page=%d&include_adult=false&query=%s", n.ApiKey.Key, page, html.EscapeString(query))
	res := strings.ReplaceAll(request, " ", "+")
	return res
}

func (n *Network) MakeSerieUrl(tvid int) string {
	request := fmt.Sprintf(BASE_URL+"tv/%d?api_key=%s", tvid, n.ApiKey.Key)
	return request
}

func (n *Network) MakeSeasonUrl(tvid int, season int) string {
	request := fmt.Sprintf(BASE_URL+"tv/%d/season/%d?api_key=%s", tvid, season, n.ApiKey.Key)
	return request
}

func (n *Network) MakeEpisodeUrl(tvid int, season int, episode int) string {
	request := fmt.Sprintf(BASE_URL+"tv/%d/season/%d/episode/%d?api_key=%s", tvid, season, episode, n.ApiKey.Key)
	return request
}

func (n *Network) GetRequestWith(url string) ([]byte, error) {

	response, responseError := http.Get(url)
	if responseError != nil {
		log.Printf("GetRequestWith: http.Get error")
		return []byte{}, responseError
	}

	responseData, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		log.Printf("GetRequestWith: ioutil.ReadAll error", readError)
		return []byte{}, readError
	}

	return responseData, nil
}
