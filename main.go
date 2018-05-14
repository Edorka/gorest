package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "net/http"
    "log"
    "fmt"
    "strconv"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"

)

// Book

type Book struct {
    ID      int `json:"id"`
    ISBN    string `json:"isbn"`
    Title   string `json:"title"`
    Author  *Author `json:"author"`
}

type Author struct {
    Firstname  string `json:"firstname"`
    Lastname   string `json:"lastname"`
}

var books []Book

type BooksApp struct {
    Access *sql.DB
}

func (app *BooksApp) Source(source *sql.DB) {
    app.Access = source
}

func (app *BooksApp) List(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    const query = `SELECT id,isbn,title from books`
    selection, err := app.Access.Query(query)
    if err != nil {
        panic(err.Error())
    }
    result := json.NewEncoder(response)
    for selection.Next() {
        var current Book = Book{} 
        err = selection.Scan(&current.ID, &current.ISBN, &current.Title)
        if err != nil {
            panic(err.Error())
        }
        result.Encode(current)
    }
}

func (app *BooksApp) Get(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    params := mux.Vars(request)
    expectedId, err := strconv.ParseInt(params["id"], 10, 32)
    if err != nil {
        panic(err.Error())
    }
    const query = `SELECT id,isbn,title from books where id = $1 `
    var result Book = Book{} 
    err = app.Access.QueryRow(query, expectedId).Scan(&result.ID, &result.ISBN, &result.Title)
    json.NewEncoder(response).Encode(result)
    
}

func (app *BooksApp) Create(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    var newBook Book
    _ = json.NewDecoder(request.Body).Decode(&newBook)
    const query = "INSERT INTO books(id, isbn, title) VALUES(?,?,?)"
    insertion, err := app.Access.Prepare(query)
    if err != nil {
        panic(err.Error())
    }
    insertion.Exec(&newBook.ID, &newBook.ISBN, &newBook.Title)
    response.WriteHeader(http.StatusOK)
    json.NewEncoder(response).Encode(newBook)
}

func (app *BooksApp) Update(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    params := mux.Vars(request)
    expectedId, err := strconv.ParseInt(params["id"], 10, 32)
    if err != nil {
        http.Error(response, fmt.Sprintf( "malformed id: %v", params["id"]), 400)
        return 
    }
    const query = "UPDATE books SET isbn=?, title=?  where id = ?"
    update, err := app.Access.Prepare(query)
    if err != nil {
        panic(err.Error())
    }
    var target Book
    _ = json.NewDecoder(request.Body).Decode(&target)
    update.Exec(&target.ISBN, &target.Title, expectedId)
    response.WriteHeader(http.StatusOK)
    json.NewEncoder(response).Encode(target)
}

func (app *BooksApp) Delete(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    params := mux.Vars(request)
    expectedId, err := strconv.ParseInt(params["id"], 10, 32)
    if err != nil {
        http.Error(response, fmt.Sprintf( "malformed id: %v", params["id"]), 400)
        return 
    }
    const query = "DELETE FROM books where id = ?"
    update, err := app.Access.Prepare(query)
    if err != nil {
        panic(err.Error())
    }
    update.Exec(expectedId)
    response.WriteHeader(http.StatusOK)
    //       fmt.Printf("not found :%d\n", expectedId)
    //    http.NotFound(response, request)
}

func main() {
    router := mux.NewRouter()
    

    books = append(books, Book{ID: 1, ISBN: "438227", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
    books = append(books, Book{ID: 2, ISBN: "454555", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})

    db, err := sql.Open("sqlite3", "file:foo.db?_loc=auto")
    if err != nil {
        panic(err)
    }
    app := BooksApp{}
    app.Source(db)

    router.HandleFunc("/api/books", app.List).Methods("GET")
    router.HandleFunc("/api/book/{id}", app.Get).Methods("GET")
    router.HandleFunc("/api/books", app.Create).Methods("POST")
    router.HandleFunc("/api/book/{id}", app.Update).Methods("PUT")
    router.HandleFunc("/api/book/{id}", app.Delete).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":3000", router))
}

