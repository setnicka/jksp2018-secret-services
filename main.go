package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/coreos/go-log/log"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/setnicka/jksp2018-secret-services/state"
)

const (
	SESSION_SECRET      = "bojovyVlkodlakCeskePosty"
	SESSION_MAX_AGE     = 3600 * 24
	SESSION_COOKIE_NAME = "cp_cookie"
	TEMPLATE_DIR        = "templates"
	STATIC_DIR          = "static"

	ORG_LOGIN    = "ksp"
	ORG_PASSWORD = "mamezakladnunamesici" // TODO: load from config file?
)

type Server struct {
	sessionStore sessions.Store
	templates    *template.Template
	state        *state.State
}

type Subdomains map[string]http.Handler

// Global singleton
var server *Server

////////////////////////////////////////////////////////////////////////////////

func (subdomains Subdomains) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	domainParts := strings.Split(r.Host, ".")

	if mux := subdomains[domainParts[0]]; mux != nil {
		// Let the appropriate mux serve the request
		mux.ServeHTTP(w, r)
	} else {
		// Handle 404
		http.Error(w, "Subdomain not found", 404)
	}
}

func main() {
	cookieStore := sessions.NewCookieStore([]byte(SESSION_SECRET))
	cookieStore.MaxAge(SESSION_MAX_AGE)
	//cookieStore.Options.Domain = ".fuf.me"

	server = &Server{
		sessionStore: cookieStore,
		state:        state.Init(),
	}

	server.Start()
}

////////////////////////////////////////////////////////////////////////////////

func newRouter(name string) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/start-hry", loginHandler)
	fs := NoListFileSystem{http.Dir(STATIC_DIR + "/" + name)}
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(fs)))
	return router
}

func (s *Server) Start() {
	log.Info("Starting server...")
	subdomains := make(Subdomains)

	orgRouter := newRouter("org")
	orgRouter.HandleFunc("/login", orgLoginHandler)
	orgRouter.HandleFunc("/teams", authOrg(orgTeamsHandler))
	orgRouter.HandleFunc("/dashboard", authOrg(orgDashboardHandler))
	subdomains["org"] = orgRouter

	// Secret services pages:
	ciaRouter := newRouter("cia")
	ciaRouter.HandleFunc("/", auth(ciaIndexHandler))
	ciaRouter.HandleFunc("/internal", auth(ciaInternalHandler))
	subdomains["cia"] = ciaRouter

	nsaRouter := newRouter("nsa")
	nsaRouter.HandleFunc("/", auth(nsaIndexHandler))
	nsaRouter.HandleFunc("/intranet", auth(nsaInternalHandler))
	subdomains["nsa"] = nsaRouter

	kgbRouter := newRouter("kgb")
	kgbRouter.HandleFunc("/", auth(kgbIndexHandler))
	kgbRouter.HandleFunc("/intranet", auth(kgbInternalHandler))
	subdomains["kgb"] = kgbRouter

	fbiRouter := newRouter("fbi")
	fbiRouter.HandleFunc("/", auth(fbiIndexHandler))
	fbiRouter.HandleFunc("/intranet", auth(fbiInternalHandler))
	subdomains["fbi"] = fbiRouter

	pplRouter := newRouter("ppl")
	pplRouter.HandleFunc("/", auth(pplIndexHandler))
	pplRouter.HandleFunc("/intranet", auth(pplInternalHandler))
	subdomains["ppl"] = pplRouter

	bisRouter := newRouter("bis")
	bisRouter.HandleFunc("/", auth(bisIndexHandler))
	bisRouter.HandleFunc("/intranet", auth(bisInternalHandler))
	subdomains["bis"] = bisRouter

	server.getTemplates()
	log.Info("Server started")

	http.ListenAndServe(":8080", subdomains)
}

func auth(handle http.HandlerFunc, renewAuth ...bool) http.HandlerFunc {
	renew := true
	if len(renewAuth) > 0 {
		renew = renewAuth[0]
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if checkSession(w, r, renew) {
			handle(w, r)
			return
		}
		http.Redirect(w, r, "/start-hry", http.StatusTemporaryRedirect)
	}
}

func authOrg(handle http.HandlerFunc, renewAuth ...bool) http.HandlerFunc {
	renew := true
	if len(renewAuth) > 0 {
		renew = renewAuth[0]
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if checkSession(w, r, renew) && isOrg(r) {
			handle(w, r)
			return
		}
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
}
