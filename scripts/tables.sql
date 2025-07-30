DROP TABLE IF EXISTS points_result CASCADE;
DROP TABLE IF EXISTS timed_result CASCADE;
DROP TABLE IF EXISTS timed_result_point_allocation CASCADE;
DROP TABLE IF EXISTS jersey CASCADE;
DROP TABLE IF EXISTS stage CASCADE;
DROP TABLE IF EXISTS split CASCADE;
DROP TABLE IF EXISTS rider CASCADE;
DROP TABLE IF EXISTS team CASCADE;
DROP TABLE IF EXISTS jersey_ranking CASCADE;
DROP TABLE IF EXISTS jersey_ranking_point_allocation CASCADE;
DROP TABLE IF EXISTS goal CASCADE;
DROP TABLE IF EXISTS power_up CASCADE;
DROP TABLE IF EXISTS fantasy_team CASCADE;
DROP TABLE IF EXISTS fantasy_team_rider CASCADE;
DROP TABLE IF EXISTS fantasy_team_goal CASCADE;
DROP TABLE IF EXISTS fantasy_team_power_up CASCADE;
DROP TABLE IF EXISTS goal_point_allocation CASCADE;
DROP TABLE IF EXISTS power_up_point_allocation CASCADE;

DROP TYPE IF EXISTS jersey_type CASCADE;
DROP TYPE IF EXISTS goal_type CASCADE;
DROP TYPE IF EXISTS classification_type CASCADE;
DROP TYPE IF EXISTS power_up_type CASCADE;

CREATE TABLE IF NOT EXISTS team(
	id SERIAL PRIMARY KEY,
	name VARCHAR UNIQUE,
	url VARCHAR,
	jersey_image bytea
);

CREATE TABLE IF NOT EXISTS rider(
	id SERIAL PRIMARY KEY,
	name VARCHAR,
	team_id INT REFERENCES team (id) NOT NULL,
	points INT,
	UNIQUE (name, team_id)
);

CREATE TABLE IF NOT EXISTS stage(
	id SERIAL PRIMARY KEY,
	name VARCHAR UNIQUE,
	date DATE,
	url VARCHAR
);

CREATE TABLE IF NOT EXISTS split(
	id SERIAL PRIMARY KEY,
	name VARCHAR,
	stage_id INT REFERENCES stage (id) NOT NULL,
	UNIQUE (stage_id, name)
);

CREATE TYPE classification_type AS ENUM ('stage', 'gc', 'youth', 'points', 'kom');

CREATE TABLE IF NOT EXISTS jersey(
	id SERIAL PRIMARY KEY,
	type classification_type UNIQUE
);

-- Stage
-- Youth Stage
CREATE TABLE IF NOT EXISTS timed_result(
	id SERIAL PRIMARY KEY,
	rider_id INT REFERENCES rider (id),
	stage_id INT REFERENCES stage (id),
	jersey_id INT REFERENCES jersey (id),
	rank INT,
	time INT,
	time_str VARCHAR,
	points INT,
	UNIQUE (rider_id, stage_id, jersey_id)
);

-- Points (intermediate and sprint finish)
-- KOM (intermediate and climb finishes)
CREATE TABLE IF NOT EXISTS points_result(
	id SERIAL PRIMARY KEY,
	rider_id INT REFERENCES rider (id),
	stage_id INT REFERENCES stage (id),
	jersey_id INT REFERENCES jersey (id),
	split_id INT REFERENCES split (id),
	rank INT,
	points INT,
	UNIQUE (rider_id, stage_id, jersey_id, split_id)
);

-- Ranking as is at the end of each stage for the following
-- GC
-- Points
-- KOM
-- Youth
CREATE TABLE IF NOT EXISTS jersey_ranking(
	id SERIAL PRIMARY KEY,
	rider_id INT REFERENCES rider (id),
	stage_id INT REFERENCES stage (id),
	jersey_id INT REFERENCES jersey (id),
	rank INT,
	points INT,
	UNIQUE (rider_id, stage_id, jersey_id)
);

CREATE TABLE IF NOT EXISTS goal(
	id SERIAL PRIMARY KEY,
	type classification_type UNIQUE
);

CREATE TYPE power_up_type AS ENUM ('stage', 'points', 'kom', 'youth');

CREATE TABLE IF NOT EXISTS power_up(
	id SERIAL PRIMARY KEY,
	type power_up_type UNIQUE
);

CREATE TABLE IF NOT EXISTS goal_point_allocation(
	id SERIAL PRIMARY KEY,
	goal_id INT REFERENCES goal (id),
	position INT,
	bonus INT,
	UNIQUE (goal_id, position)
);

CREATE TABLE IF NOT EXISTS power_up_point_allocation(
	id SERIAL PRIMARY KEY,
	goal_id INT REFERENCES power_up (id),
	position INT,
	bonus INT,
	UNIQUE (goal_id, position)
);

CREATE TABLE IF NOT EXISTS fantasy_team(
	id SERIAL PRIMARY KEY,
	name VARCHAR,
	points int
);

CREATE TABLE IF NOT EXISTS fantasy_team_rider(
	id SERIAL PRIMARY KEY,
	team_id INT REFERENCES fantasy_team(id),
	rider_id INT REFERENCES rider(id),
	UNIQUE (team_id, rider_id)
);

