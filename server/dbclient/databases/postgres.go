package databases

import (
	"database/sql"
	"fmt"
	"le-tour-dashmore-server/dbclient"
	"le-tour-dashmore-server/model"
	"strings"
)

type PostgresClient struct {
	db *sql.DB
}

// GetGoalPoints implements dbclient.DBClient.
func (p PostgresClient) GetGoalPoints(teamID int) (int, error) {
	row := p.db.QueryRow("SELECT points FROM fantasy_team_goal WHERE team_id = $1", teamID)
	var points int
	err := row.Scan(&points)
	return points, err
}

// GetJerseyRankingForRiderAtStage implements dbclient.DBClient.
func (p PostgresClient) GetJerseyRankingForRiderAtStage(riderID int, stageID int, goal model.GoalType) (int, error) {
	row := p.db.QueryRow("SELECT rank FROM jersey_ranking WHERE rider_id = $1 AND stage_id = $2 AND jersey_id = $3", riderID, stageID, goal)
	var rank int
	err := row.Scan(&rank)
	if err != nil && err.Error() == "sql: no rows in result set" {
		err = nil
	}
	return rank, err
}

// UpdateGoalPoints implements dbclient.DBClient.
func (p PostgresClient) UpdateGoalPoints(teamID int, totalGoalPoints int) error {
	_, err := p.db.Exec("UPDATE fantasy_team_goal SET points = $1 WHERE id = $2", totalGoalPoints, teamID)
	return err
}

// GetTimedResultsForRider implements dbclient.DBClient.
func (p PostgresClient) GetTimedResultsForRider(riderID int, jersey model.GoalType) ([]int, error) {
	rows, err := p.db.Query("SELECT rank FROM timed_result WHERE rider_id = $1 AND jersey_id = $2", riderID, int(jersey))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]int, 0)
	for rows.Next() {
		var rank int
		if err := rows.Scan(&rank); err != nil {
			return nil, err
		}
		res = append(res, rank)
	}
	return res, nil
}

// GetGoalPointAllocations implements dbclient.DBClient.
func (p PostgresClient) GetGoalPointAllocations(goal model.GoalType) (map[int]int, error) {
	rows, err := p.db.Query("SELECT position, bonus FROM goal_point_allocation WHERE goal_id = $1", goal)
	if err != nil {
		fmt.Println("Error getting goal point allocations", err)
		return nil, err
	}

	defer rows.Close()

	res := make(map[int]int)
	for rows.Next() {
		var position, points int
		if err := rows.Scan(&position, &points); err != nil {
			return nil, err
		}
		res[position] = points
	}

	return res, nil
}

// CreateFantasyTeamGoal implements dbclient.DBClient.
func (p PostgresClient) CreateFantasyTeamGoal(id int, goal string) error {
	_, err := p.db.Exec("INSERT INTO fantasy_team_goal(team_id, goal_id, points) VALUES ($1, (SELECT id from goal WHERE type = $2), 0)", id, goal)
	return err
}

// UpdateFantasyTeamPoints implements dbclient.DBClient.
func (p PostgresClient) UpdateFantasyTeamPoints(teamID int, points int) error {
	_, err := p.db.Exec("UPDATE fantasy_team SET points = $1 WHERE id = $2", points, teamID)
	return err
}

// GetAllFantasyTeams implements dbclient.DBClient.
func (p PostgresClient) GetAllFantasyTeams() ([]model.FantasyTeam, error) {
	rows, err := p.db.Query("SELECT fantasy_team.id, fantasy_team.name, fantasy_team_goal.goal_id, fantasy_team.points FROM fantasy_team INNER JOIN fantasy_team_goal ON fantasy_team.id = fantasy_team_goal.team_id")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]model.FantasyTeam, 0)
	for rows.Next() {
		var team model.FantasyTeam
		if err := rows.Scan(&team.ID, &team.Name, &team.Goal, &team.Points); err != nil {
			return nil, err
		}
		res = append(res, team)
	}

	for i := range res {
		riders, err := p.getAllRidersInFantasyTeam(res[i].ID)
		if err != nil {
			return nil, err
		}

		res[i].Riders = riders
	}
	return res, nil
}

