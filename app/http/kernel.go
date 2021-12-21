package http

import "github.com/wxsatellite/goweb/framework/gin"

func NewHttpEngine() (engine *gin.Engine, err error) {

	// 三种mode分别对应了不同的场景。在我们开发调试过程中，使用debug模式就可以了。
	// 在上线的时候，一定要选择release模式。
	// test可以用在测试场景中。
	gin.SetMode(gin.ReleaseMode)

	engine = gin.Default()

	// 绑定路由
	SetRoutes(engine)
	return
}
