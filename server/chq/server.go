package chq

import (
	"context"
	"log"
	"net"
	"net/http"
	"path"
	"reflect"
	"strings"

	"github.com/couchbase/gocb/v2"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/spf13/afero"
	"golang.org/x/crypto/acme/autocert"
)

type ChatterQ struct {
	// Go
	Listener net.Listener
	Autocert *autocert.Manager
	Statics  afero.Fs
	// Gin
	Gin *gin.Engine
	// DB
	Couchbase       *gocb.Cluster
	CouchbaseDefine *CouchbaseDefinition
	// JWT
	JWKGoogle      *jwk.AutoRefresh
	JWKGoogleURL   string
	JWKGoogleAppID string
	JWKGoogleToken string
	JWKSelfPrivate jwk.Key
	JWKSelfPublic  jwk.Key
}

type CouchbaseDefinition struct {
	User    CouchbaseLocation
	Session CouchbaseLocation
}
type CouchbaseLocation struct {
	Bucket     string
	Scope      string
	Collection string
}

func (cd *CouchbaseDefinition) reflectBuckets() map[string]struct{} {
	rcl := reflect.TypeOf(CouchbaseLocation{})
	tp := reflect.ValueOf(cd).Elem()
	res := make(map[string]struct{})
	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		if field.Type() == rcl {
			res[field.FieldByName("Bucket").String()] = struct{}{}
		}
	}
	return res
}
func (cl *CouchbaseLocation) ToCollection(c *gocb.Cluster) *gocb.Collection {
	var bk *gocb.Bucket
	var sp *gocb.Scope
	var co *gocb.Collection
	bk = c.Bucket(cl.Bucket)
	if len(cl.Scope) > 0 {
		sp = bk.Scope(cl.Scope)
	} else {
		sp = bk.DefaultScope()
	}
	if len(cl.Collection) > 0 {
		co = sp.Collection(cl.Collection)
	} else {
		co = sp.Collection("_default")
	}

	return co
}
func (chq *ChatterQ) Close() error {
	return NewStackError().
		Handle(func() error {
			return chq.Couchbase.Close(nil)
		}).
		Build("chatter? close failed")
}
func (chq *ChatterQ) Run() error {
	return chq.Gin.RunListener(chq.Listener)
}

func (chq *ChatterQ) Plugin(c *gin.Context) {
	SetCouchbase(c, chq.Couchbase, chq.CouchbaseDefine)
	if gpub, err := chq.JWKGoogle.Fetch(context.Background(), chq.JWKGoogleURL); err == nil {
		SetJWKGoogle(c, gpub)
	} else {
		// TODO : Warn Log
		log.Printf("[Warn] : %e\n", err)
	}
	SetJWKSelf(c, chq.JWKSelfPublic, chq.JWKSelfPrivate)
	//
	c.Next()
}
func (chq *ChatterQ) FileSystem(c *gin.Context) {
	chq.ServeFile(c, c.Param("filepath"))
}

func (chq *ChatterQ) ServeFile(c *gin.Context, fpath string) {
	file, err := chq.Statics.Open(fpath)
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if stat.IsDir() {
		if strings.HasSuffix(fpath, "index.html") {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		chq.ServeFile(c, path.Join(fpath, "index.html"))
		return
	}
	http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), file)
}
