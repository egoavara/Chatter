package chq

import "github.com/pkg/errors"

var (
	ErrV1Param             = errors.New("parameter failed")
	ErrV1ParameterRequired = errors.New("parameter required")
)

type (
	v1Response struct {
		Ok    bool        `json:"ok"`
		Value interface{} `json:"value,omitempty"`
		Error error       `json:"error,omitempty"`
	}
	v1ParamAuthLogin struct {
		OpenID       string `json:"openid,omitempty"`
		IsFirstLogin bool   `json:"is-first-login,omitempty"`
	}
	v1ResultAuthLogin struct {
		JWT string `json:"jwt"`
	}
	v1ParamAuthRefresh struct {
	}
	v1ResultAuthRefresh struct {
		JWT string `json:"jwt"`
	}
	v1ParamAuthLogout struct {
	}
)

func Ok(value interface{}) v1Response { return v1Response{Ok: true, Value: value} }
func Err(err error) v1Response        { return v1Response{Ok: false, Error: err} }

func (v *v1ParamAuthLogin) Validate() error {
	if len(v.OpenID) == 0 {
		return errors.WithMessage(ErrV1ParameterRequired, "there is no `openid` field")
	}
	return nil
}
