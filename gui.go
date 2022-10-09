package main

import (
	"fmt"
	"log"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

const (
	NR_OF_VISIBLE_SERIES = 10
	SEARCH_BUTTON_TITLE  = "Search"
	WINDOW_TITLE         = "Challenge"
	MIN_WINDOW_X         = 400
	MIN_WINDOW_Y         = 400
)

type Gui struct {
	SeriesOnGrid    int
	MaxSeriesOnGrid int

	widget       *widgets.QWidget
	scrollWidget *widgets.QWidget
	scrollLayout *widgets.QVBoxLayout
	scrollArea   *widgets.QScrollArea

	window *widgets.QMainWindow
	layout *widgets.QVBoxLayout

	input *widgets.QLineEdit
	text  *widgets.QLabel

	searchButton *widgets.QPushButton
	buttonsGrid  []*widgets.QPushButton
}

func NewGui() *Gui {
	g := Gui{}
	g.MaxSeriesOnGrid = NR_OF_VISIBLE_SERIES

	widgets.NewQApplication(len(os.Args), os.Args)

	g.window = widgets.NewQMainWindow(nil, 0)
	g.window.SetWindowTitle(WINDOW_TITLE)
	g.window.SetMinimumSize2(MIN_WINDOW_X, MIN_WINDOW_Y)

	g.layout = widgets.NewQVBoxLayout()
	g.widget = widgets.NewQWidget(nil, 0)
	g.widget.SetLayout(g.layout)

	g.scrollLayout = widgets.NewQVBoxLayout()
	g.scrollWidget = widgets.NewQWidget(nil, 0)
	g.scrollWidget.SetLayout(g.scrollLayout)

	g.input = widgets.NewQLineEdit(nil)

	g.layout.AddWidget(g.input, 0, 0)
	g.window.SetCentralWidget(g.widget)

	g.searchButton = widgets.NewQPushButton2(SEARCH_BUTTON_TITLE, nil)
	g.layout.AddWidget(g.searchButton, 0, 0)

	g.scrollArea = widgets.NewQScrollArea(nil)
	g.scrollArea.SetVerticalScrollBarPolicy(core.Qt__ScrollBarAlwaysOn)
	g.scrollArea.SetHorizontalScrollBarPolicy(core.Qt__ScrollBarAsNeeded)
	g.scrollArea.SetEnabled(true)
	g.scrollArea.SetWidgetResizable(true)
	g.layout.AddWidget(g.scrollArea, 0, 0)
	g.scrollArea.SetWidget(g.scrollWidget)

	return &g
}

func (g *Gui) ResetResultGrid() {
	log.Println("ResetResultGrid called")

	for i := 0; i < g.SeriesOnGrid; i++ {
		g.buttonsGrid[i].DisconnectClicked()
		g.buttonsGrid[i].DeleteLater()
		g.buttonsGrid[i] = nil
	}

	if g.text != nil {
		g.text.DeleteLater()
		g.text = nil
	}

	g.buttonsGrid = make([]*widgets.QPushButton, 0)
	g.MaxSeriesOnGrid = NR_OF_VISIBLE_SERIES
	g.SeriesOnGrid = 0
}

func (g *Gui) updateViewWithIdsStoring(ids []int, elemNames []string, callback func(bool, int), cont bool) {
	log.Println("updateViewWithIdsStoring", elemNames)

	for i, v := range elemNames {
		button := widgets.NewQPushButton2(v, nil)
		tvid := ids[i]
		button.ConnectClicked(func(checked bool) {
			go callback(checked, tvid)
		})
		g.buttonsGrid = append(g.buttonsGrid, button)
		g.scrollLayout.AddWidget(button, 0, 0)
	}
	g.SeriesOnGrid += len(elemNames)

}

func (g *Gui) updateView(elemNames []string, callback func(bool, int)) {
	log.Println("updateView", elemNames)

	var ids []int

	for i := 0; i < len(elemNames); i++ {
		ids = append(ids, i)
	}

	g.ResetResultGrid()
	g.updateViewWithIdsStoring(ids, elemNames, callback, false)

}

func (g *Gui) putOverview(name string, overview string) {
	log.Println("putOverview", name, overview)

	g.ResetResultGrid()

	g.text = widgets.NewQLabel(nil, 0)
	g.text.SetText(fmt.Sprintf("%s\n\n%s\n", name, overview))
	g.text.SetWordWrap(true)

	g.scrollLayout.AddWidget(g.text, 0, 0)

}
