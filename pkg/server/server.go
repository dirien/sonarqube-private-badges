package server

import (
	"encoding/base64"
	"github/dirien/sonarqube-private-badges/pkg/config"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
)

type SonarQubeProxy struct {
	httputil.ReverseProxy
	address       string
	remote        *url.URL
	authorization string
}

type SonarQubeProxyer interface {
	director(r *http.Request)
	serveHTTP(w http.ResponseWriter, r *http.Request)
	Server() *http.Server
}

func modifyResponse(r *http.Response) error {
	c := r.StatusCode
	switch {
	case c < http.StatusOK:
		panic(http.StatusBadGateway)
	case c < http.StatusInternalServerError:
		panic(http.StatusNotFound)
	default:
		panic(http.StatusBadGateway)
	}
}

func (s *SonarQubeProxy) director(req *http.Request) {
	req.URL.Scheme = s.remote.Scheme
	req.URL.Host = s.remote.Host
	req.Host = s.remote.Host
	if s.authorization != "" {
		req.Header.Add("Authorization", s.authorization)
	}
}

func (s *SonarQubeProxy) serveHTTP(w http.ResponseWriter, r *http.Request) {
	s.ServeHTTP(w, r)
}

func (s *SonarQubeProxy) Server() *http.Server {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/badges/{type}", s.serveHTTP)
	api.HandleFunc("/project_badges/{type}", s.serveHTTP)
	return &http.Server{Addr: s.address, Handler: r}
}

func basicAuthorization(token string) string {
	if token == "" {
		return ""
	}
	a := []byte(token + ":")
	b := base64.StdEncoding.EncodeToString(a)
	return "Basic " + b
}

func NewSonarQubeProxy(c *config.SonarQubeProxyConfig) *SonarQubeProxy {
	s := new(SonarQubeProxy)
	s.Director = s.director
	s.address = ":" + c.Port
	s.remote = c.Remote
	s.authorization = basicAuthorization(c.Auth)
	return s
}
