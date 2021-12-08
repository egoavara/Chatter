package chq

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
)

func SetupV1(r gin.IRouter) {
	// ============= 유저정보 ============= //
	// 로그인
	r.GET("/users")
	// 회원가입
	r.POST("/users", v1UserRegister)
	// 회원조회
	r.GET("/users/:uid")
	// 회원정보 수정
	r.PATCH("/users/:uid")
	// 회원탈퇴
	r.DELETE("/users/:uid")
	// OAuth2 : 커뮤니티 로그인 추가
	r.POST("/users/:uid/auth")
	// OAuth2 : 외부 로그인 연결 확인
	r.GET("/users/:uid/auth")
	// OAuth2 : 커뮤니티 로그인 제거
	r.DELETE("/users/:uid/auth")
	// ============= 지역기반 ============= //
	// 유저의 위치 반환
	r.GET("/geo/user/:uid")
	// 지역 정보 반환
	r.GET("/geo/area/:gid")
	// 지역내 유저수 반환
	r.GET("/geo/area/:gid/users")
	// 접속 요청
	r.POST("/geo/connect")

}

func v1UserRegister(c *gin.Context) {
	// 패러미터값 분석, 필요 리소스 준비
	var param v1ParamRegister
	var cb, def = Couchbase(c)
	var jwkpriv = JWKSelfPrivate(c)
	if err := c.ShouldBind(&param); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	co := def.User.ToCollection(cb)
	// 새로 생성할 유저 모델을 생성
	var model = &ModelUser{
		Name:        "",
		Description: "",
		Auths: ModelUserAuths{
			Google: nil,
		},
	}
	// 인증정보 검토
	atk, err := jwt.ParseHeader(c.Request.Header, "Authorization")
	if err != nil {
		c.JSON(http.StatusOK, Err(errors.New("authorization not setup")))
		return
	}
	switch atk.Issuer() {
	case "google":
		atk, err = jwt.ParseHeader(c.Request.Header, "Authorization", jwt.WithKeySet(JWKGoogle(c)))
		if err != nil {
			c.JSON(http.StatusOK, Err(errors.New("authorization by google failed")))
			return
		}
		log.Println(atk)

	default:
		c.JSON(http.StatusOK, Err(fmt.Errorf("unknown ISS '%s'", atk.Issuer())))
		return
	}
	// TODO : ModelValidate
	// 새 유저를 DB에 업로드 시작
	uid := uuid.New()
	ntk := jwt.New()
	ntk.Set(jwt.IssuerKey, JWKSelfIssuer(c))
	ntk.Set(jwt.IssuedAtKey, time.Now())
	ntk.Set(jwt.SubjectKey, uid.String())
	ntk.Set(jwt.ExpirationKey, time.Now().Add(5*time.Minute))
	sig, err := jwt.Sign(ntk, jwa.SignatureAlgorithm(jwkpriv.Algorithm()), jwkpriv.Algorithm)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusOK, Err(errors.New("sign failed")))
		return
	}
	if _, err = co.Upsert(uid.String(), model, nil); err != nil {
		c.JSON(http.StatusOK, Err(errors.New("authorization not setup")))
		return
	}
	log.Printf("New user - '%s'\n", uid.String())
	c.JSON(http.StatusOK, Ok(sig))
}
