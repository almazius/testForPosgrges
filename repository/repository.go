package repository

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
	models "testForPosgrges"
	"testForPosgrges/config"
)

var Connect *sqlx.DB

func InitPsqlDB(c *config.Config) (*sqlx.DB, error) {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		*c.Postgres.Host,
		*c.Postgres.Port,
		*c.Postgres.User,
		*c.Postgres.Password,
		*c.Postgres.DbName)
	database, err := sqlx.Connect("pgx", connectionUrl)
	if err != nil {
		return nil, err
	}
	return database, nil
}

func GetHashFromRepository(dir *models.Directory) (map[string]string, error) {
	filesHash := make(map[string]string)
	hash := ""
	for file, _ := range dir.FileHash {
		err := Connect.QueryRowx("SELECT hash FROM hashes WHERE path = $1", file).Scan(&hash)
		if err != sql.ErrNoRows && err != nil {
			if dir.LogThread != nil {
				dir.LogThread.Print(err)
			} else {
				log.Print(err)
			}
			return nil, err
		}
		if err == sql.ErrNoRows {
			filesHash[file] = ""
		} else {
			filesHash[file] = hash
		}
	}
	return filesHash, nil
}

func SetHashFromRepository(dir *models.Directory) error {
	for file, hash := range dir.FileHash {
		isExist := false
		err := Connect.QueryRowx(`SELECT EXISTS(select path from hashes where path=$1)`, file).Scan(&isExist)
		if err != nil {
			if dir.LogThread != nil {
				dir.LogThread.Print(err)
			} else {
				log.Print(err)
			}
			return err
		}
		if isExist {
			_, err := Connect.Queryx(`UPDATE hashes set path=$1, hash=$2 WHERE path=$1`, file, hash)
			if err != nil {
				if dir.LogThread != nil {
					dir.LogThread.Print(err)
				} else {
					log.Print(err)
				}
				return err
			}
		} else {
			_, err := Connect.Queryx(`INSERT INTO hashes values($1,$2)`, file, hash)
			if err != nil {
				if dir.LogThread != nil {
					dir.LogThread.Print(err)
				} else {
					log.Print(err)
				}
				return err
			}
		}
	}
	return nil
}
