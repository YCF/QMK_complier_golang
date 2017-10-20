package router

import (
	"Goose/model"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

// Response 返回的json结构
//http://stackoverflow.com/questions/942951/rest-api-error-return-good-practices
type Response struct {
	Version string      `json:"version"` //api version
	Status  string      `json:"status"`  // http always =200,the real err will be here
	Msg     string      `json:"msg"`     // for human
	Data    interface{} `json:"data"`
}

// SignUp 注册
func SignUp(c echo.Context) error {
	name := c.FormValue("name")
	pwd := c.FormValue("pwd")
	dsname := c.FormValue("dsname")
	db, _ := model.NewDB()
	err := db.QueryRow("select UserName from user where  UserName = $1", name).Scan(&name)
	if err != sql.ErrNoRows {
		u := &Response{
			Version: "0.1",
			Status:  "409",
			Msg:     "User:" + c.FormValue("name") + " is existed.",
			Data:    "",
		}
		return c.JSON(http.StatusOK, u)
	}
	tx, err := db.Begin()
	stmt, err := tx.Prepare("INSERT INTO user(UserName, Password, UserDsName,Activity) values(?,?,?,?)")
	if err != nil {
		fmt.Println(err)
	}
	defer stmt.Close()
	password := []byte(pwd)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec(name, string(hashedPassword), dsname, 1)
	tx.Commit()
	u := &Response{
		Version: "0.1",
		Status:  "200",
		Msg:     "User:" + c.FormValue("name") + " added.",
		Data:    "",
	}
	return c.JSON(http.StatusOK, u)
}
