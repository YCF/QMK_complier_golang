package main

import (
	"Goose/conf"
	"Goose/router"
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/echo-contrib/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/syntaqx/renderer"
	"github.com/unrolled/render"
)

// Response 返回的json结构
//http://stackoverflow.com/questions/942951/rest-api-error-return-good-practices
type Response struct {
	Version string      `json:"version"` //api version
	Status  string      `json:"status"`  // http always =200,the real err will be here
	Msg     string      `json:"msg"`     // for human
	Data    interface{} `json:"data"`
}

func main() {

	qmkRoot := conf.GetOption("QMK", "root")
	db, err := sql.Open("sqlite3", conf.GetOption("system", "db"))
	if err != nil {
		fmt.Println(err)
	}
	e := echo.New()
	// Enable debug logging
	e.Debug = true

	// Keeps the DefaultFuncs provided by renderer
	funcs := []template.FuncMap{renderer.DefaultFuncs}

	// Create an instance of unrolled/render with app-specific configurations.
	r := render.New(render.Options{
		//Layout:        "layout",
		Directory:     "templates",
		Extensions:    []string{".html"},
		IsDevelopment: e.Debug,
		Funcs:         funcs,
	})

	// Wrap the render instance with the Renderer compliant interface.
	e.Renderer = renderer.Wrap(r)

	// Attach middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "static")                //static file in /static
	e.File("/favicon.ico", "static/favicon.ico") //Serve favicon.ico

	store := sessions.NewCookieStore([]byte("secret"))
	if err != nil {
		panic(err)
	}
	e.Use(sessions.Sessions("echosession", store))

	// Let's give it a go!
	e.GET("/", func(c echo.Context) error {
		session := sessions.Default(c)
		var value string
		val := session.Get("name")
		if val == nil {
			c.Redirect(301, "/login")
		} else {
			c.Redirect(301, "/mitosis-plus")
		}

		return c.JSON(200, value)
	})
	e.GET("/api/logout", func(c echo.Context) error {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		return c.Redirect(301, "/")
	})
	e.POST("/api/login", func(c echo.Context) error {
		name := c.FormValue("name")
		pwd := c.FormValue("pwd")
		var Password string
		err := db.QueryRow("select Password from user where Activity='1' and UserName = $1", name).Scan(&Password)
		if err != nil {
			fmt.Println(err)
		}
		// nil means it is a match
		bHashedPwd := []byte(Password)
		bPwd := []byte(pwd)
		err = bcrypt.CompareHashAndPassword(bHashedPwd, bPwd)
		// 密码正确
		if err == nil {
			session := sessions.Default(c)
			session.Set("name", name)
			session.Save()
			u := &Response{
				Version: "0.1",
				Status:  "301",
				Msg:     "Login successful.",
				Data:    "/mitosis-plus",
			}
			return c.JSON(http.StatusOK, u)
		}
		// 密码不对
		u := &Response{
			Version: "0.1",
			Status:  "401",
			Msg:     "Wrong pwd",
			Data:    "The User Name or Password is incorrect.",
		}
		return c.JSON(http.StatusOK, u)
	})
	e.POST("/api/signup", router.SignUp)
	e.GET("/login", func(c echo.Context) error {
		return c.Render(http.StatusOK, "login", "")
	})

	e.GET("/mitosis-plus", func(c echo.Context) error {
		return c.Render(http.StatusOK, "kb", "")
	})

	e.POST("/make/mitosis-plus", func(c echo.Context) error {
		layer0 := c.FormValue("layer0")
		layer1 := c.FormValue("layer1")
		layer2 := c.FormValue("layer2")
		layer3 := c.FormValue("layer3")
		layer4 := c.FormValue("layer4")
		qmkABS, _ := filepath.Abs(qmkRoot)
		//keymapHead, err := ioutil.ReadFile(qmkABS + "/keyboards/mitosis/keymaps/defult/head.c")
		f, err := os.Open(qmkABS + "/keyboards/mitosis/keymaps/default/head.c")
		if err != nil {
			u := &Response{
				Version: "0.1",
				Status:  "400",
				Msg:     "服务器内部错误",
				Data:    err,
			}
			return c.JSON(http.StatusOK, u)
		}
		keymapHead, err := ioutil.ReadAll(f)

		f, err = os.Open(qmkABS + "/keyboards/mitosis/keymaps/default/foot.c")
		if err != nil {
			u := &Response{
				Version: "0.1",
				Status:  "400",
				Msg:     "服务器内部错误",
				Data:    err,
			}
			return c.JSON(http.StatusOK, u)
		}
		keymapFoot, err := ioutil.ReadAll(f)
		f, err = os.Open(qmkABS + "/keyboards/mitosis/keymaps/default/no.c")
		if err != nil {
			u := &Response{
				Version: "0.1",
				Status:  "400",
				Msg:     "服务器内部错误",
				Data:    err,
			}
			return c.JSON(http.StatusOK, u)
		}
		keymapNo, err := ioutil.ReadAll(f)

		if layer0 == "" {
			u := &Response{
				Version: "0.1",
				Status:  "400",
				Msg:     "请至少设置一层",
				Data:    "Layer0 undefined.",
			}
			return c.JSON(http.StatusOK, u)
		}
		if layer1 == "" || layer1 == "undefined" {
			layer1 = string(keymapNo)
		}
		if layer2 == "" || layer2 == "undefined" {
			layer2 = string(keymapNo)
		}
		if layer3 == "" || layer3 == "undefined" {
			layer4 = string(keymapNo)
		}
		if layer4 == "" || layer4 == "undefined" {
			layer4 = string(keymapNo)
		}
		layer0 = "[0] = {\n" + layer0 + "\n},\n"
		layer1 = "[1] = {\n" + layer1 + "\n},\n"
		layer2 = "[2] = {\n" + layer2 + "\n},\n"
		layer3 = "[3] = {\n" + layer3 + "\n},\n"
		layer4 = "[4] = {\n" + layer4 + "\n},\n"
		u := &Response{
			Version: "0.1",
			Status:  "200",
			Msg:     layer1,
			Data:    keymapHead,
		}
		buf := bytes.NewBuffer(keymapHead)
		buf.Write([]byte(layer0))
		buf.Write([]byte(layer1))
		buf.Write([]byte(layer2))
		buf.Write([]byte(layer3))
		buf.Write([]byte(layer4))
		buf.Write([]byte(keymapFoot))
		t := fmt.Sprint(time.Now().UnixNano())
		tmpPath := qmkABS + "/keyboards/mitosis/keymaps/" + t
		fmt.Println(tmpPath)
		err = os.MkdirAll(tmpPath, 0711)
		err = ioutil.WriteFile(tmpPath+"/keymap.c", buf.Bytes(), 0644)
		defer f.Close()

		cmd := exec.Command("/bin/bash", "-c", `cd '`+qmkABS+`';make mitosis-`+t)
		//创建获取命令输出管道
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Printf("Error:can not obtain stdout pipe for command:%s\n", err)
			return err
		}

		//执行命令
		if err := cmd.Start(); err != nil {
			fmt.Println("Error:The command is err,", err)
			return err
		}

		//读取所有输出
		bytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			fmt.Println("ReadAll Stdout:", err.Error())
			return err
		}

		if err := cmd.Wait(); err != nil {
			fmt.Println("wait:", err.Error())
			return err
		}
		//移动编译文件
		fmt.Printf("stdout:\n\n %s", bytes)
		err = os.Rename(qmkABS+"/mitosis_"+t+".hex", "static/hexfiles/mitosis_"+t+".hex")
		u = &Response{
			Version: "0.1",
			Status:  "200",
			Msg:     "make",
			Data:    bytes,
		}
		err = os.RemoveAll(tmpPath)
		if err != nil {
			return c.JSON(http.StatusOK, u)
		}
		return c.Redirect(301, "/static/hexfiles/mitosis_"+t+".hex")
	})

	// RESTful api
	e.GET("/api", func(c echo.Context) error {
		u := &Response{
			Version: "0.1",
			Status:  "200",
			Msg:     "ok",
			Data:    "Oh hello, sir! -J.A.R.V.I.S",
		}
		return c.JSON(http.StatusOK, u)
	})

	// port := cfg.Section("system").Key("port").String()
	e.Logger.Fatal(e.Start(conf.GetOption("system", "port")))
}
