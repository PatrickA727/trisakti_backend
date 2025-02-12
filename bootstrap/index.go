package bootstrap

import (
	"log"
	"os"

	// "github.com/PatrickA727/trisakti-proto/config/app_config"
	"github.com/PatrickA727/trisakti-proto/controllers/studentController"
	"github.com/PatrickA727/trisakti-proto/database"
	"github.com/PatrickA727/trisakti-proto/routes"
	"github.com/PatrickA727/trisakti-proto/store"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	// "net/http"
)

func AppBootstrap() {

	err := godotenv.Load();
	if err != nil {
		log.Println("error loading env: ", err);
	}

	PORT := os.Getenv("PORT");

	db, err :=database.ConnectDB();
	if err != nil {
		log.Fatal("error connecting database")
	}
	app := gin.Default();	// Returns gin default engine, similar to express app in js 
	
	studentStore := store.NewStudentStore(db)
	studentController := studentController.NewController(*studentStore)

	routes.InitRoute(app, *studentController);

	app.Run(PORT);
}
