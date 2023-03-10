package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mohibul75/snippetbox-go-project/internal/models"
)

type neuteredFileSystem struct{
	fs http.FileSystem
}

type application struct{
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets  *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {

	addr:= flag.String("addr",":4000","HTTP network address")
	dsn := flag.String("dsn", "web:pass@/snippetbox?parseTime=true", "MySQL data source name")

	flag.Parse()

	infoLog := log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t",log.Ldate|log.Ltime|log.Lshortfile)

	db, err:= openDB(*dsn)

	if err!=nil {
		errorLog.Fatal(err)
	}else{
		infoLog.Printf("DB conection established!!!")
	}

	defer db.Close()

	templateCache, err:= newTemplateCache()

	if err!=nil{
		errorLog.Fatal(err)
	}


	app:= &application{
		errorLog: errorLog,
		infoLog: infoLog,
		snippets: &models.SnippetModel{DB:db},
		templateCache: templateCache,
	}

	srv:= &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}

	infoLog.Printf("Starting Server on  %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string)(*sql.DB, error){
	db, err:= sql.Open("mysql",dsn)

	if err!=nil {
		return nil,err
	}

	if err=db.Ping(); err!=nil{
		return nil, err
	}
	return db, nil
}

func (nfs neuteredFileSystem) Open(path string)(http.File, error){
	f,err := nfs.fs.Open(path)

	if err!=nil {
		return nil,err
	}

	s, err:= f.Stat()

	if s.IsDir() {
		index:= filepath.Join(path,"index.html")

		if _, err := nfs.fs.Open(index); err!= nil{
			closeErr:= f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f,nil
}
