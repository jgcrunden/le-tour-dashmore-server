package main

import (
	"flag"
	"fmt"
	"le-tour-dashmore-server/dataextractor"
	"le-tour-dashmore-server/dbclient/databases"
	"le-tour-dashmore-server/model"
	"le-tour-dashmore-server/teammanager"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/magiconair/properties"
)

type TableTitleLocator struct {
	Element   string
	TitleName string
}

//var db *sql.DB
var cfg model.Config

func main() {

	confFilePtr := flag.String("f", "/etc/le-tour-dashmore-server/server.conf", "path to config file")
	flag.Parse()
	p := properties.MustLoadFile(*confFilePtr, properties.UTF8)
	var cfg model.Config

	if err := p.Decode(&cfg); err != nil {
		fmt.Println("Error reading config file", err)
		return
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}

func App() {
	dbClient := databases.NewPostgresClient(cfg.DBHost, cfg.DBUser, cfg.DBPassword)
	extractor := dataextractor.NewExtractor(dbClient, cfg)
	extractor.GetRaceDetails()

	stages, err := dbClient.GetAllStages()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, v := range stages {
		_, err := extractor.GetStageResults(v)
		if err != nil {
			fmt.Println(err)
		}
	}

	err = extractor.CalculatePoints()
	if err != nil {
		fmt.Println("Error calculating results")
	}
	manager := teammanager.NewTeamManager(dbClient)
	/*
		_ = model.FantasyTeam {
			Name:   "My team",
			Goal:   model.GC,
			Points: 0,
			Riders: []model.Rider{
				{
					Name:   "PHILIPSEN Jasper",
				},
				{
					Name:   "EVENEPOEL Remco",
				},
				{
					Name: "PEDERSEN Mads",
				},
				{
					Name: "VINGEGAARD Jonas",
				},
				{
					Name: "YATES Adam",
				},
				{
					Name: "VAN DER POEL Mathieu",
				},
				{
					Name: "VAN AERT Wout",
				},
				{
					Name: "ROGLIČ Primož",
				},
			},
		}

	*/
	// The boys #FFA500
	t := model.CreateTeam("Les Bebes Sauvages", [8]string{"POGAČAR Tadej", "YATES Adam", "DRIZNERS Jarrad", "MADOUAS Valentin", "ASKEY Lewis", "TEUNISSEN Mike", "STORER Michael", "STUYVEN Jasper"})
	err = manager.CreateFantasyTeam(t)
	// Joe green
	t = model.CreateTeam("Aldi Stroll", [8]string{"O'CONNOR Ben", "RODRÍGUEZ Carlos", "ROMEO Iván", "GRIGNARD Sébastien", "WRIGHT Fred", "DE LIE Arnaud", "REINDERS Elmar", "FOSS Tobias"})
	err = manager.CreateFantasyTeam(t)
	// Freya #ff0
	t = model.CreateTeam("Not Very Good at Cycling", [8]string{"ONLEY Oscar", "ROGLIČ Primož", "BENOOT Tiesj", "SEPÚLVEDA Eduardo", "GRADEK Kamil", "MAYRHOFER Marius", "NABERMAN Tim", "LOUVEL Matis"})
	err = manager.CreateFantasyTeam(t)
	// Josh #c24646
	t = model.CreateTeam("Hanover Elite E-Bike Team", [8]string{"JORGENSON Matteo", "EVENEPOEL Remco", "SCHACHMANN Maximilian", "CAPIOT Amaury", "COSTIOU Ewen", "HOELGAARD Markus", "CONSONNI Simone", "SKUJIŅŠ Toms"})
	err = manager.CreateFantasyTeam(t)
	// Alison #005EB8
	t = model.CreateTeam("Don't Give Me o'Ganna", [8]string{"MAS Enric", "VINGEGAARD Jonas", "GROENEWEGEN Dylan", "ARANBURU Alex", "PARET-PEINTRE Valentin", "TURGIS Anthony", "VAN MOER Brent", "CAMPENAERTS Victor"})
	err = manager.CreateFantasyTeam(t)
	// Richard purple
	t = model.CreateTeam("Tour du Monde en Luberon", [8]string{"YATES Simon", "ALMEIDA João", "HAIG Jack", "VAN DEN BERG Marijn", "POLITT Nils", "ABRAHAMSEN Jonas", "LE BERRE Mathis", "DURBRIDGE Luke"})
	err = manager.CreateFantasyTeam(t)
	//err = manager.CreateFantasyTeam(team1)
	if err != nil {
		fmt.Printf("Error creating team %s: %s\n", t.Name, err)
		return
	}

	// TODO: calculate if tour has finished, store as boolean and pass into CalculateGoalScores
	err = manager.CalculateGoalScores(true)
	if err != nil {
		fmt.Println("Unable to calculate goal scores", err)
		return
	}
	err = manager.CalculateTeamScores()
	if err != nil {
		fmt.Println(err)
		return
	}
	/*
		for k, v := range res {
			fmt.Printf("%s: GC_RANK: %v; GC_TIME: %v; PNTS_RANK: %v; PNTS: %v; KOM_RANK: %v; KOM: %v; YOUTH_RANK: %v\n", k, v.TimeRanking, v.Time, v.SprintRanking, v.Sprint, v.ClimberRanking, v.Climber, v.YoungRiderRanking)
		}
	*/
	//connStr := "user=postgres dbname=letour sslmode=disable"

}
