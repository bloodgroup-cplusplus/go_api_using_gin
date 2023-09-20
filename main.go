package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 3},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

func getBooks(c *gin.Context) {
	// we get nicely formatted json
	c.IndentedJSON(http.StatusOK, books)
}

//get book by id

func getBookById(id string) (*book, error) {
	// if the id doesn't exist return error
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
			// so that we can modify attributes
		}
	}
	return nil, errors.New("Book not found")
}

func bookById(c *gin.Context) {
	// get book by id
	id := c.Param("id")
	book, err := getBookById(id)

	// if we do have an error (retun error)
	// return a custom response saying bad request or bad response

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		// gin.H is a custom json writing element
		return

	}

	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if ok == false {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter"})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found"})
		return
	}
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)

}

func createBook(c *gin.Context) {
	var newBook book
	// from context bind the data payload from the content

	if err := c.BindJSON(&newBook); err != nil {
		return
	}
	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)

}

func main() {

	router := gin.Default()
	router.GET("/books", getBooks)
	router.POST("/books", createBook)
	router.GET("/books/:id", bookById) // set up a path parameter
	router.PATCH("/checkout", checkoutBook)
	router.Run("localhost:8080")

}
