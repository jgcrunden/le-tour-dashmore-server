package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	DBHost              string `properties:"db.host"`
	DBUser              string `properties:"db.user"`
	DBPassword          string `properties:"db.password"`
	CyclingStatsBaseURL string `properties:"cycling_stats.baseUrl"`
	TourExtension       string `properties:"cycling_stats.tour_extension"`
}

type LegacyResult struct {
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
	URL  string
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

type Jersey struct {
	ID   int
	Type string
}

type Team struct {
	ID          int
	Name        string
	TeamURL     string
	JerseyImage []byte
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

type JerseyResult interface {
	GetRider() string
	GetRank() int
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

func (r TimedResult) GetRider() string {
	return r.Rider
}

func (r TimedResult) GetRank() int {
	return r.Rank
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
	switch numOfColons {
	case 0:
		_, err = fmt.Sscanf(r.TimeStr, "%d", &s)
	case 1:
		_, err = fmt.Sscanf(r.TimeStr, "%d:%d", &m, &s)
	case 2:
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

func (r PointsResult) GetRider() string {
	return r.Rider
}

func (r PointsResult) GetRank() int {
	return r.Rank
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
