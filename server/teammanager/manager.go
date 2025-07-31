package teammanager

import (
	"fmt"
	"le-tour-dashmore-server/dbclient"
	"le-tour-dashmore-server/model"
)

type TeamManager struct {
	db dbclient.DBClient
}

func NewTeamManager(db dbclient.DBClient) TeamManager {
	return TeamManager{
		db: db,
	}
}

func (m TeamManager) CreateFantasyTeam(team model.FantasyTeam) error {
	id, err := m.db.CreateFantasyTeam(team.Name)
	if err != nil {
		fmt.Println("Error creating team", team.Name)
		return err
	}

	for _, v := range team.Riders {
		err := m.db.AddFantasyRider(id, v.Name)
		if err != nil {
			fmt.Printf("Error adding rider %s: %v", v.Name, err)
			return err
		}
	}

	err = m.db.CreateFantasyTeamGoal(id, team.Goal.String())
	if err != nil {
		fmt.Println("Error creating team goal")
		return err
	}
	return nil
}

func (m TeamManager) CalculateGoalScores(isFinished bool) error {
	fmt.Println("CalculateGoalScores called")
	teams, err := m.db.GetAllFantasyTeams()
	if err != nil {
		fmt.Println("Error getting all fantasy teams", err)
		return err
	}

	stages, err := m.db.GetAllStages()
	if err != nil {
		fmt.Println("Error getting stages", err)
		return err
	}

	for _, team := range teams {
		goalPointAllocations, err := m.db.GetGoalPointAllocations(team.Goal)
		if err != nil {
			fmt.Println("Error getting goal point allocations", err)
			return err
		}
		totalGoalPoints := 0

		if team.Goal == model.STAGE {
			// calculate goal points for stage wins
			for _, rider := range team.Riders {
				stageRankings, err := m.db.GetTimedResultsForRider(rider.ID, model.STAGE)
				if err != nil {
					fmt.Println("Error getting timed results for rider", rider.Name, err)
					return err
				}

				for _, ranking := range stageRankings {
					totalGoalPoints += goalPointAllocations[ranking]
				}
			}
		} else if isFinished {
			for _, rider := range team.Riders {
				rank, err := m.db.GetJerseyRankingForRiderAtStage(rider.ID, stages[len(stages)-1].ID, team.Goal)
				if err != nil {
					fmt.Println("Error getting jersy ranking for stage", err)
					return err
				}
				fmt.Println("rider:", rider.Name, "rank:", rank, "points:", goalPointAllocations[rank])
				totalGoalPoints += goalPointAllocations[rank]
			}
		}
		if team.Goal == model.STAGE || isFinished {
			if err := m.db.UpdateGoalPoints(team.ID, totalGoalPoints); err != nil {
				fmt.Println("Error updating goal points", err)
				return err
			}
		}
	}
	return nil
}

func (m TeamManager) CalculateTeamScores() error {
	teams, err := m.db.GetAllFantasyTeams()
	if err != nil {
		fmt.Println("Error getting all fantasy teams", err)
	}
	for _, team := range teams {
		totalPoints := 0
		for _, rider := range team.Riders {
			totalPoints += rider.Points
		}

		points, err := m.db.GetGoalPoints(team.ID)
		if err != nil {
			fmt.Println("Unable to get goal points", err)
			return err
		}

		totalPoints += points
		err = m.db.UpdateFantasyTeamPoints(team.ID, totalPoints)
		if err != nil {
			fmt.Println("Error updating team points", err)
		}

	}
	return nil
}
