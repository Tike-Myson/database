package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var database *sql.DB

//User ...
type User struct {
	UserID   uuid.UUID
	Login    string
	Email    string
	Password string
}

type Post struct {
	PostID        uuid.UUID
	AuthorID      string
	Author        string
	Title         string
	Content       string
	Date          string
	Like          int
	Dislike       int
	Rating        int
	CommentsCount int
	Categories    []string
	ImageURL      string
}

type Comment struct {
	CommentID uuid.UUID
	AuthorID  uuid.UUID
	PostID    uuid.UUID
	Author    string
	Content   string
	Date      string
	Like      int
	Dislike   int
	Rating    int
}

var AllPosts []Post
var AllComments []Comment

func GetUser(login string) uuid.UUID {
	database := GetDB()
	rows, err := database.Query("SELECT * FROM `users` WHERE `login` = $1", login)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	var UserID uuid.UUID
	var UserLogin string
	var password string
	var email string
	for rows.Next() {
		rows.Scan(&UserID, &UserLogin, &password, &email)
		database.Close()
		return UserID
	}
	database.Close()
	return uuid.Nil
}

func SortCommentsByTime(comments []Comment) []Comment {
	for i := 0; i < len(comments); i++ {
		for j := i; j < len(comments); j++ {
			layout := "2006.01.02 15:04:05"
			str1 := comments[i].Date
			str2 := comments[j].Date
			date1, err := time.Parse(layout, str1)
			if err != nil {
				log.Println(err)
			}
			date2, err := time.Parse(layout, str2)
			if err != nil {
				log.Println(err)
			}
			if date1.Before(date2) {
				comments[i], comments[j] = comments[j], comments[i]
			}
		}
	}
	return comments
}

func SortPostsByTime(posts []Post) []Post {
	for i := 0; i < len(posts); i++ {
		for j := i; j < len(posts); j++ {
			layout := "2006.01.02 15:04:05"
			str1 := posts[i].Date
			str2 := posts[j].Date
			date1, err := time.Parse(layout, str1)
			if err != nil {
				log.Println(err)
			}
			date2, err := time.Parse(layout, str2)
			if err != nil {
				log.Println(err)
			}
			if date1.Before(date2) {
				posts[i], posts[j] = posts[j], posts[i]
			}
		}
	}
	return posts
}

func GetPostByID(UserID, PostID string) []Post {
	var CurrentPost Post
	AllPosts = nil
	CreateAllTables()
	database := GetDB()
	statement := CreatePostTable()
	statement.Exec()
	rows, err := database.Query("SELECT * FROM `posts` WHERE PostID = ?", PostID)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&CurrentPost.PostID, &CurrentPost.AuthorID, &CurrentPost.Author, &CurrentPost.Title, &CurrentPost.Content, &CurrentPost.Date, &CurrentPost.ImageURL)
		rowsRating, err := database.Query("SELECT * FROM rating WHERE PostID = ?", CurrentPost.PostID)
		if err != nil {
			log.Println(err.Error())
		}
		CurrentPost.Rating = 0
		var (
			RatingID string
			UserIDD  string
			PostID   string
			Value    string
		)
		for rowsRating.Next() {

			rowsRating.Scan(&RatingID, &UserIDD, &PostID, &Value)
			count, _ := strconv.Atoi(Value)
			CurrentPost.Rating += count
		}
		CurrentPost.Like = 0
		CurrentPost.Dislike = 0
		CurrentPost.CommentsCount = 0

		rowsLikesPersonal, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "1")
		if err != nil {
			log.Println(err.Error())
		}
		for rowsLikesPersonal.Next() {
			CurrentPost.Like = 1
		}

		rowsDislikesPersonal, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "-1")
		if err != nil {
			log.Println(err.Error())
		}
		for rowsDislikesPersonal.Next() {
			CurrentPost.Dislike = 1
		}
		rowsRate, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "0")
		if err != nil {
			log.Println(err.Error())
		}
		for rowsRate.Next() {
			CurrentPost.Like = 0
			CurrentPost.Dislike = 0
		}
		CurrentPost.CommentsCount = len(GetAllComments("", CurrentPost.PostID.String()))
		CurrentPost.Categories = GetCategoryByPostID(CurrentPost.PostID.String())
		AllPosts = append(AllPosts, CurrentPost)
	}
	database.Close()
	AllPosts = SortPostsByTime(AllPosts)
	return AllPosts
}

