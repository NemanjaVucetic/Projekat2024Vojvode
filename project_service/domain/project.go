package domain

import (
	"encoding/json"
	"io"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name     string             `bson:"name" json:"name"`
	Surname  string             `bson:"surname" json:"surname"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"password"`
	UserRole string             `bson:"userRole" json:"userRole"`
	Enabled  bool               `bson:"enabled" json:"enabled"`
}

type Project struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	ManagerID  primitive.ObjectID `bson:"manager" json:"manager"`
	Members    []*User            `bson:"members" json:"members"`
	Deadline   string             `bson:"deadline" json:"deadline"`
	MaxMembers int                `bson:"maxMembers" json:"maxMembers"`
	MinMembers int                `bson:"minMembers" json:"minMembers"`
}

// Methods for encoding/decoding
func (p *Project) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Project) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
