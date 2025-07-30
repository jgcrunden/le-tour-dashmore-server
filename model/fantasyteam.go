package model

type FantasyTeam struct {
	ID     int
	Name   string
	Goal   GoalType
	Points int
	Riders []Rider
}

type GoalType int

// Goal type
const (
	STAGE GoalType = iota + 1
	GC
	YOUTH
	POINTS
	KOM
)

var goalTypeName = map[GoalType]string{
	STAGE:  "stage",
	GC:     "gc",
	YOUTH:  "youth",
	POINTS: "points",
	KOM:    "kom",
}

func (g GoalType) String() string {
	return goalTypeName[g]
}

type Goal struct {
	ID   int
	Type int
}

func CreateTeam(teamName string, riderNames [8]string) FantasyTeam {
	riders := make([]Rider, 0)

	for _, i := range riderNames {
		r := Rider{
			Name: i,
		}
		riders = append(riders, r)
	}
	return FantasyTeam{
		Name:   teamName,
		Goal:   GC,
		Points: 0,
		Riders: riders,
	}
}
