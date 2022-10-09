package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/therecipe/qt/widgets"
)

var inputChan = make(chan string, 100)
var inputToServe = make(chan string, 1)

var loadNewPage = make(chan bool, 100)

type Challenge struct {
	network Network
	gui     *Gui

	input          string
	tvid           int
	season_number  int
	episode_number int

	page                       int
	totalPages                 int
	disableScrollEventHandling bool
	noNewInput                 bool
	loadingNewPage             bool
}

func NewChallenge(apiKey ApiKeyData) *Challenge {
	c := Challenge{}
	c.gui = NewGui()
	c.page = 1
	c.totalPages = 1
	c.network = Network{ApiKey: apiKey}
	c.noNewInput = true
	c.loadingNewPage = false

	c.gui.input.ConnectTextChanged(func(text string) {
		fmt.Println("input text chaged: User input", text)
		go func(text string) {
			inputChan <- text
		}(text)

	})

	c.gui.searchButton.ConnectClicked(func(checked bool) {
		text := c.gui.input.Text()
		fmt.Println("searchButton cliecked: User input", text)
		go func(text string) {
			inputChan <- text
		}(text)
	})

	c.gui.scrollArea.VerticalScrollBar().ConnectValueChanged(func(value int) {
		if value == c.gui.scrollArea.VerticalScrollBar().Maximum() && !c.disableScrollEventHandling {
			fmt.Println("scrollArea.VerticalScrollBar reached the bottom", value)
			go func() {
				loadNewPage <- true
			}()
		}
	})

	return &c
}

func (c *Challenge) Run() {
	log.Println("Running new challenge")

	go c.inputSkipper()
	go c.serveSeriesRequests()

	c.gui.window.Show()
	widgets.QApplication_Exec()
}

func (c *Challenge) inputSkipper() {

	var nextInput string
	var prevInput string
	var pNextInput string
	var pPrevInput string

	var tmp string

	for {
		select {
		case tmp = <-inputChan:
			{
				log.Println("inputSkipper: skipping last one", prevInput)
				log.Println("inputSkipper: serving", tmp)

				c.noNewInput = false
				inputToServe <- tmp
				c.noNewInput = true

				pPrevInput = prevInput
				pNextInput = nextInput

				prevInput = nextInput
				nextInput = tmp
			}
		default:
			{
				if (nextInput == pNextInput) && (prevInput == pPrevInput) {

					log.Println("inputSkipper: no input")
					log.Println("inputSkipper: serving", tmp)

					tmp = <-inputChan

					c.noNewInput = false
					inputToServe <- tmp
					c.noNewInput = true

					pPrevInput = prevInput
					pNextInput = nextInput

					prevInput = nextInput
					nextInput = tmp

				} else {
					log.Println("inputSkipper: skipping", prevInput)

					pPrevInput = prevInput
					pNextInput = nextInput

					prevInput = pPrevInput
					nextInput = pNextInput
				}
			}
		}
	}
}

func (c *Challenge) serveSeriesRequests() {
	log.Println("serveSeriesRequests: service started")
	for {
		select {
		case c.input = <-inputToServe:
			{
				c.noNewInput = true
				c.requestSeries(false)
				break
			}
		case <-loadNewPage:
			{
				length := len(loadNewPage)
				for i := 0; i < length; i++ {

				}
				c.requestSeries(true)
				break
			}
		}
	}
}

