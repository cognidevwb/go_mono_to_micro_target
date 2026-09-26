package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Routes builds the service's router and mounts its endpoints through wire.
func Routes(d *Deps) (http.Handler, error) {
	r := gin.New()
	r.Use(gin.Recovery())
	if err := wire(r, d); err != nil {
		return nil, err
	}
	return r, nil
}
