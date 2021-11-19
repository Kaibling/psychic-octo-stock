package main

import (
	"github.com/Kaibling/psychic-octo-stock/lib/database"
	"github.com/Kaibling/psychic-octo-stock/models"
)

//"github.com/gin-gonic/gin"

func main() {

	sdb := database.NewDatabaseConnector("local.db")
	sdb.Connect()
	sdb.Migrate(&models.User{})
	//userRepo := repositories.NewUserRepository(sdb)
	//fakes.AddFakeUser(20, userRepo)
	//

	// usa := userRepo.GetUserByID("ckw6y0wxe00002z1dbxchdoqv")
	// b, _ := json.MarshalIndent(usa, "", "  ")
	// fmt.Println(string(b))
	// usa.Password = "sadsad"
	// userRepo.UpdateObject(usa)
	// usa = userRepo.GetUserByID("ckw6y0wxe00002z1dbxchdoqv")
	// b, _ = json.MarshalIndent(usa, "", "  ")
	// fmt.Println(string(b))

	// r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	// r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
