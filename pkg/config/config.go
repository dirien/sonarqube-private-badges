package config

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type SonarQubeProxyConfig struct {
	Port   string
	Auth   string
	Metric map[string]string
	Remote *url.URL
}

func getPort() string {
	port := os.Getenv("PORT")
	_, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		panic("strange PORT: " + port)
	}
	return port
}

func getAuth() string {
	return os.Getenv("AUTH")
}

func setSonarQubeHost(c *http.Client) *url.URL {
	sonarHost := os.Getenv("SONAR_HOST")
	u, err := url.Parse("https://" + sonarHost)
	if err != nil {
		panic("could not parse SONAR_HOST: " + sonarHost)
	}
	r, err := c.Head(u.String())
	if err != nil {
		panic("could not issues a HEAD to the SONAR_HOST: " + sonarHost)
	}
	switch r.StatusCode {
	case http.StatusUnauthorized:
		fallthrough
	case http.StatusOK:
		fallthrough
	case http.StatusBadRequest:
		return u
	default:
		panic("just an error SONAR_HOST=" + sonarHost)
	}
}

func LoadConfig() *SonarQubeProxyConfig {
	c := &http.Client{Timeout: 10 * time.Second}
	return &SonarQubeProxyConfig{
		Port:   getPort(),
		Auth:   getAuth(),
		Remote: setSonarQubeHost(c),
	}
}
