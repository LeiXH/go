package mapper

import (
	"fmt"

	"pkg/database"
	"pkg/models"
)

const (
	getUserInfoSqlTemplate      = `select * from staff where %s = ?`
	countUserInfoSqlTemplate    = `select count(*) from staff where %s = ?`
	markUserAsSignedSqlTemplate = `update staff set signed = 1 where %s = ?`
	countUserSignSqlTemplate    = `select count(*) from staff where signed = 1 and %s = ?`
)

type UserQueryType int

const (
	_ UserQueryType = iota
	QueryTypeQrCode
	QueryTypeWorkNumber
)

var (
	userQueryTypeMap = map[UserQueryType]string{
		QueryTypeWorkNumber: "work_number",
		QueryTypeQrCode:     "qrcode",
	}
)

func GetUserInfo(queryType UserQueryType, keyword interface{}) (*models.UserInfo, error) {

	_type, ok := userQueryTypeMap[queryType]
	if !ok {
		return nil, fmt.Errorf("user query type %d not valid", queryType)
	}

	DB := database.MySQL()

	userInfo := models.UserInfo{}
	var id  int
	var sign int
	querys := fmt.Sprintf(getUserInfoSqlTemplate, _type)
	err := DB.QueryRow(querys, keyword).Scan(&id, &userInfo.WorkNumber, &userInfo.DomainName,
		&userInfo.Department, &userInfo.SeatArea, &userInfo.Name, &userInfo.QrCode, &sign)

	if sign != 0 {
		userInfo.Signed = true
	}
	if err != nil {
		return nil, err
	}
	return &userInfo, nil
}

func SignForUser(queryType UserQueryType, keyword interface{}) error {
	_type, ok := userQueryTypeMap[queryType]
	if !ok {
		return fmt.Errorf("user query type %d not valid", queryType)
	}
	DB := database.MySQL()
	_, err := DB.Exec(fmt.Sprintf(markUserAsSignedSqlTemplate, _type), keyword)
	return err
}

func countUserInfoRecord(queryType UserQueryType, keyword interface{}) (int, error) {

	var c int

	_type, ok := userQueryTypeMap[queryType]
	if !ok {
		return -1, fmt.Errorf("user query type %d not valid", queryType)
	}

	DB := database.MySQL()

	err := DB.QueryRow(fmt.Sprintf(countUserInfoSqlTemplate, _type)).Scan(&c)
	if err != nil {
		return -1, err
	}
	return c, nil
}

func CountUserSignedRecord(queryType UserQueryType, keyword interface{}) (int, error) {

	var c int

	_type, ok := userQueryTypeMap[queryType]
	if !ok {
		return -1, fmt.Errorf("user query type %d not valid", queryType)
	}

	DB := database.MySQL()

	err :=DB.QueryRow(fmt.Sprintf(countUserSignSqlTemplate, _type), keyword).Scan(&c)

	//err := DB.Get(&c, fmt.Sprintf(countUserSignSqlTemplate, _type), keyword)
	if err != nil {
		return -1, err
	}
	return c, nil
}

func GetAllUser() ([]models.UserInfo, error) {
	var all []models.UserInfo
	DB := database.MySQL()
	rows, err :=DB.Query("select * from staff")

	if err != nil {
		return all, err
	}

	for rows.Next() {
		var one models.UserInfo
		var id  int
		var sign int64
		rows.Scan(&id, &one.WorkNumber, &one.DomainName, &one.Department, &one.SeatArea, &one.Name, &one.QrCode, &sign)
		if sign != 0 {
			one.Signed = true
		}

		all = append(all, one)
	}
	return all, nil
}

func InsertData(info *models.UserInfo) error {
	DB := database.MySQL()
	smt, err := DB.Prepare(`insert into  staff (work_number, domain_name, department, seat_area, name, qrcode, signed) values (?, ? , ? , ? , ? , ?, ?)`)

	res, err :=smt.Exec(info.WorkNumber, info.DomainName, info.Department, info.SeatArea, info.Name, info.QrCode, info.Signed)

	_, err = res.LastInsertId()

	if err != nil {
		return  err
	}
	return nil
}