func GetAllPosts(UserID uuid.UUID, username string) []Post {
	AllPosts = nil
	CreateAllTables()
	database := GetDB()
	statement := CreatePostTable()
	statement.Exec()
	rows, err := database.Query("SELECT * FROM `posts`")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()
	var CurrentPost Post
	for rows.Next() {
		rows.Scan(&CurrentPost.PostID, &CurrentPost.AuthorID, &CurrentPost.Author, &CurrentPost.Title, &CurrentPost.Content, &CurrentPost.Date, &CurrentPost.ImageURL)
		rowsRating, err := database.Query("SELECT * FROM rating WHERE PostID = ?", CurrentPost.PostID)
		if err != nil {
			log.Println(err.Error())
		}
		CurrentPost.Rating = 0
		var (
			RatingID string
			UserIDD  string
			PostID   string
			Value    string
		)
		for rowsRating.Next() {

			rowsRating.Scan(&RatingID, &UserIDD, &PostID, &Value)
			count, _ := strconv.Atoi(Value)
			CurrentPost.Rating += count
		}
		CurrentPost.Like = 0
		CurrentPost.Dislike = 0
		CurrentPost.CommentsCount = 0

		rowsLikesPersonal, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "1")
		if err != nil {
			log.Println(err.Error())
		}
		for rowsLikesPersonal.Next() {
			CurrentPost.Like = 1
		}

		rowsDislikesPersonal, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "-1")
		if err != nil {
			log.Println(err.Error())
		}
		for rowsDislikesPersonal.Next() {
			CurrentPost.Dislike = 1
		}
		rowsRate, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "0")
		if err != nil {
			log.Println(err.Error())
		}
		for rowsRate.Next() {
			CurrentPost.Like = 0
			CurrentPost.Dislike = 0
		}
		CurrentPost.CommentsCount = len(GetAllComments(UserID.String(), CurrentPost.PostID.String()))
		CurrentPost.Categories = GetCategoryByPostID(CurrentPost.PostID.String())
		AllPosts = append(AllPosts, CurrentPost)
	}
	database.Close()
	//fmt.Println(AllPosts)
	AllPosts = SortPostsByTime(AllPosts)
	return AllPosts
}

func GetAllComments(UserID, PostID string) []Comment {
	AllComments = nil
	var CurrentComment Comment
	database := GetDB()
	rows, err := database.Query("SELECT * FROM comments WHERE PostID = ?", PostID)
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		rows.Scan(&CurrentComment.CommentID, &CurrentComment.AuthorID, &CurrentComment.PostID, &CurrentComment.Author, &CurrentComment.Content, &CurrentComment.Date)
		rowsRating, err := database.Query("SELECT * FROM CommentRating WHERE PostID = ? AND CommentID = ?", PostID, CurrentComment.CommentID)
		if err != nil {
			log.Println(err.Error())
		}
		CurrentComment.Rating = 0
		var (
			RatingID  string
			UserIDD   string
			PostID    string
			CommentID string
			Value     string
		)
		for rowsRating.Next() {
			rowsRating.Scan(&RatingID, &UserIDD, &PostID, &CommentID, &Value)
			count, _ := strconv.Atoi(Value)
			CurrentComment.Rating += count
		}

		CurrentComment.Like = 0
		CurrentComment.Dislike = 0

		rowsLikesPersonal, err := database.Query("SELECT * FROM CommentRating WHERE PostID = ? AND UserID = ? AND CommentID = ? AND Value = ?", PostID, UserID, CurrentComment.CommentID, "1")
		if err != nil {
			log.Println(err.Error())
		}

		for rowsLikesPersonal.Next() {
			CurrentComment.Like = 1
		}
		rowsDislikesPersonal, err := database.Query("SELECT * FROM CommentRating WHERE PostID = ? AND UserID = ? AND CommentID = ? AND Value = ?", PostID, UserID, CurrentComment.CommentID, "-1")
		if err != nil {
			log.Println(err.Error())
		}

		for rowsDislikesPersonal.Next() {
			CurrentComment.Dislike = 1
		}
		AllComments = append(AllComments, CurrentComment)
	}
	database.Close()
	AllComments = SortCommentsByTime(AllComments)
	return AllComments
}