CREATE TABLE IF NOT EXISTS fantasy_team_goal(
	id SERIAL PRIMARY KEY,
	team_id INT REFERENCES fantasy_team (id),
	goal_id INT REFERENCES goal (id),
	points INT,
	UNIQUE (team_id)
);

CREATE TABLE IF NOT EXISTS fantasy_team_power_up(
	id SERIAL PRIMARY KEY,
	team_id INT REFERENCES team (id),
	power_id INT REFERENCES power_up (id),
	rider_id INT REFERENCES rider (id),
	stage_id INT REFERENCES stage (id)
);

INSERT INTO jersey ( type ) 
VALUES
	( 'stage' ),
	( 'gc' ),
	( 'youth' ),
	( 'points' ),
	( 'kom' );


INSERT INTO goal ( type ) 
VALUES
	( 'stage' ),
	( 'gc' ),
	( 'youth' ),
	( 'points' ),
	( 'kom' );

INSERT INTO power_up ( type )
VALUES
	( 'stage' ),
	( 'points' ),
	( 'kom' ),
	( 'youth' );

-- table to define how many points should be awarded for a stage, and youth finish
CREATE TABLE IF NOT EXISTS timed_result_point_allocation(
	id SERIAL PRIMARY KEY,
	jersey_id INT REFERENCES jersey (id),
	rank INT,
	allotted_points INT,
	UNIQUE (jersey_id, rank)
);

INSERT INTO timed_result_point_allocation (jersey_id, rank, allotted_points)
VALUES
-- stage
	(1, 1, 20),
	(1, 2, 17),
	(1, 3, 15),
	(1, 4, 13),
	(1, 5, 11),
	(1, 6, 10),
	(1, 7, 9),
	(1, 8, 8),
	(1, 9, 7),
	(1, 10, 6),
	(1, 11, 5),
	(1, 12, 4),
	(1, 13, 3),
	(1, 14, 2),
	(1, 15, 1),
-- stage youth
	(3, 1, 20),
	(3, 2, 17),
	(3, 3, 15),
	(3, 4, 13),
	(3, 5, 11),
	(3, 6, 10),
	(3, 7, 9),
	(3, 8, 8),
	(3, 9, 7),
	(3, 10, 6),
	(3, 11, 5),
	(3, 12, 4),
	(3, 13, 3),
	(3, 14, 2),
	(3, 15, 1);

-- table to define how many points should be awarded for a position in jersey ranking at the end of each stage
CREATE TABLE IF NOT EXISTS jersey_ranking_point_allocation(
	id SERIAL PRIMARY KEY,
	jersey_id INT REFERENCES jersey (id),
	rank INT,
	allotted_points INT,
	UNIQUE (jersey_id, rank)
);

INSERT INTO jersey_ranking_point_allocation (jersey_id, rank, allotted_points)
VALUES
-- gc
	(2, 1, 20),
	(2, 2, 17),
	(2, 3, 15),
	(2, 4, 13),
	(2, 5, 11),
	(2, 6, 10),
	(2, 7, 9),
	(2, 8, 8),
	(2, 9, 7),
	(2, 10, 6),
	(2, 11, 5),
	(2, 12, 4),
	(2, 13, 3),
	(2, 14, 2),
	(2, 15, 1),
-- youth
	(3, 1, 20),
	(3, 2, 17),
	(3, 3, 15),
	(3, 4, 13),
	(3, 5, 11),
	(3, 6, 10),
	(3, 7, 9),
	(3, 8, 8),
	(3, 9, 7),
	(3, 10, 6),
	(3, 11, 5),
	(3, 12, 4),
	(3, 13, 3),
	(3, 14, 2),
	(3, 15, 1),
--- points
	(4, 1, 20),
	(4, 2, 17),
	(4, 3, 15),
	(4, 4, 13),
	(4, 5, 11),
	(4, 6, 10),
	(4, 7, 9),
	(4, 8, 8),
	(4, 9, 7),
	(4, 10, 6),
	(4, 11, 5),
	(4, 12, 4),
	(4, 13, 3),
	(4, 14, 2),
	(4, 15, 1),
-- kom
	(5, 1, 20),
	(5, 2, 17),
	(5, 3, 15),
	(5, 4, 13),
	(5, 5, 11),
	(5, 6, 10),
	(5, 7, 9),
	(5, 8, 8),
	(5, 9, 7),
	(5, 10, 6),
	(5, 11, 5),
	(5, 12, 4),
	(5, 13, 3),
	(5, 14, 2),
	(5, 15, 1);


INSERT INTO goal_point_allocation (goal_id, position, bonus)
VALUES
-- stage
	(1, 1, 50),
	(1, 2, 25),
	(1, 3, 10),
-- gc
	(2, 1, 200),
	(2, 2, 100),
	(2, 3, 50),
-- youth
	(3, 1, 200),
	(3, 2, 100),
	(3, 3, 50),
--- points
	(4, 1, 200),
	(4, 2, 100),
	(4, 3, 50),
-- kom
	(5, 1, 200),
	(5, 2, 100),
	(5, 3, 50);
