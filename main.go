package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	DB *sql.DB
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	connectionString := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// s is a pointer to a server struct so the same connection can be passed to each endpoint handler
	s := &Server{DB: db}

	router := gin.Default()

	router.GET("/users", s.getUsers) // we pass the get users function that is referenced with s
	router.GET("/users/:id", s.getUsersById)

	router.Run()
}

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
}

func (s *Server) getUsersById(c *gin.Context) {
	id := c.Param("id")
	row := s.DB.QueryRow("SELECT * FROM users WHERE id=$1", id)

	var user User

	if err := row.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Address); err != nil {
		c.IndentedJSON(http.StatusNotFound, user)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (s *Server) getUsers(c *gin.Context) {
	rows, err := s.DB.Query("SELECT * FROM users")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := []User{}

	for rows.Next() {
		var user User

		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Address); err != nil {
			c.IndentedJSON(http.StatusNoContent, gin.H{"message": "scan failed"})
			return
		}

		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}