//GetDB ...
func GetDB() *sql.DB {
	if database == nil {
		database, err := sql.Open("sqlite3", "./forum.db")
		if err != nil {
			log.Fatal(err)
		}
		return database
	}
	return database
}

func GetPostsByCategory(UserID, CategoryName string) []Post {
	var CurrentPost Post
	AllPosts = nil
	database := GetDB()
	statement := CreatePostTable()
	statement.Exec()
	rowsCategory, err := database.Query("SELECT * FROM `CategoryPostLink` WHERE CategoryName = ?", CategoryName)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rowsCategory.Close()
	var PostsID []string
	var id string
	var name string
	var PostID string

	for rowsCategory.Next() {
		rowsCategory.Scan(&id, &name, &PostID)
		PostsID = append(PostsID, PostID)
	}

	for _, v := range PostsID {
		rows, err := database.Query("SELECT * FROM `posts` WHERE PostID = ?", v)
		if err != nil {
			log.Println(err.Error())
		}
		for rows.Next() {
			rows.Scan(&CurrentPost.PostID, &CurrentPost.AuthorID, &CurrentPost.Author, &CurrentPost.Title, &CurrentPost.Content, &CurrentPost.Date, &CurrentPost.ImageURL)
			rowsRating, err := database.Query("SELECT * FROM rating WHERE PostID = ?", CurrentPost.PostID)
			if err != nil {
				log.Println(err.Error())
			}
			CurrentPost.Rating = 0
			var (
				RatingID string
				UserIDD  string
				PostID   string
				Value    string
			)
			for rowsRating.Next() {

				rowsRating.Scan(&RatingID, &UserIDD, &PostID, &Value)
				count, _ := strconv.Atoi(Value)
				CurrentPost.Rating += count
			}
			CurrentPost.Like = 0
			CurrentPost.Dislike = 0
			CurrentPost.CommentsCount = 0

			rowsLikesPersonal, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "1")
			if err != nil {
				log.Println(err.Error())
			}
			for rowsLikesPersonal.Next() {
				CurrentPost.Like = 1
			}

			rowsDislikesPersonal, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "-1")
			if err != nil {
				log.Println(err.Error())
			}
			for rowsDislikesPersonal.Next() {
				CurrentPost.Dislike = 1
			}
			rowsRate, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "0")
			if err != nil {
				log.Println(err.Error())
			}
			for rowsRate.Next() {
				CurrentPost.Like = 0
				CurrentPost.Dislike = 0
			}
			CurrentPost.CommentsCount = len(GetAllComments("", CurrentPost.PostID.String()))
			CurrentPost.Categories = GetCategoryByPostID(CurrentPost.PostID.String())
			AllPosts = append(AllPosts, CurrentPost)
		}
	}
	AllPosts = SortPostsByTime(AllPosts)
	return AllPosts
}

func GetLikedPostsByUserID(UserID string, PostIDArr []string) []Post {
	AllPosts = nil
	fmt.Println("*********")
	for _, PostID := range PostIDArr {
		var CurrentPost Post
		CreateAllTables()
		database := GetDB()
		statement := CreatePostTable()
		statement.Exec()
		rows, err := database.Query("SELECT * FROM `posts` WHERE PostID = ?", PostID)
		if err != nil {
			log.Fatal(err.Error())
		}
		for rows.Next() {
			rows.Scan(&CurrentPost.PostID, &CurrentPost.AuthorID, &CurrentPost.Author, &CurrentPost.Title, &CurrentPost.Content, &CurrentPost.Date, &CurrentPost.ImageURL)
			rowsRating, err := database.Query("SELECT * FROM rating WHERE PostID = ?", CurrentPost.PostID)
			if err != nil {
				log.Println(err.Error())
			}
			CurrentPost.Rating = 0
			var (
				RatingID string
				UserIDD  string
				PostID   string
				Value    string
			)
			for rowsRating.Next() {

				rowsRating.Scan(&RatingID, &UserIDD, &PostID, &Value)
				count, _ := strconv.Atoi(Value)
				CurrentPost.Rating += count
			}
			CurrentPost.Like = 0
			CurrentPost.Dislike = 0
			CurrentPost.CommentsCount = 0

			rowsLikesPersonal, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "1")
			if err != nil {
				log.Println(err.Error())
			}
			for rowsLikesPersonal.Next() {
				CurrentPost.Like = 1
			}

			rowsDislikesPersonal, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "-1")
			if err != nil {
				log.Println(err.Error())
			}
			for rowsDislikesPersonal.Next() {
				CurrentPost.Dislike = 1
			}
			rowsRate, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "0")
			if err != nil {
				log.Println(err.Error())
			}
			for rowsRate.Next() {
				CurrentPost.Like = 0
				CurrentPost.Dislike = 0
			}
			CurrentPost.CommentsCount = len(GetAllComments("", CurrentPost.PostID.String()))
			AllPosts = append(AllPosts, CurrentPost)
		}
	}
	//database.Close()
	AllPosts = SortPostsByTime(AllPosts)
	//fmt.Println(len(AllPosts))
	return AllPosts
}