func (p PostgresClient) getAllRidersInFantasyTeam(teamID int) ([]model.Rider, error) {
	rows, err := p.db.Query(
		`SELECT rider.id, rider.name, team.name, rider.points
		FROM fantasy_team_rider
		INNER JOIN rider ON fantasy_team_rider.rider_id = rider.id
		INNER JOIN team ON rider.team_id = team.id
		WHERE fantasy_team_rider.team_id = $1`, teamID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]model.Rider, 0)
	for rows.Next() {
		var rider model.Rider
		if err := rows.Scan(&rider.ID, &rider.Name, &rider.Team, &rider.Points); err != nil {
			return nil, err
		}
		res = append(res, rider)
	}

	return res, nil
}

// GetRiderByName implements dbclient.DBClient.
func (p PostgresClient) GetRiderByName(riderName string) (model.Rider, error) {
	panic("unimplemented")
}

// AddFantasyRider implements dbclient.DBClient.
func (p PostgresClient) AddFantasyRider(id int, riderName string) error {
	_, err := p.db.Exec("INSERT INTO fantasy_team_rider(team_id, rider_id) VALUES ($1, (SELECT id from rider WHERE name = $2))", id, riderName)
	return err
}

// CreateFantasyTeam implements dbclient.DBClient.
func (p PostgresClient) CreateFantasyTeam(teamName string) (int, error) {
	var id int
	err := p.db.QueryRow("INSERT INTO fantasy_team(name, points) VALUES ($1, $2) RETURNING id", teamName, 0).Scan(&id)
	return id, err
}

// UpdateRiderPoints implements dbclient.DBClient.
func (p PostgresClient) UpdateRiderPoints(riderID int, points int) error {
	_, err := p.db.Exec("UPDATE rider SET points = $1 WHERE id = $2", points, riderID)
	return err
}

// GetJerseyRankingPoints implements dbclient.DBClient.
func (p PostgresClient) GetJerseyRankingPoints(riderID int) (int, error) {
	rows, err := p.db.Query("SELECT points FROM jersey_ranking WHERE rider_id = $1", riderID)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	total := 0
	for rows.Next() {
		var points int
		if err := rows.Scan(&points); err != nil {
			return 0, err
		}
		total += points
	}
	return total, nil
}

// GetPointsResultPoints implements dbclient.DBClient.
func (p PostgresClient) GetPointsResultPoints(riderID int) (int, error) {
	rows, err := p.db.Query("SELECT points from points_result WHERE rider_id = $1", riderID)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	total := 0
	for rows.Next() {
		var points int
		if err := rows.Scan(&points); err != nil {
			return 0, err
		}
		total += points
	}
	return total, nil
}

// GetTimedResultPoints implements dbclient.DBClient.
func (p PostgresClient) GetTimedResultPoints(riderID int) (int, error) {
	rows, err := p.db.Query("SELECT points from timed_result WHERE rider_id = $1", riderID)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	total := 0
	for rows.Next() {
		var points int
		if err := rows.Scan(&points); err != nil {
			return 0, err
		}
		total += points
	}
	return total, nil
}

// AddSplit implements dbclient.DBClient.
func (p PostgresClient) AddSplit(splitName string, stageID int) (int, error) {
	sqlStatement := "INSERT INTO split (name, stage_id) VALUES ($1, $2) RETURNING id"
	var id int
	err := p.db.QueryRow(sqlStatement, splitName, stageID).Scan(&id)
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "split_stage_id_name_key"` {
		err = nil
	}
	return id, err
}

// AddJerseyRanking implements dbclient.DBClient.
func (p PostgresClient) AddJerseyRanking(jerseyResult model.JerseyResult, stageID int, jerseyID int) error {
	rows := p.db.QueryRow("SELECT allotted_points FROM jersey_ranking_point_allocation WHERE jersey_id = $1 AND rank = $2", jerseyID, jerseyResult.GetRank())
	var points int
	if err := rows.Scan(&points); err != nil {
		points = 0
	}

	sqlStatement := "INSERT INTO jersey_ranking (rider_id, stage_id, jersey_id, rank, points) VALUES ( (SELECT id FROM rider WHERE name = $1), $2, $3, $4, $5)"
	_, err := p.db.Exec(sqlStatement, jerseyResult.GetRider(), stageID, jerseyID, jerseyResult.GetRank(), points)
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "jersey_ranking_rider_id_stage_id_jersey_id_key"` {
		err = nil
	}
	return err
}

