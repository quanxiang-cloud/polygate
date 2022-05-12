package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/quanxiang-cloud/cabin/error/errdefiner"
)

// BindBody bind gin body
func BindBody(c *gin.Context, d interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	bb, ok := b.(binding.BindingBody)
	if !ok {
		return errdefiner.NewErrorWithString(errdefiner.ErrParams, "binding type error:"+c.ContentType())
	}
	if err := c.ShouldBindBodyWith(d, bb); err != nil {
		return errdefiner.NewErrorWithString(errdefiner.ErrParams, err.Error())
	}
	return nil
}

// GetRequestArgs get request args
func GetRequestArgs(c *gin.Context, d interface{}) error {
	if d == nil {
		d = &json.RawMessage{}
	}

	// get query only in GET requestion
	if c.Request.Method == http.MethodGet {
		q := c.Request.URL.Query()
		raw := QueryToBody(q, false)
		err := json.Unmarshal([]byte(raw), d)
		return err
	}

	err := BindBody(c, d)
	return err
}
