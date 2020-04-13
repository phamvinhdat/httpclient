package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/phamvinhdat/httpclient"
	"github.com/phamvinhdat/httpclient/body"
	"github.com/phamvinhdat/httpclient/gosender"
	"github.com/phamvinhdat/httpclient/hook"
)

type Credential struct {
	Username string
	Password string
}

func main() {
	go newFakeServer()

	result := Credential{}
	client := httpclient.NewClient(
		httpclient.WithSender(gosender.New(gosender.WithTimeout(time.Second * 5))),
	)
	statusCode, err := client.Post(context.Background(), "http://localhost:8080/login",
		httpclient.WithBodyProvider(body.NewJson(Credential{
			Username: "admin",
			Password: "admin",
		})),
		httpclient.WithHookFn(hook.UnmarshalResponse(&result)),
		httpclient.WithHookFn(hook.Log()),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("http code: ", statusCode)
	fmt.Println("result: ", result)
}

func newFakeServer() {
	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		cre := Credential{}
		err := c.ShouldBind(&cre)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"err": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, cre)
	})
	_ = r.Run()
}
