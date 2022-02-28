package tiretree_test

import (
	"fmt"
	"testing"

	"github.com/quanxiang-cloud/polygate/pkg/lib/tiretree"
)

func TestTireTree(t *testing.T) {
	tt := tiretree.NewTireTree()
	n := tt.BatchInsert([]string{
		"/api/v1/orgs/user/info",
		"/api/v1/orgs/dep/:id/*p",
		"/api/v1/saml/login",
		"/api/v1/saml/login/:id/*p",
	}, tiretree.White)

	fmt.Println(n)
	println(tt.Show())
	type testCase struct {
		p      string
		expect bool
	}
	testCases := []*testCase{
		&testCase{"/api/v1/orgs/user/info", true},
		&testCase{"/api/v1/orgs/dep/foo/bar", true},
		&testCase{"/api/v1/orgs/dep/x/foo/bar", true},
		&testCase{"/api/v1/saml/login", true},
		&testCase{"/api/v1/saml/login/x/bar2", true},
	}
	for i, v := range testCases {
		_, got := tt.Match(v.p)
		if got != v.expect {
			t.Errorf("case %d %q expect %v got %v", i+1, v.p, v.expect, got)
		}
	}
}
