package jwk

import (
	"context"
	"net/http"
)

type (
	HandleContext interface {
		HandleContext(context.Context) (context.Context, error)
	}
)

type (
	// for all
	WithContext context.Context
	// for `HandleFetchKeyConfig`,
	//     `HandleFetchKeysConfig`
	//     `HandleFetchSetConfig`
	WithHTTPClient *http.Client
	// for `HandleFetchKeyConfig`
	//     `HandleFetchKeysConfig`
	WithSetToKey string
	// for `HandleFetchKeysConfig`
	WithSetToKeys []string
	// for `HandleDecodeKey`
	//     `HandleDecodeSet`
	WithStrict bool
	// for `HandleDecodeKey`
	//     `HandleDecodeSet`
	WithRecalculate bool
)
