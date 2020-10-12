package gotask

import (
	"database/sql"
	"fmt"
	"time"
)

// select driver
// DB

// DB lock

// create task

// manage tasks with DB

// integrate UI

// notifier

type Locker interface {
	Lock()
	Unlock()
}

const tableName = "gotask.mutex"

type Drive struct {
	db *sql.DB
}

func (d *Drive) Lock(name, value string, expiry time.Duration) (bool, error) {
	var tempValue string
	var tempExpiry time.Duration
	// create if not exist
	// overwrite if the lock expire
	err := d.db.QueryRow(fmt.Sprintf(`SELECT value, expiry FROM %s
		WHERE name=$1`)).Scan(&tempValue, &tempExpiry)
	if err != nil {
		if err == sql.ErrNoRows { // not exist

		}
	}
}

func (d *Drive) Unlock(name, value string) {

}
