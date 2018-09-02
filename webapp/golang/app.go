package main

import (
	"log"
	"net/http"
	"time"
	"path/filepath"
	"fmt"
	"html/template"
	"strconv"
	"io"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/Joker/jade"
	"github.com/jmoiron/sqlx"
)

const (
	VIEWS_DIR = "/root/isucon/webapp/golang/views"
)

var (
	db *sqlx.DB
	tpl_index *template.Template
	tpl_post *template.Template
)

type Article struct{
	ID int64 `db:"id"`
	Title string `db:"title"`
	Body string `db:"body"`
	CreatedAt time.Time `db:"created_at"`
}

func loadSidebarData() ([]Article, error) {
	articles := []Article{}
	err := db.Select(articles, "SELECT a FROM comment c LEFT JOIN article a ON c.article = a.id GROUP BY a.id ORDER BY MAX(c.created_at) DESC LIMIT 10")
	return articles, err
}

func loadMainData() ([]Article, error) {
	articles := []Article{}
	err := db.Select(articles, "SELECT a FROM article a ORDER BY id DESC LIMIT 10")
	return articles, err
}

func GetRoot(w http.ResponseWriter, r *http.Request) {
	sidebarItems, err := loadSidebarData()
	if err != nil {
		log.Println("Failed to get recently commented articles", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	articles, err := loadMainData()
	if err != nil {
		log.Println("Failed to get recent articles", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = tpl_index.Execute(w, map[string]interface{}{
		"sidebaritems": sidebarItems,
		"articles": articles,
	})
	if err != nil {
		log.Println("Failed to execute template for index", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	sidebarItems, err := loadSidebarData()
	if err != nil {
		log.Println("Failed to get recently commented articles", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = tpl_post.Execute(w, map[string]interface{}{
		"sidebaritems": sidebarItems,
	})
	if err != nil {
		log.Println("Failed to execute template for post", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func PostPost(w http.ResponseWriter, r *http.Request) {
  //To allocate slice for request body
  length, err := strconv.Atoi(r.Header.Get("Content-Length"))
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  //Read body data to parse json
  body := make([]byte, length)
  length, err = r.Body.Read(body)
  if err != nil && err != io.EOF {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }

  //parse json
  var jsonBody map[string]interface{}
  err = json.Unmarshal(body[:length], &jsonBody)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
	_, err = db.Exec("INSERT INTO article SET title=?, body=?", jsonBody["title"], jsonBody["body"])
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
}

func GetArticle(w http.ResponseWriter, r *http.Request) {}

func PostComment(w http.ResponseWriter, r *http.Request) {}

func FormatDate(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func init() {
	db_user := "root"
	db_password := ""
	db_host := "127.0.0.1"
	db_port := "3306"
	dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/isucon?parseTime=true&loc=Local&charset=utf8mb4",
		db_user, db_password, db_host, db_port)

	log.Printf("Connecting to db: %q", dsn)
	db, _ = sqlx.Connect("mysql", dsn)
	for {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(time.Second * 3)
	}

	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)
	log.Printf("Succeeded to connect db.")

	funcMap := template.FuncMap{
		"formatDate": FormatDate,
	}

	tpl_index_str, err := jade.ParseFile(filepath.Join(VIEWS_DIR, "index.jade"))
	if err != nil {
		log.Fatalln("Failed to parse index.jade", err)
	}
	tpl_index = template.Must(template.New("index").Funcs(funcMap).Parse(tpl_index_str))

	tpl_post_str, err := jade.ParseFile(filepath.Join(VIEWS_DIR, "post.jade"))
	if err != nil {
		log.Fatalln("Failed to parse post.jade", err)
	}
	tpl_post = template.Must(template.New("post").Funcs(funcMap).Parse(tpl_post_str))
}

func main() {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(GetRoot)
	r.Methods("GET").Path("/post").HandlerFunc(GetPost)
	r.Methods("POST").Path("/post").HandlerFunc(PostPost)
	r.Methods("GET").Path("/article/{articleid}").HandlerFunc(GetArticle)
	r.Methods("POST").Path("/comment/{articleid}").HandlerFunc(PostComment)
	log.Fatalln(http.ListenAndServe(":5000", r))
}
