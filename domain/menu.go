package domain

type MenuAction struct {
	Name string
	Idx  int
}

var MenuActionStart = MenuAction{Name: "Start", Idx: 0}
var MenuActionStop = MenuAction{Name: "Stop", Idx: 1}
var MenuActionShowProjects = MenuAction{Name: "Show Projects", Idx: 2}
var MenuActionCheckBalance = MenuAction{Name: "Check Balance", Idx: 3}
var MenuActionSetPartTime = MenuAction{Name: "Set part-time", Idx: 4}
var MenuActionChangeAPIKey = MenuAction{Name: "Change API key", Idx: 5}
var MenuActionQuit = MenuAction{Name: "Quit", Idx: 6}
var MenuActionVerifyMonth = MenuAction{Name: "Verify Month", Idx: 7}

var MenuActions = []MenuAction{
	MenuActionStart,
	MenuActionStop,
	MenuActionShowProjects,
	MenuActionCheckBalance,
	MenuActionVerifyMonth,
	MenuActionSetPartTime,
	MenuActionChangeAPIKey,
	MenuActionQuit,
}
