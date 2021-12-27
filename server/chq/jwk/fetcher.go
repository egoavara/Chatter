package jwk

import (
	"context"
	"fmt"
	"net/http"
)

const (
	ctxFromSetToKeys  = "ctxFromSetToKeys"
	ctxFromSetToKey   = "ctxFromSetToKey"
	ctxFromHTTPClient = "ctxFromClient"
)

// func FetchKey(url interface{}, handles ...HandleContext) (k *Key, err error) {
// 	var ctx = context.Background()
// 	if len(handles) > 0 {
// 		if v, ok := handles[0].(WithContext); ok {
// 			ctx = v
// 			handles = handles[1:]
// 		}
// 	}
// 	for _, handle := range handles {
// 		ctx, err = handle.HandleContext(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return FetchKeyBy(url, ctx)
// }
// func FetchKeys(url interface{}, handles ...HandleContext) (ks []*Key, err error) {
// 	var ctx = context.Background()
// 	for _, handle := range handles {
// 		ctx, err = handle.HandleContext(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return FetchKeysBy(url, ctx)
// }
func FetchSet(url interface{}, handles ...HandleContext) (s *Set, err error) {
	var ctx = context.Background()
	for _, handle := range handles {
		ctx, err = handle.HandleContext(ctx)
		if err != nil {
			return nil, err
		}
	}
	return FetchSetBy(ctx, url)
}

// func FetchKeyBy(url interface{}, ctx context.Context) (*Key, error) {
// 	hclt, ok := utilInterfaceOr(ctx.Value(ctxFromHTTPClient), http.DefaultClient).(*http.Client)
// 	if !ok {
// 		return nil, fmt.Errorf("'%s' must be %T, but got %T", ctxFromHTTPClient, http.DefaultClient, hclt)
// 	}
// 	res, err := utilResponse(url, ctx, hclt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	if res.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
// 	}

// 	ioutil.ReadAll(res.Body)
// 	return nil, nil
// }
// func FetchKeysBy(url interface{}, ctx context.Context) ([]*Key, error) {
// 	hclt, ok := utilInterfaceOr(ctx.Value(ctxFromHTTPClient), http.DefaultClient).(*http.Client)
// 	if !ok {
// 		return nil, fmt.Errorf("'%s' must be %T, but got %T", ctxFromHTTPClient, http.DefaultClient, hclt)
// 	}
// 	res, err := utilResponse(url, ctx, hclt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer res.Body.Close()
// 	if res.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
// 	}
// 	return nil, nil
// }fqwqe
func FetchSetBy(ctx context.Context, url interface{}) (*Set, error) {
	hclt, ok := utilInterfaceOr(ctx.Value(ctxFromHTTPClient), http.DefaultClient).(*http.Client)
	if !ok {
		return nil, fmt.Errorf("'%s' must be %T, but got %T", ctxFromHTTPClient, http.DefaultClient, hclt)
	}
	res, err := utilResponse(url, ctx, hclt)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
	return DecodeSetBy(ctx, res.Body)
}