func GetPostsByAuthorID(UserID string) []Post {
	var CurrentPost Post
	AllPosts = nil
	CreateAllTables()
	database := GetDB()
	statement := CreatePostTable()
	statement.Exec()
	rows, err := database.Query("SELECT * FROM `posts` WHERE UserID = ?", UserID)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&CurrentPost.PostID, &CurrentPost.AuthorID, &CurrentPost.Author, &CurrentPost.Title, &CurrentPost.Content, &CurrentPost.Date, &CurrentPost.ImageURL)
		rowsRating, err := database.Query("SELECT * FROM rating WHERE PostID = ?", CurrentPost.PostID)
		if err != nil {
			log.Println(err.Error())
		}
		CurrentPost.Rating = 0
		var (
			RatingID string
			UserIDD  string
			PostID   string
			Value    string
		)
		for rowsRating.Next() {

			rowsRating.Scan(&RatingID, &UserIDD, &PostID, &Value)
			count, _ := strconv.Atoi(Value)
			CurrentPost.Rating += count
		}
		CurrentPost.Like = 0
		CurrentPost.Dislike = 0
		CurrentPost.CommentsCount = 0

		rowsLikesPersonal, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "1")
		if err != nil {
			log.Println(err.Error())
		}
		for rowsLikesPersonal.Next() {
			CurrentPost.Like = 1
		}

		rowsDislikesPersonal, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "-1")
		if err != nil {
			log.Println(err.Error())
		}
		for rowsDislikesPersonal.Next() {
			CurrentPost.Dislike = 1
		}
		rowsRate, err := database.Query("SELECT * FROM rating WHERE PostID = ? AND UserID = ? AND Value = ? ", CurrentPost.PostID, UserID, "0")
		if err != nil {
			log.Println(err.Error())
		}
		for rowsRate.Next() {
			CurrentPost.Like = 0
			CurrentPost.Dislike = 0
		}
		CurrentPost.CommentsCount = len(GetAllComments("", CurrentPost.PostID.String()))
		AllPosts = append(AllPosts, CurrentPost)
	}
	database.Close()
	AllPosts = SortPostsByTime(AllPosts)
	return AllPosts
}

//IsValidateUser return User struct
func IsValidateUser(login, passwordChecking string) (User, error) {
	var ValidationUser User
	database := GetDB()
	rowsLogin, err := database.Query("SELECT * FROM `users` WHERE `login` = $1", login)
	if err != nil {
		log.Println(err)
		database.Close()
		return ValidationUser, err
	}
	defer rowsLogin.Close()
	var UserID uuid.UUID
	var UserLogin string
	var password string
	var email string
	for rowsLogin.Next() {
		rowsLogin.Scan(&UserID, &UserLogin, &password, &email)
		d := []byte(passwordChecking)
		pB := []byte(password)
		err := bcrypt.CompareHashAndPassword(pB, d)
		if err != nil {
			return ValidationUser, err
		}
	}
	rowsEmail, err := database.Query("SELECT * FROM `users` WHERE `email` = $1", login)
	if err != nil {
		log.Println(err)
		database.Close()
		return ValidationUser, err
	}
	defer rowsEmail.Close()

	for rowsEmail.Next() {
		rowsEmail.Scan(&UserID, &UserLogin, &password, &email)
		d := []byte(passwordChecking)
		pB := []byte(password)
		err := bcrypt.CompareHashAndPassword(pB, d)
		if err != nil {
			return ValidationUser, err
		}
	}
	ValidationUser.UserID = UserID
	ValidationUser.Login = UserLogin
	ValidationUser.Email = email
	ValidationUser.Password = passwordChecking
	database.Close()
	return ValidationUser, err
}

