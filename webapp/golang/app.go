package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func GetRoot(w http.ResponseWriter, r *http.Request) {}

func GetPost(w http.ResponseWriter, r *http.Request) {}

func PostPost(w http.ResponseWriter, r *http.Request) {}

func GetArticle(w http.ResponseWriter, r *http.Request) {}

func PostComment(w http.ResponseWriter, r *http.Request) {}

func main() {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(GetRoot)
	r.Methods("GET").Path("/post").HandlerFunc(GetPost)
	r.Methods("POST").Path("/post").HandlerFunc(PostPost)
	r.Methods("GET").Path("/article/{articleid}").HandlerFunc(GetArticle)
	r.Methods("POST").Path("/comment/{articleid}").HandlerFunc(PostComment)
	log.Fatal(http.ListenAndServe(":5000", r))
}
