package chq

type (
	ModelUser struct {
		Name        string         `json:"name"`
		Description string         `json:"description"`
		Auths       ModelUserAuths `json:"auths"`
	}
	ModelUserAuths struct {
		Google *ModelUserAuthGoogle `json:"google"`
		// Naver *struct {
		// 	Identifier string
		// }
		// Kakao *struct {
		// 	Identifier string
		// }
		// Facebook *struct {
		// 	Identifier string
		// }
	}
	ModelUserAuthGoogle struct {
		SUB string `json:"sub"`
	}
)
