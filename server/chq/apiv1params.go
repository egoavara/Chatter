package chq

type (
	v1Response struct {
		Ok    bool        `json:"ok"`
		Value interface{} `json:"value,omitempty"`
		Error error       `json:"error,omitempty"`
	}
	v1ParamRegister struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

func Ok(value interface{}) v1Response { return v1Response{Ok: true, Value: value} }
func Err(err error) v1Response        { return v1Response{Ok: false, Error: err} }