//HashPassword ...
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CreateAllTables() {
	statement := CreateUsersTable()
	statement.Exec()
	statement = CreatePostTable()
	statement.Exec()
	statement = CreateCommentTable()
	statement.Exec()
	statement = CreateRatingTable()
	statement.Exec()
	statement = CreateRatingCommentTable()
	statement.Exec()
	statement = CreateCategoryTable()
	statement.Exec()
	statement = CreateCategoryPostLinkTable()
	statement.Exec()
}

func GetCategory() []string {
	var Categories []string
	database := GetDB()
	statement := CreateCommentTable()
	statement.Exec()
	rows, err := database.Query("SELECT * FROM categories")
	if err != nil {
		log.Println(err.Error())
		database.Close()
		return Categories
	}
	var (
		CategoryID   string
		CategoryName string
	)
	for rows.Next() {
		rows.Scan(&CategoryID, &CategoryName)
		Categories = append(Categories, CategoryName)
	}
	database.Close()
	return Categories
}

func CreateCategoryTable() *sql.Stmt {
	database := GetDB()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS categories (CategoryID UUID NOT NULL PRIMARY KEY, CategoryName TEXT)")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return statement
}

func GetCategoryByPostID(PostID string) []string {
	database := GetDB()
	var answer []string
	rows, err := database.Query("SELECT * FROM CategoryPostLink WHERE PostID = ?", PostID)
	if err != nil {
		log.Println(err.Error())
		database.Close()
		return nil
	}
	var id string
	var name string
	var post string
	for rows.Next() {
		rows.Scan(&id, &name, &post)
		answer = append(answer, name)
	}
	database.Close()
	return answer
}

func CreateCategoryPostLinkTable() *sql.Stmt {
	database := GetDB()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS CategoryPostLink (CategoryPostLinkID UUID PRIMARY KEY, CategoryName TEXT, PostID UUID, FOREIGN KEY (PostID) REFERENCES posts(PostID))")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return statement
}

//CreateCategoryPostLink ...
func CreateCategoryPostLink(CategoryPostLinkID, PostID uuid.UUID, CategoryName string) bool {
	database := GetDB()
	statement := CreateCategoryPostLinkTable()
	statement.Exec()
	statement, err := database.Prepare("INSERT INTO CategoryPostLink (CategoryPostLinkID, CategoryName, PostID) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err.Error())
		database.Close()
		return false
	}
	statement.Exec(CategoryPostLinkID, CategoryName, PostID)
	database.Close()
	return true
}

func CreateCommentTable() *sql.Stmt {
	database := GetDB()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS comments (CommentID UUID PRIMARY KEY NOT NULL, UserID UUID, PostID UUID, username TEXT, content TEXT, date TEXT, FOREIGN KEY (UserID) REFERENCES users(UserID), FOREIGN KEY (PostID) REFERENCES posts(PostID) )")
	if err != nil {
		log.Fatal(err.Error())
		return statement
	}
	return statement
}

