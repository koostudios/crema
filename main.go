package main

import (
	"bytes"
	"encoding/json"
	"github.com/badoux/checkmail"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Config struct {
	Protocol string
	ApiKey   string
	Url      string
}

type Message struct {
	Status string `json:"status"`
	Body   string `json:"message"`
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		msg := Message{
			Status: "ok",
			Body:   "Hello, World!",
		}
		sendJson(w, msg)
	}).Methods("GET")

	r.HandleFunc("/{email}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := vars["email"]

		err := checkmail.ValidateFormat(email)
		if err != nil {
			msg := Message{
				Status: "error",
				Body:   "Email in incorrect format.",
			}
			sendJson(w, msg)
			return
		}

		msg := Message{
			Status: "ok",
			Body:   "Thanks for using our service, " + email + ". Please make sure your form has the method=POST attribute",
		}
		sendJson(w, msg)
	}).Methods("GET")

	r.HandleFunc("/{email}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		email := vars["email"]

		err := checkmail.ValidateFormat(email)
		if err != nil {
			msg := Message{
				Status: "error",
				Body:   "Email in incorrect format.",
			}
			sendJson(w, msg)
			return
		}

		r.ParseForm()
		var doc bytes.Buffer
		t, _ := template.ParseFiles("email.gtpl")
		t.Execute(&doc, r.Form)
		html := doc.String()

		file, err := ioutil.ReadFile("./config.json")
		if err != nil {
			msg := Message{
				Status: "error",
				Body:   "Could not find or read config.json.",
			}
			sendJson(w, msg)
			return
		}

		// Set up decode
		dec := json.NewDecoder(strings.NewReader(string(file)))
		config := Config{}

		// Decode json file contents
		if err = dec.Decode(&config); err != nil {
			msg := Message{
				Status: "error",
				Body:   "Could not decode config.json",
			}
			sendJson(w, msg)
			return
		}

		resp, err := http.PostForm(config.Protocol+config.ApiKey+"@"+config.Url, url.Values{
			"from":    {"Crema Forms <crema@koostudios.com>"},
			"to":      {email},
			"subject": {"New submission"},
			"html":    {html},
		})
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			msg := Message{
				Status: "error",
				Body:   "There was an error sending the email.",
			}
			sendJson(w, msg)
			return
		}

		msg := Message{
			Status: "ok",
			Body:   string(body),
		}

		sendJson(w, msg)
	})

	log.Fatal(http.ListenAndServe(":1337", r))
}

func sendJson(w http.ResponseWriter, msg Message) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}
