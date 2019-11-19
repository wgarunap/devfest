package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-errors/errors"
	"github.com/gorilla/mux"
	"github.com/pickme-go/log"
	"io/ioutil"
	"net/http"
	"strconv"
)

var person map[int]Person

type Person struct {
	ID          int    `json:"id,omitempty"`
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
	ID int `json:"id,omitempty"`
}

type GetResponse struct {
	Firstname   string `json:"firstname,omitempty"`
	Lastname    string `json:"lastname,omitempty"`
	Contactinfo `json:"contactinfo,omitempty"`
}

type HandlerPost struct {
}

type HandlerGet struct {
}

func (HandlerPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	p.ID = len(person) + 1

	person[p.ID] = p

	w.WriteHeader(http.StatusOK)

	res := PostResponse{ID: p.ID}
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)

}

func (HandlerGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var bid int

	params := mux.Vars(r)
	if !(len(params) > 0) {
		err = errors.New("book id missing in request")
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = fmt.Fprint(w, err)
		return
	}

	bidStr, ok := params["id"]
	if !ok {
		err = errors.New("book id missing in request")
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

	person, ok := person[bid]
	if !ok {
		log.Info("no records for given book id")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "no records for given book id")
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
