package remotes

import (
	"bytes"
	"io"
	"net/http"

	"github.com/quanxiang-cloud/polygate/pkg/basic/consts"
	"github.com/quanxiang-cloud/polygate/pkg/basic/errcode"
	"github.com/quanxiang-cloud/polygate/pkg/basic/header"
	"github.com/quanxiang-cloud/polygate/pkg/basic/polysign"
	"github.com/quanxiang-cloud/polygate/pkg/lib/httputil"

	"github.com/gin-gonic/gin"
)

var _ Client = (*authorizer)(nil)

// NewAuthClient create a client to authrize by key or token
func NewAuthClient(byKey, byToken *httputil.HTTPClient) Client {
	return &authorizer{
		byKey:   (*authKey)(byKey),
		byToken: (*authToken)(byToken),
	}
}

type authorizer struct {
	byKey   *authKey
	byToken *authToken
}

func (r *authorizer) Request(c *gin.Context) error {
	if token := c.GetHeader(header.HeaderAccessToken); token != "" {
		if err := r.byToken.request(c, nil); err != nil {
			return err
		}
	} else {
		dnsName, ok := c.Params.Get(consts.PathArgServiceName)
		if !(ok && dnsName == "polyapi") { //NOTE: dont allow auth by key except polyapi
			return errcode.ErrInputMissingArg.FmtError("header", header.HeaderAccessToken)
		}

		if err := r.authBykey(c); err != nil {
			return err
		}
	}
	return nil
}

func (r *authorizer) authBykey(c *gin.Context) error {

	signInfo, err := parsePolySignatureInfo(c)
	if err != nil {
		return err
	}

	signature := signInfo.Signature
	signInfo.Signature = "" //NOTE: remove x_polyapi_signature from body

	body, err := signInfo.marshalJSON()
	if err != nil {
		return err
	}

	arg := &requestArg{
		accessKeyID: signInfo.AccessKeyID,
		signature:   signature,
		body:        body,
	}

	if err := r.byKey.request(c, arg); err != nil {
		return err
	}
	return nil
}

//------------------------------------------------------------------------------

type authKey httputil.HTTPClient

type authKeyReq struct {
	Key string `json:"key"`
}

type authKeyResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Secret   string `json:"secret"`
		UserInfo struct {
			UserID       string `json:"userID"`
			UserName     string `json:"userName"`
			DepartmentID string `json:"departmentID"`
		} `json:"userInfo"`
	} `json:"data"`
}

func (r *authKey) request(c *gin.Context, arg *requestArg) error {
	req := &http.Request{
		Method: http.MethodPost,
		Header: http.Header{},
		URL:    r.URI,
		Body:   io.NopCloser(bytes.NewReader(arg.body)),
	}
	req.Header.Set(header.HeaderContentType, header.MIMEJSON)
	req.Header.Set(polysign.XHeaderPolySignKeyID, arg.accessKeyID)

	resp, err := r.Client.Do(req)
	if err != nil {
		log.PutError(err, "AuthByKey.Request", ginRequestID(c))
		return err
	}

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(resp.StatusCode)
		return errcode.ErrInternal.FmtError(resp.Status)
	}

	expect := resp.Header.Get(polysign.XInternalHeaderPolySignSignature)
	got := arg.signature
	if expect == "" || got != expect {
		return errcode.ErrInputArgValidateMismatch.FmtError("body", polysign.XBodyPolySignSignature)
	}

	cloneProfile(&c.Request.Header, resp.Header)
	c.Request.Header.Set(header.HeaderAccessKeyID, header.HeaderPrefixAccessKeyID+arg.accessKeyID) // auth by access key
	delete(c.Request.Header, polysign.XHeaderPolySignKeyID)
	delete(c.Request.Header, polysign.XHeaderPolySignMethod)
	delete(c.Request.Header, polysign.XHeaderPolySignVersion)
	delete(c.Request.Header, polysign.XHeaderPolySignTimestamp)

	return nil
}

//------------------------------------------------------------------------------

// AuthToken auth by token at remote
type authToken httputil.HTTPClient

func (r *authToken) request(c *gin.Context, arg *requestArg) error {
	req := &http.Request{
		Header: c.Request.Header.Clone(),
		URL:    r.URI,
	}

	req.Header.Set(header.HeaderContentType, header.MIMEJSON)

	resp, err := r.Client.Do(req)
	if err != nil {
		log.PutError(err, "AuthByToken.Request", ginRequestID(c))
		return err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(resp.StatusCode)
		return errcode.ErrInternal.FmtError(resp.Status)
	}

	delete(c.Request.Header, header.HeaderAccessToken)
	cloneProfile(&c.Request.Header, resp.Header)

	return nil
}
