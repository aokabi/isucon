package main

import (
	"log"
	"net/http"
	"time"
	"path/filepath"
	"fmt"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/Joker/jade"
	"github.com/jmoiron/sqlx"
)

const (
	RECENT_COMMENTED_ARTICLES = "SELECT a.id, a.title FROM comment c LEFT JOIN article a ON c.article = a.id GROUP BY a.id ORDER BY MAX(c.created_at) DESC LIMIT 10"
	RECENT_ARTICLES = "SELECT id,title,body,created_at FROM article ORDER BY id DESC LIMIT 10"
	VIEWS_DIR = "/root/isucon/webapp/golang/views"
)

var (
	db *sqlx.DB
	tpl_index *template.Template
)

type Article struct{
	ID int64 `db:"id"`
	Title string `db:"title"`
	Body string `db:"body"`
	CreatedAt time.Time `db:"created_at"`
}

func loadSidebarData() ([]Article, error) {
	articles := []Article{}
	err := db.Select(articles, RECENT_COMMENTED_ARTICLES)
	return articles, err
}

func loadMainData() ([]Article, error) {
	articles := []Article{}
	err := db.Select(articles, RECENT_ARTICLES)
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

func GetPost(w http.ResponseWriter, r *http.Request) {}

func PostPost(w http.ResponseWriter, r *http.Request) {}

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

	tpl_index_str, err := jade.ParseFile(filepath.Join(VIEWS_DIR, "index.jade"))
	if err != nil {
		log.Fatal("Failed to parse index.jade", err)
	}
	funcMap := template.FuncMap{
		"formatDate": FormatDate,
	}
	tpl_index = template.Must(template.New("index").Funcs(funcMap).Parse(tpl_index_str))
}

func main() {
	r := mux.NewRouter()
	r.Methods("GET").Path("/").HandlerFunc(GetRoot)
	r.Methods("GET").Path("/post").HandlerFunc(GetPost)
	r.Methods("POST").Path("/post").HandlerFunc(PostPost)
	r.Methods("GET").Path("/article/{articleid}").HandlerFunc(GetArticle)
	r.Methods("POST").Path("/comment/{articleid}").HandlerFunc(PostComment)
	log.Fatal(http.ListenAndServe(":5000", r))
}