//CreateComment ...
func CreateComment(CommentID, UserID, PostID, username, content, date string) bool {
	database := GetDB()
	statement := CreateCommentTable()
	statement.Exec()
	statement, err := database.Prepare("INSERT INTO comments (CommentID, UserID, PostID, username, content, date) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	statement.Exec(CommentID, UserID, PostID, username, content, date)
	return true
}

func CreatePostTable() *sql.Stmt {
	database := GetDB()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS posts (PostID UUID NOT NULL PRIMARY KEY, UserID UUID, username TEXT, title TEXT, content TEXT, date TEXT, imgUrl TEXT, FOREIGN KEY (UserID) REFERENCES users(UserID))")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return statement
}

//CreatePost ...
func CreatePost(PostID string, UserID uuid.UUID, username string, title string, content string, date string, imgUrl string) bool {
	database := GetDB()
	statement := CreatePostTable()
	statement.Exec()
	statement, err := database.Prepare("INSERT INTO posts (PostID, UserID, username, title, content, date, imgUrl) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	statement.Exec(PostID, UserID, username, title, content, date, imgUrl)
	return true
}

func CreateRatingTable() *sql.Stmt {
	database := GetDB()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS rating (RatingID UUID PRIMARY KEY NOT NULL, UserID UUID, PostID UUID, Value TEXT, FOREIGN KEY (UserID) REFERENCES users(UserID), FOREIGN KEY (PostID) REFERENCES posts(PostID))")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return statement
}

func GetPostID(UserID string) []string {
	var PostIDArr []string
	database := GetDB()
	statement := CreateRatingTable()
	statement.Exec()
	rows, err := database.Query("SELECT * FROM rating WHERE UserID = ? AND Value = ?", UserID, 1)
	if err != nil {
		log.Println(err.Error())
		database.Close()
		return nil
	}
	var RatingID string
	var UserIDD string
	var PostIDD string
	var Value string
	for rows.Next() {
		rows.Scan(&RatingID, &UserIDD, &PostIDD, &Value)
		PostIDArr = append(PostIDArr, PostIDD)
	}
	database.Close()
	return PostIDArr
}

//CreateLike ...
func CreateLike(RatingID, UserID, PostID string) bool {
	database := GetDB()
	statement := CreateRatingTable()
	statement.Exec()
	if !IsRatingExist(UserID, PostID) {
		statement, err := database.Prepare("INSERT INTO rating (RatingID, UserID, PostID, Value) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err.Error())
			database.Close()
			return false
		}
		statement.Exec(RatingID, UserID, PostID, 1)
		return true
	} else {
		count := 0
		rows, err := database.Query("SELECT * FROM rating WHERE UserID = ? AND PostID = ?", UserID, PostID)
		if err != nil {
			log.Println(err.Error())
		}
		var RatingID string
		var UserIDD string
		var PostIDD string
		var Value string
		for rows.Next() {
			rows.Scan(&RatingID, &UserIDD, &PostIDD, &Value)
			count, _ = strconv.Atoi(Value)
		}
		count += 1
		_, err = database.Exec("UPDATE rating SET Value = ? WHERE UserID = ? AND PostID = ?", count, UserID, PostID)
		if err != nil {
			log.Println(err.Error())
			database.Close()
			return false
		}
		database.Close()
		return true
	}
}

//CreateDislike ...
func CreateDislike(RatingID, UserID, PostID string) bool {
	database := GetDB()
	statement := CreateRatingTable()
	statement.Exec()
	if !IsRatingExist(UserID, PostID) {
		statement, err := database.Prepare("INSERT INTO rating (RatingID, UserID, PostID, Value) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err.Error())
			database.Close()
			return false
		}
		statement.Exec(RatingID, UserID, PostID, -1)
		database.Close()
		return true
	} else {
		count := 0
		rows, err := database.Query("SELECT * FROM rating WHERE UserID = ? AND PostID = ?", UserID, PostID)
		if err != nil {
			log.Println(err.Error())
		}
		var RatingID string
		var UserIDD string
		var PostIDD string
		var Value string
		for rows.Next() {
			rows.Scan(&RatingID, &UserIDD, &PostIDD, &Value)
			count, _ = strconv.Atoi(Value)
		}
		count -= 1
		_, err = database.Exec("UPDATE rating SET Value = ? WHERE UserID = ? AND PostID = ?", count, UserID, PostID)
		if err != nil {
			log.Println(err.Error())
			database.Close()
			return false
		}
		database.Close()
		return true
	}
}

func IsRatingExist(UserID, PostID string) bool {
	count := 0
	database := GetDB()
	rows, err := database.Query("SELECT * FROM rating WHERE UserID = ? AND PostID = ?", UserID, PostID)
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		count++
	}
	if count > 0 {
		database.Close()
		return true
	}
	database.Close()
	return false
}

func CreateRatingCommentTable() *sql.Stmt {
	database := GetDB()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS CommentRating (RatingID UUID PRIMARY KEY NOT NULL, UserID UUID, PostID UUID, CommentID UUID, Value TEXT, FOREIGN KEY (UserID) REFERENCES users(UserID), FOREIGN KEY (PostID) REFERENCES posts(PostID), FOREIGN KEY (CommentID) REFERENCES comments(CommentID))")
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return statement
}

//CreateLike ...
func CreateCommentLike(RatingID, UserID, PostID, CommentID string) bool {
	database := GetDB()
	statement := CreateRatingCommentTable()
	statement.Exec()
	if !IsCommentRatingExist(UserID, PostID, CommentID) {
		statement, err := database.Prepare("INSERT INTO CommentRating (RatingID, UserID, PostID, CommentID, Value) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			log.Fatal(err.Error())
			database.Close()
			return false
		}
		statement.Exec(RatingID, UserID, PostID, CommentID, 1)
		database.Close()
		return true
	} else {
		count := 0
		rows, err := database.Query("SELECT * FROM CommentRating WHERE UserID = ? AND PostID = ? AND CommentID = ?", UserID, PostID, CommentID)
		if err != nil {
			log.Println(err.Error())
		}
		var RatingID string
		var UserIDD string
		var PostIDD string
		var CommentIDD string
		var Value string
		for rows.Next() {
			rows.Scan(&RatingID, &UserIDD, &CommentIDD, &PostIDD, &Value)
			count, _ = strconv.Atoi(Value)
		}
		count += 1
		_, err = database.Exec("UPDATE CommentRating SET Value = ? WHERE UserID = ? AND PostID = ? AND CommentID = ?", count, UserID, PostID, CommentID)
		if err != nil {
			log.Println(err.Error())
			database.Close()
			return false
		}
		database.Close()
		return true
	}
}

//CreateDislike ...
func CreateCommentDislike(RatingID, UserID, PostID, CommentID string) bool {
	database := GetDB()
	statement := CreateRatingCommentTable()
	statement.Exec()
	if !IsCommentRatingExist(UserID, PostID, CommentID) {
		statement, err := database.Prepare("INSERT INTO CommentRating (RatingID, UserID, PostID, CommentID, Value) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			fmt.Println("3")
			log.Fatal(err.Error())
			database.Close()
			return false
		}
		statement.Exec(RatingID, UserID, PostID, CommentID, -1)
		database.Close()
		return true
	} else {
		count := 0
		rows, err := database.Query("SELECT * FROM CommentRating WHERE UserID = ? AND PostID = ? AND CommentID = ?", UserID, PostID, CommentID)
		if err != nil {
			log.Println(err.Error())
		}
		var RatingID string
		var UserIDD string
		var PostIDD string
		var CommentIDD string
		var Value string
		for rows.Next() {
			rows.Scan(&RatingID, &UserIDD, &CommentIDD, &PostIDD, &Value)
			count, _ = strconv.Atoi(Value)
		}
		count -= 1
		_, err = database.Exec("UPDATE CommentRating SET Value = ? WHERE UserID = ? AND PostID = ? AND CommentID = ?", count, UserID, PostID, CommentID)
		if err != nil {
			log.Println(err.Error())
			database.Close()
			return false
		}
		database.Close()
		return true
	}
}

func IsCommentRatingExist(UserID, PostID, CommentID string) bool {
	count := 0
	database := GetDB()
	rows, err := database.Query("SELECT * FROM CommentRating WHERE UserID = ? AND PostID = ? AND CommentID = ?", UserID, PostID, CommentID)
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		count++
	}
	if count > 0 {
		database.Close()
		return true
	}
	database.Close()
	return false
}

