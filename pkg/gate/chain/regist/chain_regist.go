package regist

import (
	"errors"

	"github.com/quanxiang-cloud/polygate/pkg/gate/chain"
)

// MustRegist regist a chain node
func MustRegist(name string, h chain.Handler) {
	if err := inst.Reg(name, h); err != nil {
		panic(err)
	}
}

// ToChain return registed chain nodes
func ToChain() *chain.Node {
	return inst.ToChain()
}

//------------------------------------------------------------------------------

var inst chainList

// ChainReg if the data for reg chain node
type chainReg struct {
	Name    string
	Handler chain.Handler
}

type chainList []*chainReg

func (l *chainList) Reg(name string, h chain.Handler) error {
	if h.Name() == "" && name == "" {
		return errors.New("Reg noname handler")
	}
	*l = append(*l, &chainReg{
		Name:    name,
		Handler: h,
	})
	return nil
}

func (l chainList) ToChain() *chain.Node {
	var c *chain.Node
	for i := len(l) - 1; i >= 0; i-- {
		p := l[i]
		c = &chain.Node{
			Name: p.Name,
			H:    p.Handler,
			Next: c,
		}
	}
	return c
}
