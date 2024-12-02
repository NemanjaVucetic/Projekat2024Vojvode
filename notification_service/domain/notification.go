package domain

import (
	"encoding/json"
	"io"
	"time"

	"github.com/gocql/gocql"
)

type Notification struct {
	ID        gocql.UUID
	UserID    gocql.UUID
	Message   string
	Status    string
	CreatedAt time.Time
}

type Notifications []*Notification

func (p *Notifications) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Notification) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
