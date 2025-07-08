package main

import (
	"database/sql"
	"fmt"
	"le-tour-dashmore-server/dataextractor"
	"le-tour-dashmore-server/model"

	_ "github.com/lib/pq"
)

type TableTitleLocator struct {
	Element   string
	TitleName string
}

var db *sql.DB

func main() {

	res, err := dataextractor.GetStageResults("http://localhost:8000/stage-2.html")
	if err != nil {
		fmt.Println(err)
	}

	for k, v := range res {
		fmt.Printf("%s: GC_RANK: %v; GC_TIME: %v; PNTS_RANK: %v; PNTS: %v; KOM_RANK: %v; KOM: %v; YOUTH_RANK: %v\n", k, v.TimeRanking, v.Time, v.SprintRanking, v.Sprint, v.ClimberRanking, v.Climber, v.YoungRiderRanking)
	}
	//connStr := "user=postgres dbname=letour sslmode=disable"
	/*
		connStr := "postgres://<username>:<password>@<server>/letour?sslmode=disable"
		var err error
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			fmt.Println(err)
		}

	*/
	dataextractor.GetRaceDetails()

}

func addRiderToDatabase(rider *model.Rider) error {
	sqlStatement := "INSERT INTO rider (name, team, points) VALUES ($1, $2, $3)"
	_, err := db.Exec(sqlStatement, rider.Name, rider.Team, rider.Points)
	return err
}

func addStagesToDatabase(stage *model.Stage) error {
	sqlStatement := "INSERT INTO stage (name, date) VALUES ($1, $2)"
	_, err := db.Exec(sqlStatement, stage.Name, stage.Date)
	return err
}

func getAllRidersFromDatabase() error {
	rows, err := db.Query("SELECT * FROM rider")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var rider model.Rider
		if err := rows.Scan(&rider.ID, &rider.Name, &rider.Team, &rider.Points); err != nil {
			fmt.Println(err)
		}
		fmt.Println(rider)
	}
	return nil
}

func getAllStagesFromDatabase() error {
	rows, err := db.Query("SELECT * FROM stage")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var stage model.Stage
		if err := rows.Scan(&stage.ID, &stage.Name, &stage.Date); err != nil {
			fmt.Println(err)
		}
		fmt.Println(stage)
	}
	return nil
}
