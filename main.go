package main

import (
	"flag"
	"github.com/gorilla/mux"
	"net/http"
	"math/rand"
	"time"
	"log"
	"encoding/json"
	"errors"
)

const baseURL = "http://media.oboobs.ru"
var boobSource = []string{"http://api.oboobs.ru/noise/1", "http://api.oboobs.ru/boobs/0/1/random"}

type ARGs struct {
	Username string
	Password string
	ListenHostname string
}

func getBoobsLink() (string, error) {
	source := boobSource[rand.Int31n(int32(len(boobSource)))]
	rq, err := http.NewRequest("GET", source, nil)
	if err != nil {
		return "", err
	}

	hCl := http.Client{Timeout: time.Second * 5}
	rsp, err := hCl.Do(rq)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()

	var m []map[string]interface{}
	if err := json.NewDecoder(rsp.Body).Decode(&m); err != nil {
		return "", err
	}

	if preview, ok := m[0]["preview"]; ok {
		if pStr, ok := preview.(string); ok {
			return baseURL + "/" + pStr, nil
		}
	}

	return "", errors.New("preview not found in response")
}

type dfResponse struct {
	Speech string `json:"speech"`
}

func main() {
	var args ARGs
	flag.StringVar(&args.Username, "u", "", "basic auth user")
	flag.StringVar(&args.Password, "p", "", "basic auth password")
	flag.StringVar(&args.ListenHostname, "h", ":8081", "listen hostname")
	flag.Parse()
	rand.Seed(time.Now().Unix())

	r := mux.NewRouter()

	r.Methods("POST").Path("/boobs/v1/random").HandlerFunc(func (w http.ResponseWriter, rq *http.Request){
		u, p, ok := rq.BasicAuth()
		if !ok || u != args.Username || p != args.Password {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		boobLink, err := getBoobsLink()
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(dfResponse{
			Speech: boobLink,
		})
	})

	if err := http.ListenAndServe(args.ListenHostname, r); err != nil {
		log.Fatalln(err)
	}
}