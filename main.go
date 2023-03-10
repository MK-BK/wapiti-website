package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"wapiti/manager"
	"wapiti/models"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

var (
	c   = cache.New(time.Hour, 30*time.Minute)
	log = logrus.New()
)

func main() {
	e := gin.Default()

	e.GET("/:id", GetResult)
	e.POST("/", Excute)

	if err := http.ListenAndServe(":18080", e); err != nil {
		panic(err)
	}
}

func Excute(ctx *gin.Context) {
	request := &models.Request{}

	if err := ctx.ShouldBind(request); err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}

	m := manager.NewManager(request)
	c.SetDefault(m.Uid, m.Response)

	go func() {
		start := time.Now()

		log.Infof("Job: %+v start", m.Uid)
		defer log.Infof("Job: %+v end, spend: %+v", m.Uid, time.Now().Sub(start))

		if err := m.Handler(context.Background()); err != nil {
			m.Response.Status = models.StatusFailed
			m.Response.Error = err.Error()
		} else {
			m.Response.Status = models.StatusComplete
		}
		c.SetDefault(m.Uid, m.Response)

	}()

	ctx.JSON(http.StatusOK, map[string]string{
		"uuid": m.Uid,
	})
}

func GetResult(ctx *gin.Context) {
	id := ctx.Param("id")

	value, ok := c.Get(id)
	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, errors.New(fmt.Sprintf("not found: %s", id)))
	}

	ctx.JSON(http.StatusOK, value)
}
