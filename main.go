package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"wapiti/manager"
	"wapiti/models"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

var (
	c   = cache.New(5*time.Minute, 10*time.Minute)
	log = logrus.New()
)

func main() {
	startServer()
}

func startServer() {
	num := runtime.NumCPU()

	runtime.GOMAXPROCS(num << 2)

	router := httprouter.New()

	router.POST("/", Excute)
	router.GET("/:id", GetResult)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Panic(err)
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
	go func() {
		m.Collector()
		c.Set(m.Uuid.String(), m.Response, cache.DefaultExpiration)
	}()

	w.Header().Set("content-type", "application/json")
	w.Header().Set("uuid", m.Uuid.String())
	w.WriteHeader(http.StatusOK)
}

func GetResult(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	key := params.ByName("id")

	w.Header().Set("content-type", "application/json")
	if value, ok := c.Get(key); !ok {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(errors.New(fmt.Sprintf("not found: %s", key)).Error()))
	} else {
		b, _ := json.Marshal(value)
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
