package tiretree_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/quanxiang-cloud/polygate/pkg/lib/tiretree"
)

func TestTireTree(t *testing.T) {
	tt := tiretree.NewTireTree()
	if err := tt.BatchInsert([]string{
		"/api/v1/orgs/user/info",
		"/api/v1/orgs/dep/:id/*p",
		"/api/v1/saml/login",
		"/api/v1/saml/login/:id/*p",
	}, tiretree.White); err != nil {
		t.Error(err)
	}

	if err := tt.BatchInsertKV(
		map[string]string{
			"/api/v1/polyapi/*": "form",
		},
	); err != nil {
		t.Error(err)
	}

	fmt.Println(tt.Show())
	type testCase struct {
		p         string
		expect    bool
		expectVal string
	}
	testCases := []*testCase{
		&testCase{"/api/v1/orgs/user/info", true, tiretree.White},
		&testCase{"/api/v1/orgs/dep/foo/bar", true, tiretree.White},
		&testCase{"/api/v1/orgs/dep/x/foo/bar", true, tiretree.White},
		&testCase{"/api/v1/saml/login", true, tiretree.White},
		&testCase{"/api/v1/saml/login/x/bar2", true, tiretree.White},
		&testCase{"/api/v1/polyapi/balabala", true, "form"},
		&testCase{"/api/v1/polyapi", false, ""},
		&testCase{"/api/v1/polyapi/", true, "form"},
	}
	for i, v := range testCases {
		val, got := tt.Match(v.p)
		if got != v.expect || val != v.expectVal {
			t.Errorf("case %d %q expect %v,%v got %v,%v",
				i+1, v.p, v.expect, v.expectVal, got, val)
		}
	}
}

func TestFastSplitKeys(t *testing.T) {
	type testCase struct {
		p      string
		expect []string
	}
	testCases := []*testCase{
		&testCase{"/api/v1/orgs/user/info", []string{"api", "v1", "orgs", "user", "info"}},
		&testCase{"/api/v1/orgs/dep/foo/bar", []string{"api", "v1", "orgs", "dep", "foo", "bar"}},
		&testCase{"/api/v1/orgs/dep/x/foo/bar", []string{"api", "v1", "orgs", "dep", "x", "foo", "bar"}},
		&testCase{"/api/v1/saml/login", []string{"api", "v1", "saml", "login"}},
		&testCase{"/api/v1/saml/login/x/bar2", []string{"api", "v1", "saml", "login", "x", "bar2"}},
		&testCase{"/api/v1/polyapi/balabala", []string{"api", "v1", "polyapi", "balabala"}},
		&testCase{"/api/v1/polyapi", []string{"api", "v1", "polyapi"}},
		&testCase{"/api/v1/polyapi/", []string{"api", "v1", "polyapi", ""}},
	}
	for i, v := range testCases {
		got := tiretree.FastSplit(v.p, '/', 0)
		if !reflect.DeepEqual(got, v.expect) {
			t.Errorf("case %d %q expect %v got %#v",
				i+1, v.p, v.expect, got)
		}
	}
}
