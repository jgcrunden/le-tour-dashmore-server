package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	Time              string
	TimeRanking       int `json:"time_ranking"`
	SprintRanking     int `json:"sprint_ranking"`
	YoungRiderRanking int `json:"young_rider_ranking"`
	ClimberRanking    int `json:"climber_ranking"`
	Sprint            int `json:"sprint"`
	Climber           int `json:"climber"`
}

type DataItem interface {
	SetField(fieldName string, value string)
}

type Stage struct {
	ID   int
	Name string
	Date time.Time
}

func NewStage() *Stage {
	return &Stage{}
}

func (s *Stage) SetField(fieldName string, value string) {
	switch fieldName {
	case "Date":
		// TODO remove hardcoded 2025 year
		t, err := time.Parse("2006-02/01", fmt.Sprintf("%d-%s", 2025, value))
		if err != nil {
			fmt.Println(err)
			return
		}
		s.Date = t
	case "Stage":
		s.Name = value
	}
}

type Rider struct {
	ID     int
	Name   string
	Team   string
	Points int
}

func NewRider() *Rider {
	return &Rider{}
}

func (r *Rider) SetField(fieldName string, value string) {
	switch fieldName {
	case "Rider":
		r.Name = value
	case "Team":
		r.Team = value
	case "Points":
		points, _ := strconv.Atoi(value)
		r.Points = points
	}
}

type GCResult struct {
	Time int
	Rank int
}

type TimedResult struct {
	Rank    int
	Rider   string
	TimeStr string
	Time    int
}

func NewTimedResult() *TimedResult {
	return &TimedResult{}
}

func (r *TimedResult) SetField(fieldName string, value string) {
	switch fieldName {
	case "Rnk":
		rank, _ := strconv.Atoi(value)
		r.Rank = rank
	case "Rider":
		r.Rider = value
	case "Time":
		r.TimeStr = value
	}
}

func (r *TimedResult) ParseTime(leadersTime int) error {
	var h, m, s int
	var err error
	numOfColons := strings.Count(r.TimeStr, ":")
	if numOfColons == 0 {
		_, err = fmt.Sscanf(r.TimeStr, "%d", &s)
	} else if numOfColons == 1 {
		_, err = fmt.Sscanf(r.TimeStr, "%d:%d", &m, &s)
	} else if numOfColons == 2 {
		_, err = fmt.Sscanf(r.TimeStr, "%d:%d:%d", &h, &m, &s)
	}

	if err != nil {
		return err
	}
	res := h*3600 + m*60 + s
	if r.Rank != 1 {
		res += leadersTime
	}
	r.Time = res
	return nil
}

func (r *TimedResult) FixDitto(previousTime string, leadersTime string) string {
	newTime := r.TimeStr
	if r.TimeStr == ",," && previousTime == leadersTime {
		newTime = "0:00"
	} else if r.TimeStr == ",," && previousTime != leadersTime {
		newTime = previousTime
	}
	if r.Rank != 1 {
		r.TimeStr = /*"+" +*/ newTime
	}
	return newTime
}

type PointsResult struct {
	Rank   int
	Rider  string
	Points int
}

func NewPointsResult() *PointsResult {
	return &PointsResult{}
}

func (r *PointsResult) SetField(fieldName string, value string) {
	switch fieldName {
	case "Rnk":
		r.Rank, _ = strconv.Atoi(value)
	case "Rider":
		r.Rider = value
	case "Pnt":
		r.Points, _ = strconv.Atoi(value)
	}
}
