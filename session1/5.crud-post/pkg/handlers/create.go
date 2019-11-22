package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/pickme-go/log"
	"io/ioutil"
	"net/http"
	"sync/atomic"
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

type HandlerPost struct {
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

	p.ID = atomic.AddInt64(&counter, 1)

	PersonMap[p.ID] = p

	w.WriteHeader(http.StatusOK)

	res := PostResponse{ID: p.ID}
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)

}
