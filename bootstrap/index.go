package bootstrap

import (
	"log"
	"os"
	// "time"

	// "github.com/gin-contrib/cors"

	// "github.com/PatrickA727/trisakti-proto/config/app_config"
	"github.com/PatrickA727/trisakti-proto/controllers/adminController"
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
		log.Println("error loading env: ", err.Error());
	}

	PORT := os.Getenv("PORT");

	db, err :=database.ConnectDB();
	if err != nil {
		log.Fatal("error connecting database")
	}
	app := gin.Default();	// Returns gin default engine, similar to express app in js 
	
	err = app.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Fatal(err)
	}

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:5173", "https://bkn25q2k-5173.asse.devtunnels.ms/"}, // Change to specific ngrok URL in production
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "ngrok-skip-browser-warning"},
	// 	ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))

	studentStore := store.NewStudentStore(db)
	AdminStore := store.NewAdminStore(db)
	studentController := studentController.NewController(*studentStore)
	adminController := adminController.NewAdminController(*AdminStore)

	routes.InitRoute(app, *studentController, *adminController);
	routes.InitAdminRoute(app, adminController);

	app.Run(PORT);
}
