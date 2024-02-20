package pkgs

import "github.com/gin-gonic/gin"

func SetupGin() {

	g := gin.New()
	//gin.SetMode(gin.DebugMode)

	g.GET("/gin-service1", func(c *gin.Context) {
		c.String(200, "Hello Gin!")
	})

	err := g.Run("0.0.0.0:9001")
	if err != nil {
		panic(err)
	}
}