func (c *Challenge) requestSeries(cont bool) {
	if cont {
		log.Println("requestSeries:", c.input, "continuation")
	} else {
		log.Println("requestSeries: new input", c.input)
	}

	c.disableScrollEventHandling = false

	if strings.TrimSpace(c.input) == "" && cont {
		return
	}

	if strings.TrimSpace(c.input) == "" {
		log.Println("requestSeries: empty input - won't serve")
		c.page = 1
		c.totalPages = 1
		c.gui.ResetResultGrid()
		return
	}

	if cont {
		c.gui.MaxSeriesOnGrid += NR_OF_VISIBLE_SERIES
	} else {
		c.page = 1
		c.totalPages = 1
		c.gui.ResetResultGrid()
	}

	for ; (c.gui.SeriesOnGrid <= c.gui.MaxSeriesOnGrid) && (c.page <= c.totalPages) && c.noNewInput; c.page++ {
		log.Println("requestSeries: getting page", c.page, "for input", c.input)

		responseData, err := c.network.GetRequestWith(c.network.MakeSeriesSearchUrl(c.page, c.input))
		if err != nil {
			log.Println("clickedOnSerie: network error ", err)
			continue
		}

		var serieSearchInfo SerieSearchInfo
		serieSearchResponse := MyUnmarshalForSeries(responseData, &serieSearchInfo)

		c.totalPages = serieSearchResponse.TotalPages

		var seriesNames []string
		var ids []int

		for _, v := range serieSearchResponse.Results {
			ids = append(ids, v.Id)
			seriesNames = append(seriesNames, v.Title)
		}

		c.gui.updateViewWithIdsStoring(ids, seriesNames, c.clickedOnSerie, cont)

	}

	if !c.noNewInput {
		log.Println("requestSeries: interrupted requests for", c.input, "because of the new input")
	}

	c.loadingNewPage = false

	c.gui.window.Update()

}

func (c *Challenge) clickedOnSerie(checked bool, tvid int) {
	log.Println("clickedOnSerie", tvid)

	c.disableScrollEventHandling = true
	c.tvid = tvid

	responseData, err := c.network.GetRequestWith(c.network.MakeSerieUrl(c.tvid))
	if err != nil {
		log.Println("clickedOnSerie: network error ", err)
		return
	}

	var serieInfo SerieInfo

	unmarshallErr := json.Unmarshal(responseData, &serieInfo)
	if unmarshallErr != nil {
		log.Fatalln("clickedOnSerie: json.Unmarshal Failed with", unmarshallErr)
		return
	}

	var seasons []string
	for _, v := range serieInfo.Seasons {
		seasons = append(seasons, fmt.Sprintf("%s", v.Name))
	}

	c.gui.updateView(seasons, c.clickedOnSeason)
}

func (c *Challenge) clickedOnSeason(checked bool, season_number int) {
	log.Println("clickedOnSeason", season_number)

	c.season_number = season_number
	responseData, err := c.network.GetRequestWith(c.network.MakeSeasonUrl(c.tvid, c.season_number))
	if err != nil {
		log.Println("clickedOnSerie: network error ", err)
		return
	}

	var seasonInfo SeasonInfo

	unmarshallErr := json.Unmarshal(responseData, &seasonInfo)
	if unmarshallErr != nil {
		log.Fatalln("clickedOnSeason: json.Unmarshal Failed with", unmarshallErr)
		return
	}

	var episodes []string
	for i, v := range seasonInfo.Episodes {
		episodes = append(episodes, fmt.Sprintf("E%d: %s", i+1, v.Name))
	}

	c.gui.updateView(episodes, c.clickedOnEpisode)

}

func (c *Challenge) clickedOnEpisode(checked bool, episode_number int) {
	log.Println("clickedOnEpisode", episode_number)

	c.episode_number = episode_number
	responseData, err := c.network.GetRequestWith(c.network.MakeEpisodeUrl(c.tvid, c.season_number, c.episode_number))
	if err != nil {
		log.Println("clickedOnSerie: network error ", err)
		return
	}

	var episodeInfo EpisodeInfo

	unmarshallErr := json.Unmarshal(responseData, &episodeInfo)
	if unmarshallErr != nil {
		log.Fatalln("clickedOnSerie: json.Unmarshal Failed with", unmarshallErr)
		return
	}

	c.gui.putOverview(episodeInfo.Name, episodeInfo.Overview)
}
