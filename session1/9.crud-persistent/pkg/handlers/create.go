package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pickme-go/log"
	"github.com/wgarunap/devfest/session1/8.crud-persistent/pkg/adapters"
	"github.com/wgarunap/devfest/session1/8.crud-persistent/pkg/models"
	"github.com/wgarunap/devfest/session1/8.crud-persistent/pkg/repositories"
)

var PersonMap map[int]models.Person

type PostResponse struct {
	ID int `json:"id,omitempty"`
}

type GetResponse struct {
	Firstname          string `json:"firstname,omitempty"`
	Lastname           string `json:"lastname,omitempty"`
	models.Contactinfo `json:"contactinfo,omitempty"`
}

type GetAllResponse struct {
	Data map[int]models.Person `json:"data"`
}

type HandlerPost struct {
}

type HandlerGet struct {
}

type HandlerGetAll struct {
}

type HandlerPut struct {
}

type HandlerDelete struct {
}

var personRepository = repositories.NewPersonRepository(adapters.NewDbConnection())

func (HandlerPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// p.ID = len(PersonMap) + 1

	// PersonMap[p.ID] = p

	fmt.Println(p)
	id, err := personRepository.AddPerson(p)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)

	res := PostResponse{ID: id}
	b, _ := json.Marshal(res)
	_, _ = w.Write(b)

}

func (HandlerGet) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	recordExist, person, err := personRepository.GetPersonByID(bid)

	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	if !recordExist {
		log.Error(err)
		w.WriteHeader(http.StatusNotFound)
		// _, _ = fmt.Fprint(w, err)
		return
	}

	res := GetResponse{}
	res.Phone = person.Phone
	res.City = person.City
	res.Firstname = person.Firstname
	res.Lastname = person.Lastname
	res.Zipcode = person.Zipcode
	res.Contactinfo = person.Contactinfo
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	b, _ := json.Marshal(res)
	_, _ = w.Write(b)

}

func (HandlerGetAll) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// res := GetAllResponse{}

	persons, err := personRepository.GetAll()
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	// res.Data = PersonMap
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(persons)
	_, _ = w.Write(b)

}

func (HandlerPut) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// _, ok = PersonMap[bid]
	// if !ok {
	// 	log.Info("no records for given  id to update")
	// 	w.WriteHeader(http.StatusOK)
	// 	_, _ = fmt.Fprint(w, "no records for given  id to update")
	// 	return
	// }

	// PersonMap[bid] = p

	success, err := personRepository.UpdatePerson(bid, p)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	if !success {
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "success updating info for id: %v", bid)

}

func (HandlerDelete) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	_, err = personRepository.DeletePerson(bid)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprint(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
