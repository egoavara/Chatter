package chq

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"path/filepath"
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/spf13/afero"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/http2"
)

type Config struct {
	Engine struct {
		Net struct {
			Address string `json:"address"`
		} `json:"net"`
		TLS struct {
			Mode          string   `json:"mode"`
			Certification string   `json:"certification,omitempty"`
			PrivateKey    string   `json:"private-key,omitempty"`
			CacheDir      string   `json:"cache-dir,omitempty"`
			Domains       []string `json:"domains,omitempty"`
			HTTP2         bool     `json:"http2"`
		} `json:"tls"`
		Favicon string `json:"favicon"`
		Robots  string `json:"robots"`
		Statics string `json:"statics"`
	} `json:"engine"`
	Couchbase struct {
		Address   string              `json:"address"`
		Username  string              `json:"username"`
		Password  string              `json:"password"`
		Define    CouchbaseDefinition `json:"define"`
		WaitUntil int                 `json:"wait-until"`
	} `json:"couchbase"`
	Auth struct {
		Private json.RawMessage `json:"jwk-private"`
		Public  json.RawMessage `json:"jwk-public"`
	} `json:"auth"`
	Google struct {
		WebCert struct {
			URL         string `json:"url"`
			MinInterval int    `json:"min-interval"`
		} `json:"jwk-web"`
		AppID string `json:"appid"`
		Token string `json:"token"`
	} `json:"google"`
}

func New(r io.Reader) (*Config, error) {
	var c Config
	err := json.NewDecoder(r).Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
func MapReadCloser(rc io.ReadCloser, err error) (*Config, error) {
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return New(rc)
}

func (s *Config) ChatterQ() (*ChatterQ, error) {
	q := &ChatterQ{
		Listener:        nil,
		Autocert:        nil,
		Gin:             nil,
		Couchbase:       nil,
		CouchbaseDefine: &s.Couchbase.Define,
		JWKGoogle:       nil,
		JWKGoogleURL:    "",
		JWKGoogleAppID:  "",
		JWKGoogleToken:  "",
		JWKSelfPrivate:  nil,
		JWKSelfPublic:   nil,
	}
	err := NewStackError().
		// FS
		Handle(func() error {
			base, err := filepath.Abs(s.Engine.Statics)
			if err != nil {
				return err
			}
			q.Statics = afero.NewBasePathFs(afero.NewOsFs(), base)
			return nil
		}).
		// Gin
		Handle(func() error {
			q.Gin = gin.New()
			q.Gin.Use(q.Plugin)
			URoute().Methods("GET", "HEAD").Paths("/", "/index.html").Build(q.Gin, func(c *gin.Context) { q.ServeFile(c, "index.html") })
			URoute().Methods("GET", "HEAD").Paths("/robots.txt").Build(q.Gin, func(c *gin.Context) { q.ServeFile(c, "robots.txt") })
			URoute().Methods("GET", "HEAD").Paths("/favicon.ico").Build(q.Gin, func(c *gin.Context) { q.ServeFile(c, "favicon.ico") })
			URoute().Methods("GET", "HEAD").Paths("/statics/*filepath").Build(q.Gin, q.FileSystem)
			SetupV1(q.Gin.Group("/api"))
			return nil
		}).
		// TCP Listener Config
		Handle(func() error {
			var tlscfg = &tls.Config{
				NextProtos: []string{"http/1.1"},
			}
			if s.Engine.TLS.HTTP2 {
				tlscfg.NextProtos = append([]string{http2.NextProtoTLS}, tlscfg.NextProtos...)
			}
			switch s.Engine.TLS.Mode {
			case "Let's Encrypt":
				q.Autocert = &autocert.Manager{
					Prompt:     autocert.AcceptTOS,
					HostPolicy: autocert.HostWhitelist(s.Engine.TLS.Domains...),
					Cache:      autocert.DirCache(s.Engine.TLS.CacheDir),
				}
				tlscfg.GetCertificate = q.Autocert.GetCertificate
			case "X509":
				cert, err := tls.LoadX509KeyPair(s.Engine.TLS.Certification, s.Engine.TLS.PrivateKey)
				if err != nil {
					return err
				}
				tlscfg.Certificates = append(tlscfg.Certificates, cert)
			}

			listen, err := tls.Listen("tcp", s.Engine.Net.Address, tlscfg)
			if err != nil {
				return err
			}
			q.Listener = listen
			return nil
		}, func() error { return q.Listener.Close() }).
		// Couchbase
		Handle(func() error {
			var err error
			q.Couchbase, err = gocb.Connect(s.Couchbase.Address, gocb.ClusterOptions{
				Username: s.Couchbase.Username,
				Password: s.Couchbase.Password,
			})
			if err != nil {
				return err
			}
			//
			timeout := time.Duration(s.Couchbase.WaitUntil) * time.Second
			if timeout == 0 {
				timeout = time.Minute
			}
			till := make(chan struct{})
			go func() {
				tick := 1 * time.Second
				ticker := time.NewTicker(tick)
				ticksum := time.Duration(0)
				defer ticker.Stop()
			TICKER_LOOP:
				for {
					select {
					case <-ticker.C:
						ticksum += tick
						log.Printf("[couchbase] Ticking : %v\n", ticksum)
					case <-till:
						break TICKER_LOOP
					}
				}
			}()
			for bk := range q.CouchbaseDefine.reflectBuckets() {
				err = q.Couchbase.Bucket(bk).WaitUntilReady(timeout, nil)
			}
			if err != nil {
				return err
			}
			close(till)

			return nil
		}, func() error { return q.Couchbase.Close(nil) }).
		// Google API, Cert
		Handle(func() error {

			var err error
			q.JWKGoogleAppID = s.Google.AppID
			q.JWKGoogleToken = s.Google.Token
			q.JWKGoogle = jwk.NewAutoRefresh(context.Background())
			q.JWKGoogleURL = s.Google.WebCert.URL
			q.JWKGoogle.Configure(s.Google.WebCert.URL, jwk.WithMinRefreshInterval(time.Duration(s.Google.WebCert.MinInterval)*time.Second))
			_, err = q.JWKGoogle.Refresh(context.Background(), s.Google.WebCert.URL)
			if err != nil {
				return err
			}
			return nil
		}).
		// Self JWK
		Handle(func() error {
			var err error
			q.JWKSelfPrivate, err = jwk.ParseKey(s.Auth.Private)
			if err != nil {
				return err
			}
			q.JWKSelfPublic, err = jwk.ParseKey(s.Auth.Public)
			if err != nil {
				return err
			}
			return nil
		}).
		Build(`"config" to "chatter? server" failed`)
	if err != nil {
		return nil, err
	}
	return q, nil
}
