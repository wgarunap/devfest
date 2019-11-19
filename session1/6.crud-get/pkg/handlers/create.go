package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pickme-go/log"
	"io/ioutil"
	"net/http"
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

type Response struct {
	ID int `json:"id,omitempty"`
}

type Handler struct {
}

func (Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	res := Response{ID: p.ID}
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)

}
