package model

type MainMenuAction struct {
	Name string
	Idx  int
}

var MainMenuActionStart = MainMenuAction{Name: "Start", Idx: 0}
var MainMenuActionStop = MainMenuAction{Name: "Stop", Idx: 1}
var MainMenuActionShowProjects = MainMenuAction{Name: "Show Projects", Idx: 2}
var MainMenuActionCheckBalance = MainMenuAction{Name: "Check Balance", Idx: 3}
var MainMenuActionSetPartTime = MainMenuAction{Name: "Set Part Time", Idx: 4}
var MainMenuActionChangeAPIKey = MainMenuAction{Name: "Change API key", Idx: 5}
var MainMenuActionQuit = MainMenuAction{Name: "Quit", Idx: 6}

var MainMenuActions = []MainMenuAction{
	MainMenuActionStart,
	MainMenuActionStop,
	MainMenuActionShowProjects,
	MainMenuActionCheckBalance,
	MainMenuActionSetPartTime,
	MainMenuActionChangeAPIKey,
	MainMenuActionQuit,
}
