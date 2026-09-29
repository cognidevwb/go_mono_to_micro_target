package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// mounts are routes a file of this package adds beside wire's — a saga's
// compensating operation, say. Each appends itself from an init func, so
// wire (the monolith's own wiring) is never edited to reach it.
var mounts []func(r *gin.Engine, d *Deps)

// Routes builds the service's router and mounts its endpoints through wire.
func Routes(d *Deps) (http.Handler, error) {
	r := gin.New()
	r.Use(gin.Recovery())
	if err := wire(r, d); err != nil {
		return nil, err
	}
	for _, m := range mounts {
		m(r, d)
	}
	return r, nil
}
