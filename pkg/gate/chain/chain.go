package chain

import (
	"github.com/quanxiang-cloud/polygate/pkg/basic/errcode"

	"github.com/gin-gonic/gin"
)

// Handler is the gate handler
type Handler interface {
	Name() string
	Handle(c *gin.Context) error
}

// HandleFunc adapt func to Handler
type HandleFunc func(c *gin.Context) error

// Handle the the adapt function
func (f HandleFunc) Handle(c *gin.Context) error {
	return f(c)
}

// Name is the name of handler
func (f HandleFunc) Name() string {
	return ""
}

// Node is the chain node
type Node struct {
	H    Handler
	Name string
	Next *Node
}

// GetName return name of node
func (n *Node) GetName() string {
	if n.Name != "" {
		return n.Name
	}
	return n.H.Name()
}

// Handle handle the chain
func (n *Node) Handle(c *gin.Context) error {
	for p := n; p != nil; p = p.Next {
		if err := p.Handle(c); err != nil {
			return errcode.ErrGateError.FmtError(n.GetName(), err.Error())
		}
	}
	return nil
}
