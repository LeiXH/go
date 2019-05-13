package controller

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"net/http"
	"pkg/logger"
	"pkg/models"
	"strconv"

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
		_ = services.DoSIgnForUser(mapper.QueryTypeTelephone, userInfo.Telephone)
		printer.PrintUserLabel(userInfo)
	}
	if queryType == mapper.QueryTypeQrCode {
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
	userInfo, err := services.QueryUserInfo(mapper.QueryTypeTelephone, workNum)
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
		user.Signed, _ = strconv.Atoi(one["signed"])
		user.Code = one["code"]
		user.Company = one["company"]
		user.Telephone = one["telephone"]
		user.Degree, _ = strconv.Atoi(one["degree"])
		user.Mark = one["mark"]
		re :=services.IsExist(user.Telephone)
		if re == true {
			err := services.InsertData(&user)
			if err != nil {
				logger.Infof("insert error %s", err)
			}
		}
	}
	return c.JSON(http.StatusOK, response.RenderSuccess(""))
}

func All(c echo.Context) error  {

	nums, _ := strconv.Atoi(c.FormValue("num"))
	users, ok :=services.GetAll(nums)

	if ok {
		return c.JSON(http.StatusOK, response.RenderSuccess(users))
	}
	return c.JSON(http.StatusOK, response.RenderSuccess(""))
}

