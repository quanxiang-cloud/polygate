package remotes

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/quanxiang-cloud/polygate/pkg/basic/errcode"
	"github.com/quanxiang-cloud/polygate/pkg/basic/polysign"
	"github.com/quanxiang-cloud/polygate/pkg/lib/httputil"

	"github.com/gin-gonic/gin"
)

type signatureInfo polysign.PolySignatureInfo

// Verify verify the signature
func (s *signatureInfo) Verify() error {
	if s.Signature == "" {
		return errcode.ErrInputMissingArg.FmtError("body", polysign.XBodyPolySignSignature)
	}
	then, err := time.Parse(polysign.XHeaderPolySignTimestampFmt, s.Timestamp)
	if err != nil {
		return errcode.ErrDataFormatInvalid.FmtError(
			"header", polysign.XHeaderPolySignTimestamp,
			polysign.XHeaderPolySignTimestampFmt)
	}
	now := time.Now().UTC()
	elapse := now.Sub(then)
	if !(elapse > 0 && elapse < polysign.PolySignatureTimeout) {
		diff := elapse / time.Second * time.Second
		return errcode.ErrInputValueExpired.FmtError("header", polysign.XHeaderPolySignTimestamp, diff.String())
	}
	switch s.SignMethod {
	case polysign.XHeaderPolySignMethodVal:
	default:
		return errcode.ErrInputValueInvalid.FmtError("header", polysign.XHeaderPolySignMethod)
	}
	switch s.SignVersion {
	case polysign.XHeaderPolySignVersionVal:
	default:
		return errcode.ErrInputValueInvalid.FmtError("header", polysign.XHeaderPolySignVersion)
	}

	return nil
}

func (s *signatureInfo) marshalJSON() ([]byte, error) {
	return json.Marshal(s)
}

// get and delete field
func popString(data map[string]interface{}, name string) string {
	if d, ok := data[name]; ok {
		delete(data, name)
		if s, ok := d.(string); ok {
			return s
		}
	}
	return ""
}

func parsePolySignatureInfo(c *gin.Context) (*signatureInfo, error) {
	var err error
	getHeader := func(header string) string {
		ret := c.GetHeader(header)
		if ret == "" && err == nil {
			err = errcode.ErrInputMissingArg.FmtError("header", header)
		}
		return ret
	}

	var signInfo signatureInfo
	if err := httputil.GetRequestArgs(c, &signInfo.Body); err != nil {
		return nil, errcode.ErrParameterError.FmtError(err.Error())
	}

	// BUG: http: proxy error: net/http: HTTP/1.x transport connection broken: http: ContentLength=177 with Body length 0
	if c.Request.Method != http.MethodGet {
		var body json.RawMessage
		if err := httputil.BindBody(c, &body); err != nil {
			return nil, err
		}
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
	}

	signInfo.Signature = popString(signInfo.Body, polysign.XBodyPolySignSignature)
	signInfo.AccessKeyID = getHeader(polysign.XHeaderPolySignKeyID)
	signInfo.Timestamp = getHeader(polysign.XHeaderPolySignTimestamp)
	signInfo.SignMethod = getHeader(polysign.XHeaderPolySignMethod)
	signInfo.SignVersion = getHeader(polysign.XHeaderPolySignVersion)
	if err != nil {
		return nil, err
	}

	if err := signInfo.Verify(); err != nil {
		return nil, err
	}
	return &signInfo, nil
}
