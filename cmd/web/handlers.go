package main

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"
	"github.com/mohibul75/snippetbox-go-project/internal/models"
)

func (app *application)home(w http.ResponseWriter, r *http.Request){

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err:= app.snippets.Latest()
	if err!=nil{
		app.serverError(w,err)
		return
	}

	data:= app.newTemplateData(r)
	data.Snippets=snippets

	app.render(w, http.StatusOK, "home.tmpl",data)

}

func (app *application)snippetView(w http.ResponseWriter, r *http.Request){

	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err!= nil || id<1 {
		app.notFound(w)
		return
	}

	snippet, err:= app.snippets.Get(id)

	if err!=nil{
		if errors.Is(err, models.ErrNoRecord){
			app.notFound(w)
		}else {
			app.serverError(w,err)
		}
		return
	}

	data:= app.newTemplateData(r)
	data.Snippet=snippet

	app.render(w, http.StatusOK, "view.tmpl",data)

}

func (app *application) snippetCretePost(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost {
		w.Header().Set("Allow","POST")
		app.clientError(w,http.StatusMethodNotAllowed)

		return
	}

	title:="0 snail"
	content:= "0 Snail\nClimb Mount,\nBut Slow"
	expires:=7

	id, err:= app.snippets.Insert(title,content,expires)
	if err!=nil{
		app.serverError(w,err)
		return
	}

	http.Redirect(w,r,fmt.Sprintf("/snippet/view?id=%d",id),http.StatusSeeOther)
}

func (app *application)snippetCrete(w http.ResponseWriter, r *http.Request){

	w.Write([]byte("Display form for creating snippet"))
}
