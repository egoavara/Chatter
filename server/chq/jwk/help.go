package jwk

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
)

func utilResponse(rurl interface{}, ctx context.Context, clt *http.Client) (*http.Response, error) {
	var req *http.Request
	var err error
	switch v := rurl.(type) {
	case string:
		req, err = http.NewRequestWithContext(ctx, "GET", v, nil)
	case url.URL:
		req, err = http.NewRequestWithContext(ctx, "GET", v.String(), nil)
	case *url.URL:
		req, err = http.NewRequestWithContext(ctx, "GET", v.String(), nil)
	default:
		err = fmt.Errorf("unsupported url type %T", v)
	}
	if err != nil {
		return nil, err
	}
	return clt.Do(req)
}

func utilInterfaceFn(i interface{}, fndefault func() interface{}) interface{} {
	if i == nil {
		return fndefault()
	}
	return i
}

func utilInterfaceOr(i interface{}, vdefault interface{}) interface{} {
	if i == nil {
		return vdefault
	}
	return i
}

func utilConsumeStr(m map[string]interface{}, k string) (string, bool) {
	if v, ok := m[k]; ok {
		if s, ok := v.(string); ok {
			delete(m, k)
			return s, true
		}
	}
	return "", false
}

func utilConsumeB64url(m map[string]interface{}, k string) ([]byte, bool) {
	if s, ok := utilConsumeStr(m, k); ok {
		bts, err := base64.RawURLEncoding.DecodeString(s)
		if err != nil {
			fmt.Printf("%c\n", s[340])
			fmt.Println(err)
			return nil, false
		}
		return bts, true
	}
	return nil, false
}
