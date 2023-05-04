package main

import (
	"encoding/json"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"plumelog/config"
	"plumelog/elastic"
	"plumelog/kafka"
	"plumelog/log"
	"plumelog/model"
	"plumelog/utils"
	"strconv"
	"time"
)

var (
	upGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	healths []model.HealthGuard
)

func main() {
	model.ConnMap = make(map[string]*websocket.Conn)
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Static("/assets", "dist/assets")
	r.StaticFile("/favicon.ico", "dist/favicon.ico")
	r.Use(static.Serve("", static.LocalFile("dist/index.html", true)))
	r.Use(utils.CorsHandler())
	r.POST("/login", login)
	r.Use(utils.Auth())
	r.GET("/ws", wsHandle)
	r.GET("/health", queryHealth)
	r.GET("/hello", hello)
	r.POST("/query", query)
	r.GET("/appNames", getAppNames)
	r.POST("/errors", queryErrors)
	r.POST("/urls", queryUrls)
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.POST("/download", downloadLog)
	r.NoRoute(func(c *gin.Context) {
		utils.WriteHtml(c)
	})

	defer r.Run(":" + strconv.Itoa(config.Conf.Gin.Port))
}

func login(c *gin.Context) {
	user := model.UserInfo{}
	if err := c.BindJSON(&user); err != nil {
		log.Error.Println("user login bind json err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	if user.Username != config.Conf.Gin.Username || user.Password != config.Conf.Gin.Password {
		log.Error.Println("username or password error")
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "username or password error"})
		return
	}
	token := utils.GenerateToken(user.Username, user.Password)
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "login success", "data": token})
}

func hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": "hello world"})
}

func wsHandle(c *gin.Context) {
	//升级get请求为webSocket协议
	query := c.Request.URL.Query()
	token := query["token"][0]
	appName := query["appName"][0]
	env := query["env"][0]
	conn, _ := createConn(c)
	returnSuccess(conn)
	key := token + "_" + appName
	group := kafka.ConsumerMap[config.Conf.GroupId]
	//ch := group.WsChMap[key]
	model.ConnMap[key] = conn
	ch := make(chan model.PlumelogInfo)
	group.WsChMap[key] = ch
	go hearBeat(conn, token, appName)
	go func() {
		for {
			select {
			case plumelogInfo := <-ch:
				if appName == plumelogInfo.AppName && (env == "" || env == plumelogInfo.Env) {
					bytes, err := json.Marshal(&plumelogInfo)
					err = conn.WriteMessage(1, bytes)
					if err != nil {
						log.Error.Println(err)
						return
					}
				}
			}
		}
	}()
}

func createConn(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	return conn, err
}

func returnSuccess(conn *websocket.Conn) {
	info := &model.PlumelogInfo{
		Seq:       1,
		LogLevel:  "INFO",
		TraceId:   "start",
		DateTime:  time.Now().Format("2006-01-02 15:04:05"),
		ClassName: "ws",
		Method:    "start",
	}
	info.Content = "====================================================滚动日志连接成功!====================================================="
	bytes, _ := json.Marshal(&info)
	conn.WriteMessage(1, bytes)
}

// heartBeat 心跳检测
func hearBeat(conn *websocket.Conn, token, appName string) {
	for {
		key := token + "_" + appName
		_, _, err := conn.ReadMessage()
		if err != nil {
			conn.Close()
			log.Error.Println(token, "is disconnect:", err)
			delete(model.ConnMap, key)
			return
		}
	}
}

func query(c *gin.Context) {
	req := model.PlumelogInfoPageReq{}
	if err := c.BindJSON(&req); err != nil {
		log.Error.Println("plume log page bind json err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	page, err := elastic.QueryPage(&req)
	if err != nil {
		log.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": page})
}

func getAppNames(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": elastic.GetAppNames()})
}

func queryErrors(c *gin.Context) {
	req := model.ErrorsReq{}
	if err := c.BindJSON(&req); err != nil {
		log.Error.Println("plume log errors bind json err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	errors, err := elastic.QueryErrors(&req)
	if err != nil {
		log.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": errors})
}

func downloadLog(c *gin.Context) {
	req := model.PlumelogInfoPageReq{}
	if err := c.BindJSON(&req); err != nil {
		log.Error.Println("plume log page bind json err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	data, err := elastic.DownloadLog(&req)
	if err != nil {
		log.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	c.Header("Transfer-Encoding", "chunked")
	c.Writer.Header().Add("Content-Type", "text/plain;charset=utf-8")
	for _, datum := range data {
		c.Writer.Write([]byte(datum))
	}
	c.Writer.Flush()
}

func queryHealth(c *gin.Context) {
	conn, _ := createConn(c)
	writeHealth(conn)
	query := c.Request.URL.Query()
	token := query["token"][0]
	group := kafka.ConsumerMap[config.Conf.GroupId]
	ch := make(chan []byte)
	group.HealthMap[token] = ch
	go func() {
		for {
			select {
			case bytes := <-ch:
				healthGuard := model.HealthGuard{}
				json.Unmarshal(bytes, &healthGuard)
				if healthGuard.Status == "UP" {
					healths = append(healths, healthGuard)
				} else {
					healths = delSlice(healths, healthGuard)
				}
				marshal, _ := json.Marshal(&healths)
				log.Info.Println(string(marshal))
				err := conn.WriteMessage(1, marshal)
				if err != nil {
					log.Error.Println(err)
				}
			}
		}
	}()
}

func writeHealth(conn *websocket.Conn) error {
	healths = make([]model.HealthGuard, 0)
	health, err := elastic.QueryHealth()
	if err != nil {
		conn.WriteMessage(1, []byte("查询es失败!err:"+err.Error()))
	}
	marshal, _ := json.Marshal(&health)
	err = conn.WriteMessage(1, marshal)
	healths = health
	return err
}

func delSlice(a []model.HealthGuard, healthGuard model.HealthGuard) []model.HealthGuard {
	tgt := a[:0]
	for _, health := range a {
		if healthGuard.AppName == health.AppName && healthGuard.Env == health.Env && healthGuard.Ip == health.Ip {
			continue
		}
		tgt = append(tgt, health)
	}
	return tgt
}

func queryUrls(c *gin.Context) {
	req := model.UrlStatReq{}
	if err := c.BindJSON(&req); err != nil {
		log.Error.Println("plume log errors bind json err:", err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	urls, err := elastic.QueryUrlCount(&req)
	if err != nil {
		log.Error.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": urls})
}