// AddPointsResult implements dbclient.DBClient.
func (p PostgresClient) AddPointsResult(pointsResult *model.PointsResult, stageID int, jerseyID int, splitID int) error {
	sqlStatement := "INSERT INTO points_result (rider_id, stage_id, jersey_id, split_id, rank, points) VALUES ( (SELECT id FROM rider WHERE name = $1), $2, $3, $4, $5, $6)"
	_, err := p.db.Exec(sqlStatement, pointsResult.Rider, stageID, jerseyID, splitID, pointsResult.Rank, pointsResult.Points)
	if err != nil && err.Error() == `pq: insert or update on table "points_result" violates foreign key constraint "points_result_split_id_fkey"` {
		err = nil
	}
	return err
}

// AddTimedResult implements dbclient.DBClient.
func (p PostgresClient) AddTimedResult(timedResult *model.TimedResult, stageID int, jerseyID int) error {
	rows := p.db.QueryRow("SELECT allotted_points FROM timed_result_point_allocation WHERE jersey_id = $1 AND rank = $2", jerseyID, timedResult.Rank)
	var points int
	if err := rows.Scan(&points); err != nil {
		points = 0
	}

	sqlStatement := "INSERT INTO timed_result (rider_id, stage_id, jersey_id, rank, time, time_str, points) VALUES ( (SELECT id FROM rider WHERE name = $1), $2, $3, $4, $5, $6, $7)"
	_, err := p.db.Exec(sqlStatement, timedResult.Rider, stageID, jerseyID, timedResult.Rank, timedResult.Time, timedResult.TimeStr, points)
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "timed_result_rider_id_stage_id_jersey_id_key"` {
		err = nil
	}
	return err
}

// GetJersey implements dbclient.DBClient.
func (p PostgresClient) GetJersey(name string) (model.Jersey, error) {
	name = strings.ToLower(name)
	res := p.db.QueryRow(`SELECT * FROM jersey WHERE type = $1`, name)
	if res == nil {
		fmt.Println("Unable to find jersey", name)
	}

	jersey := model.Jersey{}
	err := res.Scan(&jersey.ID, &jersey.Type)
	return jersey, err
}

func NewPostgresClient(host string, username string, password string) dbclient.DBClient {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/letour?sslmode=disable", username, password, host)
	var err error
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}
	return PostgresClient{
		db,
	}
}

func (p PostgresClient) AddRider(rider *model.Rider) error {
	sqlStatement := "INSERT INTO rider (name, team_id, points) VALUES ( $1, ( SELECT id FROM team WHERE name = $2), $3 );"
	_, err := p.db.Exec(sqlStatement, rider.Name, rider.Team, rider.Points)
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "rider_name_team_id_key"` {
		err = nil
	}
	return err
}

func (p PostgresClient) AddTeam(team *model.Team) error {
	sqlStatement := "INSERT INTO team (name, url, jersey_image) VALUES ($1, $2, $3)"
	_, err := p.db.Exec(sqlStatement, team.Name, team.TeamURL, team.JerseyImage)
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "team_name_key"` {
		err = nil
	}
	return err
}

func (p PostgresClient) GetAllRiders() ([]model.Rider, error) {
	res := make([]model.Rider, 0)
	rows, err := p.db.Query("SELECT * FROM rider")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rider model.Rider
		if err := rows.Scan(&rider.ID, &rider.Name, &rider.Team, &rider.Points); err != nil {
			return nil, err
		}
		res = append(res, rider)
	}
	return res, nil
}

func (p PostgresClient) AddStage(stage *model.Stage) error {
	sqlStatement := "INSERT INTO stage (name, date, url) VALUES ($1, $2, $3)"
	_, err := p.db.Exec(sqlStatement, stage.Name, stage.Date, stage.URL)
	if err != nil && err.Error() == `pq: duplicate key value violates unique constraint "stage_name_key"` {
		err = nil
	}
	return err
}

func (p PostgresClient) GetAllStages() ([]model.Stage, error) {
	res := make([]model.Stage, 0)
	rows, err := p.db.Query("SELECT * FROM stage")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var stage model.Stage
		if err := rows.Scan(&stage.ID, &stage.Name, &stage.Date, &stage.URL); err != nil {
			return nil, err
		}
		res = append(res, stage)
	}
	return res, nil
}
