package chq

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupV1(r gin.IRouter) {
	// ============= 유저정보 ============= //
	// 로그인
	r.POST("/auth", v1AuthLogin)
	// JWT 수명 주기 재설정
	r.PUT("/auth", v1hMustAuthorized, v1AuthRefresh)
	// 로그아웃
	r.DELETE("/auth", v1hMustAuthorized, v1AuthLogout)
	// 유저정보 요청(자신)
	r.GET("/users", v1hMustAuthorized)
	// 유저정보 수정(자신)
	r.PUT("/users", v1hMustAuthorized)
	// 유저탈퇴(자신)
	r.DELETE("/users")
	//
	// 회원조회
	r.GET("/users/:uid", v1hOptionalAuthorized)
	// ============= 유저정보:메신저 ============= //
	r.GET("/users/:uid/talk", v1hMustAuthorized)
	r.PUT("/users/:uid/talk", v1hMustAuthorized)
	r.POST("/users/:uid/talk", v1hMustAuthorized)
	// ============= 그룹 ============= //
	// 그룹 목록 조회
	r.GET("/groups", v1hOptionalAuthorized)
	//
	r.GET("/groups/:gid", v1hOptionalAuthorized)
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
func v1hMustAuthorized(c *gin.Context) {
	// kset := JWKSelf(c)
	// // token, err := jwt.ParseHeader(c.Request.Header, "Authorization", jwt.WithVerify(jwa.SignatureAlgorithm(kset.Algorithm()), kset))
	// if err != nil {
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	return
	// }

	// c.AbortWithStatus(http.StatusUnauthorized)
	c.Next()
}
func v1hOptionalAuthorized(c *gin.Context) {
	c.Next()
}
func v1AuthLogin(c *gin.Context) {
	var (
		params v1ParamAuthLogin
	)
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusOK, Err(ErrV1Param))
		return
	}
	if err := params.Validate(); err != nil {
		c.JSON(http.StatusOK, Err(err))
		return
	}
	// kset := JWKSelf(c)

	c.JSON(200, "Hello, World")
}
func v1AuthRefresh(c *gin.Context) {
	var (
		params v1ParamAuthRefresh
	)
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusOK, Err(ErrV1Param))
		return
	}
}

func v1AuthLogout(c *gin.Context) {
	var (
		params v1ParamAuthLogout
	)
	if err := c.ShouldBind(&params); err != nil {
		c.JSON(http.StatusOK, Err(ErrV1Param))
		return
	}
}
