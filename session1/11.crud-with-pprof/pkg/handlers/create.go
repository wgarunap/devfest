package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/pickme-go/log"
	"github.com/wgarunap/devfest/session1/11.crud-with-pprof/metrics"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var counter int64
var PersonMap map[int64]Person

type Person struct {
	ID          int64    `json:"id,omitempty"`
	Firstname   string `json:"firstname,omitempty"`
	Lastname    string `json:"lastname,omitempty"`
	Contactinfo `json:"contactinfo,omitempty"`
}
type Contactinfo struct {
	City    string `json:"city,omitempty"`
	Zipcode string `json:"zipcode,omitempty"`
	Phone   string `json:"phone,omitempty"`
}

type PostResponse struct {
	ID int64 `json:"id,omitempty"`
}

type GetResponse struct {
	Firstname   string `json:"firstname,omitempty"`
	Lastname    string `json:"lastname,omitempty"`
	Contactinfo `json:"contactinfo,omitempty"`
}

type GetAllResponse struct {
	Data map[int64]Person `json:"data"`
}

type HandlerPost struct {
	MetricsCounter metrics.Metricer
}

type HandlerGet struct {
	MetricsCounter metrics.Metricer
}

type HandlerGetAll struct {
	MetricsCounter metrics.Metricer
}

type HandlerPut struct {
	MetricsCounter metrics.Metricer
}

type HandlerDelete struct {
	MetricsCounter metrics.Metricer
}

func (hp HandlerPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer hp.MetricsCounter.CountLatency(time.Now(), []string{`post`})

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	var p Person

	err = json.Unmarshal(data, &p)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	p.ID = atomic.AddInt64(&counter, 1)

	PersonMap[p.ID] = p

	w.WriteHeader(http.StatusOK)

	res := PostResponse{ID: p.ID}
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)

}

func (hg HandlerGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer hg.MetricsCounter.CountLatency(time.Now(), []string{`get`})

	var err error
	var bid int

	params := mux.Vars(r)
	if !(len(params) > 0) {
		err = errors.New("id missing in request")
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	bidStr, ok := params["id"]
	if !ok {
		err = errors.New("id missing in request")
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	bid, err = strconv.Atoi(bidStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	person, ok := PersonMap[int64(bid)]
	if !ok {
		log.Info("no records for given  id")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "no records for given  id")
		return
	}

	res := GetResponse{}
	res.Phone = person.Phone
	res.City = person.City
	res.Firstname = person.Firstname
	res.Lastname = person.Lastname
	res.Zipcode = person.Zipcode
	res.Contactinfo = person.Contactinfo
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)

}

func (hga HandlerGetAll) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer hga.MetricsCounter.CountLatency(time.Now(), []string{`get_all`})

	res := GetAllResponse{}
	res.Data = PersonMap
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)

}

func (hpu HandlerPut) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer hpu.MetricsCounter.CountLatency(time.Now(), []string{`put`})

	var err error
	var bid int

	params := mux.Vars(r)
	if !(len(params) > 0) {
		err = errors.New("id missing in request")
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	bidStr, ok := params["id"]
	if !ok {
		err = errors.New("id missing in request")
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	bid, err = strconv.Atoi(bidStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	var p Person

	err = json.Unmarshal(data, &p)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	_, ok = PersonMap[int64(bid)]
	if !ok {
		log.Info("no records for given  id to update")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "no records for given  id to update")
		return
	}

	PersonMap[int64(bid)] = p
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "success updating info for id: %v", bid)

}

func (hd HandlerDelete) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer hd.MetricsCounter.CountLatency(time.Now(), []string{`delete`})

	var err error
	var bid int

	params := mux.Vars(r)
	if !(len(params) > 0) {
		err = errors.New("id missing in request")
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	bidStr, ok := params["id"]
	if !ok {
		err = errors.New("id missing in request")
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	bid, err = strconv.Atoi(bidStr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	_, ok = PersonMap[int64(bid)]
	if !ok {
		log.Info("nothing to delete")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "nothing to delete")
		return
	}

	delete(PersonMap, int64(bid))
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "success deleting info for id: %v", bid)

}
