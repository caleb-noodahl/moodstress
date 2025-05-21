package app

import (
	"log"

	"github.com/NimbleMarkets/ntcharts/barchart"
	"github.com/NimbleMarkets/ntcharts/linechart/timeserieslinechart"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/cockroachdb/pebble"
)

type App struct {
	conf                map[string]string
	showForm            bool
	showTable           bool
	showBarChart        bool
	showTimeseriesChart bool
	form                *huh.Form
	table               table.Model
	barchart            barchart.Model
	timeseries          timeserieslinechart.Model
	selected            string
	intralogview        *IntraLogView
	db                  *pebble.DB
	height              int
	width               int
}

func NewApp(conf map[string]string) *App {
	db, err := pebble.Open("db", &pebble.Options{})
	if err != nil {
		log.Panic(err)
	}

	return &App{
		conf:         conf,
		intralogview: NewIntraLogView(),
		db:           db,
		table:        table.New(),
	}
}

func (a *App) Init() tea.Cmd {
	a.showTable = false
	a.showForm = true

	a.form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().Options(
				huh.NewOption("Today", "today"),
				huh.NewOption("Log", "log"),
				huh.NewOption("Metrics", "metrics"),
				huh.NewOption("Graph", "graph"),
			).Value(&a.selected),
		),
	)
	a.form.SubmitCmd = func() tea.Msg {
		switch a.selected {
		case "today":
			a.Today()
		case "log":
			a.IntraLog()
		case "metrics":
		case "graph":
			a.Graph()
		}
		return a.form.Init()
	}
	return tea.Batch([]tea.Cmd{a.form.Init(), tea.WindowSize()}...)
}

func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.height = msg.Height
		a.width = msg.Width
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			return a, tea.Quit
		case "enter":
			if a.table.Focused() {
				a.table.Blur()
				return a, a.Init()
			}
		}
	}

	cmds := []tea.Cmd{}
	var cmd tea.Cmd

	form, cmd := a.form.Update(msg)
	if form, ok := form.(*huh.Form); ok {
		a.form = form
	}
	cmds = append(cmds, cmd)

	tbl, cmd := a.table.Update(msg)
	a.table = tbl
	cmds = append(cmds, cmd)

	cht, cmd := a.barchart.Update(msg)
	a.barchart = cht
	cmds = append(cmds, cmd)

	return a, tea.Batch(cmds...)
}

func (a *App) View() string {
	result := ""
	if a.showForm {
		result = a.form.View()
	}
	if a.showTable {
		result += "\n" + a.table.View()
	}
	if a.showBarChart {
		result += "\n" + a.barchart.View()
	}
	if a.showTimeseriesChart {
		result += "\n" + a.timeseries.View()
	}
	return result
	//return fmt.Sprintf("%s\n\n%s", a.form.View(), a.table.View())
}
