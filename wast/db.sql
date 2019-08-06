/*
|----------------------------------------------------------
| Schema
|----------------------------------------------------------
*/
DROP TABLE IF EXISTS ranks;
DROP TABLE IF EXISTS stats;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS players;

CREATE TABLE players (
  id SERIAL NOT NULL,
  nickname VARCHAR(255) NOT NULL,
  discord_tag VARCHAR(255) NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE games (
  id SERIAL NOT NULL,
  created_at DATE NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE stats (
  id SERIAL NOT NULL,
  kills SMALLINT,
  deaths SMALLINT,
  damage SMALLINT,
  winner BOOLEAN,
  player_id INTEGER NOT NULL,
  game_id INTEGER NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_player_id FOREIGN KEY(player_id) REFERENCES players(id),
  CONSTRAINT fk_game_id FOREIGN KEY(game_id) REFERENCES games(id)
);

CREATE TABLE ranks (
  id SERIAL NOT NULL,
  name VARCHAR(255) NOT NULL,
  min_elo INTEGER NOT NULL,
  max_elo INTEGER NOT NULL,
  PRIMARY KEY(id)
);

/*
|----------------------------------------------------------
| Seeds
|----------------------------------------------------------
*/
INSERT INTO ranks (name, min_elo, max_elo) VALUES ('Leonded', 0, 19);
INSERT INTO ranks (name, min_elo, max_elo) VALUES ('Achier', 20, 39);
INSERT INTO ranks (name, min_elo, max_elo) VALUES ('Paouf', 40, 59);
INSERT INTO ranks (name, min_elo, max_elo) VALUES ('Padégueu', 60, 79);
INSERT INTO ranks (name, min_elo, max_elo) VALUES ('Moyen moins', 80, 99);
INSERT INTO ranks (name, min_elo, max_elo) VALUES ('Bien mais pas top', 100, 119);
INSERT INTO ranks (name, min_elo, max_elo) VALUES ('Barack Obanane', 120, 139);
INSERT INTO ranks (name, min_elo, max_elo) VALUES ('Maître Bananier', 140, 159);
INSERT INTO ranks (name, min_elo, max_elo) VALUES ('Banana Jesus', 160, 179);
INSERT INTO ranks (name, min_elo, max_elo) VALUES ('Worms God', 180, 200);