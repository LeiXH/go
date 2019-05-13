package models

type UserInfo struct {
	Name       string        `db:"name" json:"name"`
	Code       string         `db:"code" json:"code"`
	Company string           `db:"company" json:"company"`
	Telephone string         `db:"telephone" json:"telephone"`
	Degree   int          `db:"degree" json:"degree"`
	Signed     int           `db:"signed" json:"signed"`
	Mark     string           `db:"mark" json:"mark"`
}
