package jwk

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"
)

type (
	innerJWKSet struct {
		Keys []json.RawMessage `json:"keys"`
	}
	innerJWKSetEachKey struct {
		Index int
		Value *Key
	}
	innerJWKKeyLeft struct {
		Value *Key
		Map   map[string]interface{}
	}
)

const (
	ctxJSONStrict  = "ctxJSONStrict"
	ctxRecalculate = "ctxRecalculate"
)

func DecodeKeyBy(ctx context.Context, reader io.Reader) (*Key, error) {
	isStrict := utilInterfaceOr(ctx.Value(ctxJSONStrict), false).(bool)
	recalculate := utilInterfaceOr(ctx.Value(ctxRecalculate), false).(bool)
	doneResult := make(chan interface{}, 1)
	// Step 1 : JSON Decode
	go func() {
		var data map[string]interface{}
		err := json.NewDecoder(reader).Decode(&data)
		select {
		case <-ctx.Done():
			close(doneResult)
		default:
			if err != nil {
				doneResult <- err
			} else {
				doneResult <- data
			}
		}

	}()
	//
	for v := range doneResult {
		switch v := v.(type) {
		// Step 2 : decoded json to jwk key
		case map[string]interface{}:
			go func() {
				var key = new(Key)
				// kty
				if kty, ok := utilConsumeStr(v, "kty"); ok {
					key.KeyType = KeyType(kty)
				} else {
					doneResult <- fmt.Errorf("'kty' must be string")
					return
				}
				// use
				if kty, ok := utilConsumeStr(v, "use"); ok {
					key.KeyUse = KeyUse(kty)
				} else {
					doneResult <- fmt.Errorf("'use' must be string")
					return
				}
				switch key.KeyType {
				case KeyTypeRSA:
					pubk, err := decodeRSAPubKey(v)
					if err != nil {
						doneResult <- fmt.Errorf("'use' must be string")
						return
					}
					prik, err := decodeRSAPriKey(v, pubk, recalculate)
					if err != nil {
						doneResult <- fmt.Errorf("'use' must be string")
						return
					}
					if prik == nil {
						key.Raw = pubk
					} else {
						key.Raw = prik
					}
				case KeyTypeEC:
					// TODO : EC
					panic("unimplemented")
				case KeyTypeOctet:
					// TODO : EC
					panic("unimplemented")
				}
				// jobs done
				select {
				case <-ctx.Done():
					close(doneResult)
				default:
					doneResult <- &innerJWKKeyLeft{
						Value: key,
						Map:   v,
					}
				}
			}()
		// Step 3 : Check Strict Mode
		case *innerJWKKeyLeft:
			if isStrict && len(v.Map) > 0 {
				return nil, fmt.Errorf("strict mode on, there is unconsumed values : %v", v.Map)
			}
			return v.Value, nil
		case error:
			return nil, v
		}
	}
	return nil, errors.New("context done")
}

// map to rsa public key
func decodeRSAPubKey(data map[string]interface{}) (*rsa.PublicKey, error) {
	var n, e *big.Int
	// public
	if bn, ok := utilConsumeB64url(data, "n"); ok {
		n = new(big.Int).SetBytes(bn)
	} else {
		return nil, fmt.Errorf("'n' not exist")
	}
	if be, ok := utilConsumeB64url(data, "e"); ok {
		e = new(big.Int).SetBytes(be)
	} else {
		return nil, fmt.Errorf("'e' not exist")
	}
	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

// decodeRSAPriKey must after decodeRSAPubKey
// it return nil if it is not private key
// for example, when rsa public key, it return nil, nil
// `recalculate` true when you need to do `rsa.PrivateKey.Precompute` manualy, but it automaticaly set this value to true when there is no precomputed values
func decodeRSAPriKey(data map[string]interface{}, pubk *rsa.PublicKey, recalculate bool) (*rsa.PrivateKey, error) {
	var d, p, q, dp, dq, qi *big.Int

	if bd, ok := utilConsumeB64url(data, "d"); ok {
		d = new(big.Int).SetBytes(bd)
	} else {
		// D must exist when this `data` is RSA private key
		return nil, nil
	}
	if bp, ok := utilConsumeB64url(data, "p"); ok {
		p = new(big.Int).SetBytes(bp)
	} else {
		return nil, fmt.Errorf("'p' not exist")
	}
	if bq, ok := utilConsumeB64url(data, "q"); ok {
		q = new(big.Int).SetBytes(bq)
	} else {
		return nil, fmt.Errorf("'q' not exist")
	}
	if !recalculate {
		if bdp, ok := utilConsumeB64url(data, "dp"); ok {
			dp = new(big.Int).SetBytes(bdp)
		} else {
			recalculate = true
		}
		if bdq, ok := utilConsumeB64url(data, "dq"); ok {
			dq = new(big.Int).SetBytes(bdq)
		} else {
			recalculate = true
		}
		if bqi, ok := utilConsumeB64url(data, "qi"); ok {
			qi = new(big.Int).SetBytes(bqi)
		} else {
			recalculate = true
		}
	}
	prik := &rsa.PrivateKey{
		PublicKey: *pubk,
		D:         d,
		Primes:    []*big.Int{p, q},
	}
	if recalculate {
		prik.Precompute()
	} else {
		prik.Precomputed.Dp = dp
		prik.Precomputed.Dq = dq
		prik.Precomputed.Qinv = qi
	}
	return prik, nil
}

func DecodeSetBy(ctx context.Context, reader io.Reader) (*Set, error) {
	doneResult := make(chan interface{}, 1)
	isStrict := utilInterfaceOr(ctx.Value(ctxJSONStrict), false).(bool)
	// Step 1 : parse raw text to json unmarshaled struct
	go func() {
		var data innerJWKSet
		dec := json.NewDecoder(reader)
		if isStrict {
			dec.DisallowUnknownFields()
		}
		err := dec.Decode(&data)
		select {
		case <-ctx.Done():
			close(doneResult)
		default:
			if err != nil {
				doneResult <- err
			} else {
				doneResult <- &data
			}
		}
	}()
	var (
		result = new(Set)
		works  = 0
	)
	for v := range doneResult {
		switch v := v.(type) {
		// Step 2 : make workers for each key
		case *innerJWKSet:
			works = len(v.Keys)
			result.Keys = make([]*Key, len(v.Keys))
			for i, key := range v.Keys {
				go func(i int, key string) {
					res, err := DecodeKeyBy(ctx, strings.NewReader(key))
					select {
					case <-ctx.Done():
						close(doneResult)
					default:
						if err != nil {
							doneResult <- err
						} else {
							doneResult <- &innerJWKSetEachKey{
								Index: i,
								Value: res,
							}
						}
					}
				}(i, string(key))
			}
		// Step 3 : handle each keys and if job is done, return result of set
		case *innerJWKSetEachKey:
			works -= 1
			result.Keys[v.Index] = v.Value
			if works <= 0 {
				return result, nil
			}
		case error:
			return nil, v
		default:
			panic("unreachable")
		}
	}
	return nil, errors.New("context done")
}
