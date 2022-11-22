package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"wapiti/manager"
	"wapiti/models"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

var (
	c   = cache.New(time.Hour, 30*time.Minute)
	log = logrus.New()
)

var timeOut = &models.TIMEOUT

func main() {
	env := os.Getenv("WAPITI_TIMEOUT")
	if env != "" {
		timeout, err := time.ParseDuration(env)
		if err != nil {
			log.Error(err)
		} else {
			timeOut = &timeout
		}
	}

	log.Warn("Env TIMEOUT:", timeOut)

	router := httprouter.New()

	router.POST("/", Excute)
	router.GET("/:id", GetResult)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func Excute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	request := &models.Request{}
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	m := manager.NewManager(request)
	c.SetDefault(m.Uid, m.Response)

	go func() {
		m.Handler()
		if m.Err != nil {
			m.Response.Status = models.StatusFailed
			m.Response.Error = m.Err.Error()
		} else {
			m.Response.Status = models.StatusComplete
		}
		c.SetDefault(m.Uid, m.Response)
	}()

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]string{
		"uuid": m.Uid,
	}
	b, _ := json.Marshal(resp)
	w.Write(b)
}

func GetResult(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("id")

	w.Header().Set("content-type", "application/json")

	value, ok := c.Get(key)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errors.New(fmt.Sprintf("not found: %s", key)).Error()))
		return
	}

	b, _ := json.Marshal(value)
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
