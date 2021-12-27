package chq

import (
	"time"

	"github.com/couchbase/gocb/v2"
	"github.com/gin-gonic/gin"
)

func SetLogoutCache(c *gin.Context, cache map[string]time.Time) { c.Set("logout-cache", cache) }
func LogoutCache(c *gin.Context) map[string]time.Time {
	return c.Value("logout-cache").(map[string]time.Time)
}

// func JWKGoogle(c *gin.Context) jwk.Set {
// 	return c.Value("jwk-google").(jwk.Set)
// }
// func SetJWKGoogle(c *gin.Context, s jwk.Set) {
// 	c.Set("jwk-google", s)
// }

func JWKSelfIssuer(c *gin.Context) string {
	// TODO : Load from config
	return "https://chatterq.net"
}

// func JWKSelf(c *gin.Context) jwk.Key {
// 	return c.Value("jwk-self").(jwk.Key)
// }
// func SetJWKSelf(c *gin.Context, set jwk.Key) {
// 	c.Set("jwk-self", set)
// }

func Couchbase(c *gin.Context) (*gocb.Cluster, *CouchbaseDefinition) {
	return c.Value("couchbase").(*gocb.Cluster), c.Value("couchbase-definition").(*CouchbaseDefinition)
}

func SetCouchbase(c *gin.Context, cb *gocb.Cluster, def *CouchbaseDefinition) {
	c.Set("couchbase", cb)
	c.Set("couchbase-definition", def)
}

type UniRoute struct {
	methods []string
	paths   []string
}

func URoute() *UniRoute {
	return &UniRoute{
		methods: nil,
		paths:   nil,
	}
}
func (h *UniRoute) Methods(methods ...string) *UniRoute {
	h.methods = append(h.methods, methods...)
	return h
}
func (h *UniRoute) Paths(paths ...string) *UniRoute {
	h.paths = append(h.paths, paths...)
	return h
}
func (h *UniRoute) Build(r gin.IRouter, handlers ...gin.HandlerFunc) {
	for _, method := range h.methods {
		for _, path := range h.paths {
			r.Handle(method, path, handlers...)
		}
	}
}
