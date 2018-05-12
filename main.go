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

func GetBooks(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    db, err := sql.Open("sqlite3", "file:foo.db?_loc=auto")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    const query = `SELECT id,isbn,title from books`
    selection, err := db.Query(query)
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

func GetBook(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    params := mux.Vars(request)
    expectedId, err := strconv.ParseInt(params["id"], 10, 32)
    if err != nil {
        return
    }
    db, err := sql.Open("sqlite3", "file:foo.db?_loc=auto")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    const query = `SELECT id,isbn,title from books where id = $1 `
    var retval Book
    err = db.QueryRow(query, expectedId).Scan(&retval.ID, &retval.ISBN, &retval.Title)
    json.NewEncoder(response).Encode(retval)
    // response.WriteHeader(http.StatusNotFound)
    //json.NewEncoder(response).Encode(&Book{}) 
}

func GetNextId(books []Book) int {
    result := 0
    for _, item := range books {
        if result < item.ID {
            result = item.ID
        }
    }
    return result + 1
}

func GetPositionById(id int) int {
    var result int = -1
    for position, item := range books {
        if id == item.ID {
            return position
        }
    }
    return result
}

func CreateBook(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    var newBook Book
    _ = json.NewDecoder(request.Body).Decode(&newBook)
    db, err := sql.Open("sqlite3", "file:foo.db?_loc=auto")
    if err != nil {
        panic(err.Error())
    }
    insertion, err := db.Prepare("INSERT INTO books(id, isbn, title) VALUES(?,?,?)")
    if err != nil {
        panic(err.Error())
    }
    insertion.Exec(&newBook.ID, &newBook.ISBN, &newBook.Title)
    defer db.Close()
    response.WriteHeader(http.StatusOK)
    json.NewEncoder(response).Encode(newBook)
}

func UpdateBook(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    params := mux.Vars(request)
    expectedId, err := strconv.ParseInt(params["id"], 10, 32)
    if err != nil {
        http.Error(response, fmt.Sprintf( "malformed id: %v", params["id"]), 400)
        return 
    }
    var position int = GetPositionById(int(expectedId))
    if position != -1 {
        json.NewDecoder(request.Body).Decode(&books[position])
    } else {
        fmt.Printf("not found :%d\n", expectedId)
        http.NotFound(response, request)
        return 
    }
    json.NewEncoder(response).Encode(&books[position])
    return
}

func DeleteBook(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    params := mux.Vars(request)
    expectedId, err := strconv.ParseInt(params["id"], 10, 32)
    if err != nil {
        http.Error(response, fmt.Sprintf( "malformed id: %v", params["id"]), 400)
        return 
    }
    var position int = GetPositionById(int(expectedId))
    if position != -1 {
        response.WriteHeader(http.StatusOK)
        books = books[:position+copy(books[position:], books[position+1:])]
    } else {
        fmt.Printf("not found :%d\n", expectedId)
        http.NotFound(response, request)
    }
}

func main() {
    router := mux.NewRouter()
    

    books = append(books, Book{ID: 1, ISBN: "438227", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
    books = append(books, Book{ID: 2, ISBN: "454555", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})


    router.HandleFunc("/api/books", GetBooks).Methods("GET")
    router.HandleFunc("/api/book/{id}", GetBook).Methods("GET")
    router.HandleFunc("/api/books", CreateBook).Methods("POST")
    router.HandleFunc("/api/book/{id}", UpdateBook).Methods("PUT")
    router.HandleFunc("/api/book/{id}", DeleteBook).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":3000", router))
}

