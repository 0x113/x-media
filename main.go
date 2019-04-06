package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/0x113/x-media/auth"
	"github.com/0x113/x-media/database/mysql"
	"github.com/0x113/x-media/video"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
		FullTimestamp:   true,
	})
}

func main() {
	jwt_secret := flag.String("jwt-key", "", "Key for generating JWT")
	flag.Parse()
	if *jwt_secret == "" {
		fmt.Println("jwt-key cannot be empty. Run server using -jwt-key flag.")
		os.Exit(0)
	}

	conn := mysqlConnection("xmedia_user", "password", "127.0.0.1", "3306", "xmedia")
	defer conn.Close()

	/* authentication */
	authRepo := mysql.NewMySQLAuthRepository(conn, *jwt_secret)
	authService := auth.NewAuthService(authRepo, *jwt_secret)
	authHandler := auth.NewAuthHandler(authService)

	/* video */
	videoRepo := mysql.NewMySQLVideoRepository(conn, *jwt_secret)
	videoService := video.NewVideoService(videoRepo)
	videoHandler := video.NewVideoHandler(videoService)

	router := mux.NewRouter().StrictSlash(true)

	/* authentication routes */
	router.HandleFunc("/user/create", authHandler.Create).Methods("POST", "GET")
	router.HandleFunc("/user/token/generate", authHandler.GenerateJWT).Methods("POST")

	/* vidoe routes */
	router.HandleFunc("/api/movies/update", videoHandler.UpdateMovies).Methods("GET")

	http.Handle("/", accessControl(router))

	log.Infoln("Launching server on port :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Errorf("ListenAndServe: %s", err)
	}
}

func mysqlConnection(username, password, host, port, dbname string) *sql.DB {
	sqlStmt := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname)
	db, err := sql.Open("mysql", sqlStmt)
	if err != nil {
		log.Errorf("Error while connection to the database: %s", err)
	}
	return db
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
