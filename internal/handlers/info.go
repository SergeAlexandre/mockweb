package handlers

import (
	"github.com/go-logr/logr"
	"html/template"
	"mockweb/pkg/skserver"
	"net/http"
	"time"
)

var nbrHit = 0

var _ skserver.LoggingHandler = &InfoHandler{}

type InfoHandler struct {
	log logr.Logger
}

func (h *InfoHandler) GetLog() logr.Logger {
	return h.log
}

func (h *InfoHandler) SetLog(logger logr.Logger) {
	h.log = logger
}

type infoTmplData struct {
	Now        string
	GlobalHits int
}

func (h *InfoHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	nbrHit += 1

	info := &infoTmplData{
		Now:        time.Now().Format("2006-01-02T15:04:05"),
		GlobalHits: nbrHit,
	}
	renderInfo(writer, info)
}

var infoTmpl = template.Must(template.New("info.html").Parse(`<html>
  <head>
    <style>
		/* make pre wrap */
		pre {
		 white-space: pre-wrap;       /* css-3 */
		 white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
		 white-space: -pre-wrap;      /* Opera 4-6 */
		 white-space: -o-pre-wrap;    /* Opera 7 */
		 word-wrap: break-word;       /* Internet Explorer 5.5+ */
		}
    </style>
  </head>
  <body>
	<table>
		<tr><td>Now:</td><td>{{ .Now }}</td></tr>
		<tr><td>Global hits:</td><td>{{ .GlobalHits }}</td></tr>
  </body>
</html>
`))

func renderInfo(w http.ResponseWriter, info *infoTmplData) {
	renderTemplate(w, infoTmpl, info)
}
