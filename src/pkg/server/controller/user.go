package controller

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"net/http"
	"pkg/logger"
	"pkg/models"

	"pkg/server/controller/form"
	"pkg/server/controller/response"
	"pkg/services"
	"pkg/services/mapper"
	"pkg/services/printer"
)

func GinDoUserSignManually( c echo.Context) error {

	signForm := new(form.SignForm)
	if err := c.Bind(signForm); err != nil {
		return c.JSON(http.StatusBadRequest, response.RenderError(response.INVALID_PARAM, err))
	}

	queryType := mapper.UserQueryType(signForm.SignType)
	userInfo, err := services.QueryUserInfo(queryType, signForm.Keyword)
	if err != nil {
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, err))
	}

	if signed := services.IsUserSigned(userInfo); !signed {
		_ = services.DoSIgnForUser(mapper.QueryTypeWorkNumber, userInfo.WorkNumber)
		printer.PrintUserLabel(userInfo)
	}
	if queryType == mapper.QueryTypeWorkNumber {
		if webSocketConned.Load() == true {
			userInfoChan <- *userInfo
		}
	}

	return c.JSON(http.StatusOK, response.RenderSuccess(userInfo))
}

func GinJustPrintUserLabel(c echo.Context) error{
	workNum := c.FormValue("num")
	if len(workNum) == 0 {
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, errors.New("no args")))
	}
	userInfo, err := services.QueryUserInfo(mapper.QueryTypeWorkNumber, workNum)
	if err != nil {
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, err))
	}
	printer.PrintUserLabel(userInfo)
	return c.JSON(http.StatusOK, response.RenderSuccess(""))
}

func GinFaceSign(c echo.Context) error {
	key := c.FormValue("key")
	if len(key) == 0 {
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, errors.New("no args")))
	}
	userInfo, err := services.QueryUserInfo(mapper.QueryTypeWorkNumber, key)

	if err != nil {
		return c.JSON(http.StatusOK, response.RenderError(response.INVALID_PARAM, err))
	}

	userInfoChan <- *userInfo

	printer.PrintUserLabel(userInfo)
	return c.JSON(http.StatusOK, response.RenderSuccess(userInfo))
}

func Import(c echo.Context) error  {
	data := c.FormValue("data")
	var dataMap []map[string]string

	var user models.UserInfo

	json.Unmarshal([]byte(data),&dataMap)

	for _, one := range dataMap {
		user.Name = one["name"]
		if one["signed"] == "true" {
			user.Signed = true
		} else {
			user.Signed = false
		}
		user.QrCode = one["qr_code"]
		user.SeatArea = one["seat_area"]
		user.DomainName = one["domain_name"]
		user.WorkNumber = one["work_number"]
		user.Department = one["department"]
		err := services.InsertData(&user)
		if err != nil {
			logger.Infof("insert error %s", err)
		}
	}
	return c.JSON(http.StatusOK, response.RenderSuccess(""))
}

func All(c echo.Context) error  {

	users, ok :=services.GetAll()

	if ok {
		return c.JSON(http.StatusOK, response.RenderSuccess(users))
	}
	return c.JSON(http.StatusOK, response.RenderSuccess(""))
}