func CreateUsersTable() *sql.Stmt {
	database := GetDB()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (UserID UUID PRIMARY KEY, login TEXT, password TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err.Error())
	}
	return statement
}

//CreateUser ...
func CreateUser(UserID uuid.UUID, login, password, email string) bool {
	database := GetDB()
	if IsUserExist(login, email) {
		return false
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (UserID UUID PRIMARY KEY, login TEXT, password TEXT, email TEXT)")
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	statement.Exec()
	//log.Println("Users Table successfully created")
	statement, err = database.Prepare("INSERT INTO users (UserID, login, password, email) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err.Error())
		return false
	}
	statement.Exec(UserID, login, password, email)
	database.Close()
	//log.Println("User " + login + " successfully created")
	return true
}

func IsUserExist(login, email string) bool {
	database := GetDB()
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS users (UserID UUID PRIMARY KEY, login TEXT, password TEXT, email TEXT)")
	if err != nil {
		log.Println(err.Error())
	}
	statement.Exec()
	rowsLogin, err := database.Query("SELECT * FROM users WHERE login = ?", login)
	if err != nil {
		log.Println(err.Error())
	}
	//var (
	//	id string
	//	login string
	//	password string
	//	email string
	//)
	count := 0
	for rowsLogin.Next() {
		count++
	}
	rowsEmail, err := database.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		log.Println(err.Error())
	}
	for rowsEmail.Next() {
		count++
	}
	if count > 0 {
		database.Close()
		return true
	}
	database.Close()
	return false
}
