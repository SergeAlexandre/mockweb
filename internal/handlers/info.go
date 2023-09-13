package handlers

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-logr/logr"
	"html/template"
	"mockweb/internal/config"
	"net/http"
	"net/http/httputil"
	"time"
)

var nbrHit = 0
var nbrSession = 0

const (
	sessionIdKey       = "sessionId"
	sessionHitCountKey = "sessionHitCount"
)

var _ http.Handler = &InfoHandler{}

type InfoHandler struct {
	Log            logr.Logger
	SessionManager *scs.SessionManager
}

type infoTmplData struct {
	Now          string
	GlobalHits   int
	Name         string
	Headers      map[string][]string
	SessionCount int
	SessionId    int
	SessionHits  int
}

func (h *InfoHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	dumpRequest(request, h.Log)

	var sessionHitCount int
	nbrHit += 1
	sessionId := h.SessionManager.GetInt(request.Context(), sessionIdKey)
	if sessionId == 0 {
		// New session
		nbrSession += 1
		sessionId = nbrSession
		h.SessionManager.Put(request.Context(), sessionIdKey, sessionId)
		sessionHitCount = 1
	} else {
		sessionHitCount = h.SessionManager.GetInt(request.Context(), sessionHitCountKey) + 1
	}
	h.SessionManager.Put(request.Context(), sessionHitCountKey, sessionHitCount)

	info := &infoTmplData{
		Name:         config.Conf.Name,
		Now:          time.Now().Format("2006-01-02T15:04:05"),
		GlobalHits:   nbrHit,
		SessionCount: nbrSession,
		SessionId:    sessionId,
		SessionHits:  sessionHitCount,
		Headers:      request.Header,
	}
	if config.Conf.NoCache {
		writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate;")
		writer.Header().Set("pragma", "no-cache")
	}
	writer.WriteHeader(http.StatusOK)
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
		<tr><td>Name:</td><td>{{ .Name }}</td></tr>
		<tr><td>Now:</td><td>{{ .Now }}</td></tr>
		<tr><td>Global hits:</td><td>{{ .GlobalHits }}</td></tr>
		<tr><td>Session count:</td><td>{{ .SessionCount }}</td></tr>
		<tr><td>Session ID:</td><td>{{ .SessionId }}</td></tr>
		<tr><td>Session hits:</td><td>{{ .SessionHits }}</td></tr>
	</table>
	<h3>Headers:</h3>
	<table>
		{{ range $key, $value := .Headers }}
			<tr><td>{{ $key }}:</td><td>{{ $value }}</td></tr>
		{{ end }}
	</table>
  </body>
</html>
`))

func renderInfo(w http.ResponseWriter, info *infoTmplData) {
	renderTemplate(w, infoTmpl, info)
}

func dumpRequest(r *http.Request, log logr.Logger) {
	if log.V(0).Enabled() {
		if log.V(1).Enabled() {
			dump, err := httputil.DumpRequest(r, true)
			if err != nil {
				log.Error(err, "Error on httputil.DumpRequest(...)")
			}
			log.V(1).Info("-----> HTTP Request", "request", dump)
			//log.V(2).Info(fmt.Sprintf("%q", dump))
			//for hdr := range r.Header {
			//	httpLog.V(2).Info(fmt.Sprintf("Header:%s - > %v", hdr, r.Header[hdr]))
			//}
		} else {
			log.V(0).Info("-----> HTTP Request", "method", r.Method, "uri", r.RequestURI, "remote", r.RemoteAddr)
		}
	}
}
