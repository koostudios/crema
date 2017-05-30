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
	"os"
)

// From config.json
type Config struct {
	Protocol string `json:"protocol"`
	ApiKey   string `json:"api-key"`
	Url      string `json:"url"`
}

// Send message
type Message struct {
	Status string `json:"status"`
	Body   string `json:"message"`
}

func main() {
	r := mux.NewRouter()

	// ALL /
	// Sends a Hello, Crema
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sendJson(w, r, Message{
			Status: "ok",
			Body:   "Hello, Crema!",
		})
	})

	// GET /{email}
	// Both test validates the email and reminds user to POST to /{email}
	r.HandleFunc("/{email}", func(w http.ResponseWriter, r *http.Request) {
		email := mux.Vars(r)["email"]

		// Validates email and sends an error message if it is not validated
		if err := checkmail.ValidateFormat(email); err != nil {
			sendJson(w, r, Message{
				Status: "error",
				Body:   "Email in incorrect format.",
			})
			return
		}

		// If the email is validated, a reminder for the user to POST to the account
		sendJson(w, r, Message{
			Status: "ok",
			Body:   "Thanks for using our service, " + email + ". Please make sure your form has the method=POST attribute",
		})
	}).Methods("GET")

	// POST /{email}
	// Main function of crema forms - validates the email, executes the template, sends template through email
	r.HandleFunc("/{email}", func(w http.ResponseWriter, r *http.Request) {
		email := mux.Vars(r)["email"]

		// Validates email and sens an error message if it is not validated
		if err := checkmail.ValidateFormat(email); err != nil {
			sendJson(w, r, Message{
				Status: "error",
				Body:   "Email in incorrect format.",
			})
			return
		}

		// Parse the form and execute it into the email template
		r.ParseMultipartForm((1 << 10) * 24)
		var doc bytes.Buffer
		t, err := template.ParseFiles("email.tmpl")
		if err != nil {
			sendJson(w, r, Message{
				Status: "error",
				Body:   "Could not parse template..",
			})
		}
		t.Execute(&doc, r.Form)
		html := doc.String()

		config := Config{
			Protocol: os.Getenv("PROTOCOL"),
			ApiKey:   os.Getenv("APIKEY"),
			Url:      os.Getenv("URL"),
		}

		// Post to Mailchimp
		resp, err := http.PostForm(config.Protocol+config.ApiKey+"@"+config.Url, url.Values{
			"from":    {"crema@koostudios.com"},
			"to":      {email},
			"subject": {"New submission"},
			"html":    {html},
		})
		if err != nil {
			sendJson(w, r, Message{
				Status: "error",
				Body:   "Could not send email.",
			})
		}

		// Read response from Mailchimp
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			sendJson(w, r, Message{
				Status: "error",
				Body:   "There was an error sending the email.",
			})
			return
		}

		defer resp.Body.Close()

		// Sends response from Mailchimp
		sendJson(w, r, Message{
			Status: "ok",
			Body:   string(body),
		})
	})

	// Starts the Application
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Function sendJson sets the content type and relevant headers before sending a Message
func sendJson(w http.ResponseWriter, r *http.Request, msg Message) {
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Set Access-Control-Allow-Origin
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	// Exposes the required headers
	if expose := r.URL.Query()["access-control-expose-headers"]; len(expose) != 0 {
		w.Header()["Access-Control-Expose-Headers"] = expose
		for key := range expose {
			w.Header()[expose[key]] = r.URL.Query()["__amp_source_origin"]
		}
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}
