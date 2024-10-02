package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Comment struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}

type Post struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	LikeCount   int       `json:"like_count"`
	Comments    []Comment `json:"comments"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getAllPostsHandler(w, r)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func getAllPosts(db *sql.DB) ([]Post, error) {
	query := `
		SELECT
			p.id AS post_id,
			p.username AS post_username,
			p.title,
			p.description,
			p.created_at AS post_created_at,
			COUNT(DISTINCT l.id) AS like_count,
			c.id AS comment_id,
			c.user_id AS comment_user_id,
			c.username AS comment_username,
			c.comment,
			c.created_at AS comment_created_at
		FROM posts p
		LEFT JOIN comments c ON p.id = c.post_id
		LEFT JOIN likes l ON p.id = l.post_id
		GROUP BY p.id, c.id
		ORDER BY p.created_at DESC, c.created_at ASC;`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	postMap := make(map[int]*Post)

	for rows.Next() {
		var (
			postID              int
			postUsername        string
			title               string
			description         string
			postCreatedAt       string
			likeCount           int
			commentID           sql.NullInt64
			commentUserID       sql.NullInt64
			commentUsername     sql.NullString
			comment             sql.NullString
			commentCreatedAt    sql.NullString
		)

		if err := rows.Scan(&postID, &postUsername, &title, &description, &postCreatedAt, &likeCount, &commentID, &commentUserID, &commentUsername, &comment, &commentCreatedAt); err != nil {
			return nil, err
		}

		if _, exists := postMap[postID]; !exists {
			postMap[postID] = &Post{
				ID:          postID,
				Username:    postUsername,
				Title:       title,
				Description: description,
				CreatedAt:   postCreatedAt,
				LikeCount:   likeCount,
				Comments:    []Comment{},
			}
		}

		if commentID.Valid {
			postMap[postID].Comments = append(postMap[postID].Comments, Comment{
				ID:        int(commentID.Int64),
				UserID:    int(commentUserID.Int64),
				Username:  commentUsername.String,
				Comment:   comment.String,
				CreatedAt: commentCreatedAt.String,
			})
		}
	}

	posts := make([]Post, 0, len(postMap))
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	return posts, nil
}

func getAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	if err := godotenv.Load(); err != nil {
		http.Error(w, "Error loading .env file", http.StatusInternalServerError)
		return
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbDatabase
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		http.Error(w, "Error opening database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		http.Error(w, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := getAllPosts(db)
	if err != nil {
		http.Error(w, "Error fetching posts and comments: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func main() {
	http.HandleFunc("/posts", Handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
