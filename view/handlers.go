package view

import (
	"html/template"
	"net/http"
)

var params = struct {
	Title string
}{"Labyrinth Demo"}

var playTempl = template.Must(template.New("t").Parse(play_html))

func init() {
	http.HandleFunc("/", playHtmlHandle)
}

func playHtmlHandle(w http.ResponseWriter, r *http.Request) {
	playTempl.Execute(w, params)
}
