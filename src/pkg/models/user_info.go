package models

import "pkg/models/types"

type UserInfo struct {
	WorkNumber string        `db:"work_number" json:"work_number"`
	DomainName string        `db:"domain_name" json:"domain_name"`
	Department string        `db:"department" json:"department"`
	SeatArea   string        `db:"seat_area" json:"seat_area"`
	Name       string        `db:"name" json:"name"`
	QrCode     string        `db:"qrcode" json:"qr_code"`
	Signed     types.SQLBool `db:"signed" json:"signed"`
}
