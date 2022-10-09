package main

import (
	"encoding/json"
	"log"
)

type SerieWithTitle struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	OriginalTitle string `json:"original_title"`
	MediaType     string `json:"media_type"`
}

type SerieWithName struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	MediaType    string `json:"media_type"`
}

type Season struct {
	Name string `json:"name"`
}

type SerieInfo struct {
	NumberOfSeasons int      `json:"number_of_seasons"`
	Seasons         []Season `json:"seasons"`
}

type EpisodeInfo struct {
	EpisodeNumber int     `json:"episode_number"`
	Name          string  `json:"name"`
	Overview      string  `json:"overview"`
	VoteCount     int     `json:"vote_count"`
	VoteAverage   float32 `json:"vote_average"`
	StillPath     string  `json:"still_path"`
	AirDate       string  `json:"air_date"`
}

type SeasonInfo struct {
	Episodes []EpisodeInfo `json:"episodes"`
}

type SerieSearchResponseWhatLeft struct {
	Results []SerieWithName
}

type SerieSearchResponseRaw struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`

	Results  []SerieWithTitle `json:"results"`
	WhatLeft json.RawMessage
}

type SerieSearchInfo struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`

	Results []SerieWithTitle
}

func MyUnmarshalForSeries(data []byte, msr *SerieSearchInfo) SerieSearchInfo {

	unmarshallErr1 := json.Unmarshal(data, msr)
	if unmarshallErr1 != nil {
		log.Fatalln("MyUnmarshalForSeries: json.Unmarshal Failed with", unmarshallErr1)
		return SerieSearchInfo{}
	}

	var whatLeft SerieSearchResponseWhatLeft

	unmarshallErr2 := json.Unmarshal(data, &whatLeft)
	if unmarshallErr2 != nil {
		log.Fatalln("MyUnmarshalForSeries: json.Unmarshal Failed with", unmarshallErr2)
		return SerieSearchInfo{}
	}

	for _, v := range whatLeft.Results {
		msr.Results = append(msr.Results, SerieWithTitle{v.Id, v.Name, v.OriginalName, v.MediaType})
	}

	var validResult SerieSearchInfo
	for _, v := range msr.Results {

		if v.MediaType != "tv" {
			continue
		}

		if v.Title != "" || v.OriginalTitle != "" {
			validResult.Results = append(validResult.Results, v)
		} else {
			continue
		}
	}

	validResult.Page = msr.Page
	validResult.TotalResults = msr.TotalResults
	validResult.TotalPages = msr.TotalPages

	return validResult
}
