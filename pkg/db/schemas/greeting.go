package schemas

import (
	"time"

	"github.com/gocql/gocql"
)

type Greeting struct {
	ID        gocql.UUID
	Message   string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
