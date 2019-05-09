package mapper

import (
	"fmt"

	"pkg/database"
	"pkg/models"
)

const (
	getUserInfoSqlTemplate      = `select * from briefing where %s = ?`
	countUserInfoSqlTemplate    = `select count(*) from briefing where %s = ?`
	markUserAsSignedSqlTemplate = `update briefing set signed = 1 where %s = ?`
	countUserSignSqlTemplate    = `select count(*) from briefing where signed = 1 and %s = ?`
	queryUsersOffset            = `select * from briefing limit ? offset ?`
	queryTelephoneIsExist       = `select count(*) from briefing where telephone = ? `
)

type UserQueryType int

const (
	_ UserQueryType = iota
	QueryTypeWorkNumber
	QueryTypeQrCode
	QueryTypeTelephone
)

var (
	userQueryTypeMap = map[UserQueryType]string{
		QueryTypeWorkNumber: "work_number",
		QueryTypeQrCode:     "code",
		QueryTypeTelephone : "telephone",
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
	querys := fmt.Sprintf(getUserInfoSqlTemplate, _type)
	err := DB.QueryRow(querys, keyword).Scan(&id, &userInfo.Name, &userInfo.Code,
		&userInfo.Company, &userInfo.Telephone, &userInfo.Degree, &userInfo.Signed)

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

func GetAllUser(nums int) ([]models.UserInfo, error) {
	var all []models.UserInfo
	DB := database.MySQL()

	var limit  int
	if nums == -1 {
		limit = 0
	} else {
		limit = 100
	}

	rows, err :=DB.Query(queryUsersOffset, limit, nums * limit )

	if err != nil {
		return all, err
	}

	for rows.Next() {
		var one models.UserInfo
		var id  int
		rows.Scan(&id, &one.Name, &one.Code, &one.Company, &one.Telephone, &one.Degree, &one.Signed)

		all = append(all, one)
	}
	return all, nil
}

func InsertData(info *models.UserInfo) error {
	DB := database.MySQL()
	smt, err := DB.Prepare(`insert into  briefing (name, code, company, telephone, degree, signed) values (?, ? , ? , ? , ? , ?)`)

	res, err :=smt.Exec(info.Name, info.Code, info.Company, info.Telephone, info.Degree, info.Signed)

	_, err = res.LastInsertId()

	if err != nil {
		return  err
	}
	return nil
}

func GetUserByTelephone(telephone string) (int ,  error)  {

	var c int

	DB := database.MySQL()

	err :=DB.QueryRow(queryTelephoneIsExist, telephone).Scan(&c)

	if err != nil {
		return -1, err
	}
	return c, nil
}
