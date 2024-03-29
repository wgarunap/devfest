package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/wgarunap/devfest/session3/11.crud-with-pprof/metrics"
	"github.com/wgarunap/devfest/session3/11.crud-with-pprof/pkg/adapters"
	"github.com/wgarunap/devfest/session3/11.crud-with-pprof/pkg/models"
	"github.com/wgarunap/devfest/session3/11.crud-with-pprof/pkg/repositories"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pickme-go/log"
)

var counter int64
var PersonMap map[int64]models.Person

type PostResponse struct {
	ID int64 `json:"id,omitempty"`
}

type GetResponse struct {
	ID                 int    `json:"id,omitempty"`
	Firstname          string `json:"firstname,omitempty"`
	Lastname           string `json:"lastname,omitempty"`
	models.ContactInfo `json:"contactinfo,omitempty"`
}

type GetAllResponse struct {
	Data map[int64]models.Person `json:"data"`
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

var personRepository = repositories.NewPersonRepository(adapters.NewDbConnection())

func (hp HandlerPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer hp.MetricsCounter.CountLatency(time.Now(), []string{`post`})

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	var p models.Person

	err = json.Unmarshal(data, &p)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	personID, err := personRepository.Add(p)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Insert success. Response with HTTP status 201 (created)
	w.WriteHeader(http.StatusCreated)
	pr := PostResponse{
		ID: int64(personID),
	}

	responseData, err := json.Marshal(pr)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(responseData)
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

	recordExist, person, err := personRepository.GetByID(bid)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	if !recordExist {
		log.Info("no records for given book id")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "no records for given book id")
		return
	}

	res := GetResponse{}
	res.ID = person.ID
	res.Phone = person.Phone
	res.City = person.City
	res.Firstname = person.Firstname
	res.Lastname = person.Lastname
	res.AreaCode = person.AreaCode
	res.ContactInfo = person.ContactInfo

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)

}

func (hga HandlerGetAll) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer hga.MetricsCounter.CountLatency(time.Now(), []string{`get_all`})

	respData := []GetResponse{}

	persons, err := personRepository.GetAll()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	for _, person := range persons {
		personResp := GetResponse{}
		personResp.ID = person.ID
		personResp.Firstname = person.Firstname
		personResp.Lastname = person.Lastname
		personResp.ContactInfo = person.ContactInfo

		respData = append(respData, personResp)
	}

	jsonResp, err := json.Marshal(respData)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)

}

func (hp HandlerPut) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer hp.MetricsCounter.CountLatency(time.Now(), []string{`put`})

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

	var p models.Person

	err = json.Unmarshal(data, &p)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	ok, err = personRepository.Update(bid, p)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	if !ok {
		log.Info("nothing updates")
		w.WriteHeader(http.StatusNotModified)
		_, _ = fmt.Fprint(w, "nothing updated")
		return
	}

	// PersonMap[int64(bid)] = p
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

	ok, err = personRepository.Delete(bid)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	if !ok {
		log.Info("nothing updates")
		w.WriteHeader(http.StatusNotModified)
		_, _ = fmt.Fprint(w, "nothing updated")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
