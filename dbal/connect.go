package dbal

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/team-gary/x/sqlcon"
)

// Connect is a wrapper for connecting to different SQL drivers.
func Connect(db string, logger logrus.FieldLogger, memf func() error, sqlf func(db *sqlx.DB) error) error {
	if db == "memory" {
		return memf()
	} else if db == "" {
		return errors.New("No database DSN provided")
	}

	scheme := sqlcon.GetDriverName(db)
	switch scheme {
	case "postgres":
		fallthrough
	case "cockroach":
		fallthrough
	case "mysql":
		c, err := sqlcon.NewSQLConnection(db, logger)
		if err != nil {
			return errors.WithStack(err)
		}

		cdb, err := c.GetDatabaseRetry(time.Second*15, time.Minute*5)
		if err != nil {
			return errors.WithStack(err)
		}

		return sqlf(cdb)
	}

	return errors.Errorf("The provided database DSN %s can not be handled", db)
}
