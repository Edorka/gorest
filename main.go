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
    "gorest/books"
)



type BooksApp struct {
    Access *sql.DB
}

func (app *BooksApp) Source(source *sql.DB) {
    app.Access = source
}

func (app *BooksApp) List(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    selection, err := books.List(app.Access)
    if err != nil {
        panic(err.Error())
    }
    result := json.NewEncoder(response)
    for selection.Next() {
        var current books.Book = books.Book{} 
        err = selection.Scan(&current.ID, &current.ISBN, &current.Title)
        if err != nil {
            panic(err.Error())
        }
        result.Encode(current)
    }
}

func (app *BooksApp) Get(response http.ResponseWriter, request *http.Request) {
    var expectedId int64
    var err error
    var result books.Book
    response.Header().Set("Content-Type", "application/json")
    params := mux.Vars(request)
    expectedId, err = strconv.ParseInt(params["id"], 10, 32)
    if err != nil {
        panic(err.Error())
    }
    result, err = books.Get(app.Access, int(expectedId))
    json.NewEncoder(response).Encode(result)
}

func (app *BooksApp) Create(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    var result, newBook books.Book
    var err error
    err = json.NewDecoder(request.Body).Decode(&newBook)
    if err != nil {
        panic(err.Error())
    }
    result, err = books.Create(app.Access, newBook)
    response.WriteHeader(http.StatusOK)
    json.NewEncoder(response).Encode(result)
}

func (app *BooksApp) Update(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    params := mux.Vars(request)
    expectedId, err := strconv.ParseInt(params["id"], 10, 32)
    if err != nil {
        http.Error(response, fmt.Sprintf( "malformed id: %v", params["id"]), 400)
        return 
    }
    var target, result books.Book
    _ = json.NewDecoder(request.Body).Decode(&target)
    result, err = books.Update(app.Access, int(expectedId), target)
    response.WriteHeader(http.StatusOK)
    json.NewEncoder(response).Encode(result)
}

func (app *BooksApp) Delete(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("Content-Type", "application/json")
    params := mux.Vars(request)
    expectedId, err := strconv.ParseInt(params["id"], 10, 32)
    if err != nil {
        http.Error(response, fmt.Sprintf( "malformed id: %v", params["id"]), 400)
        return 
    }
    var done bool;
    done, err = books.Delete(app.Access, int(expectedId))
    response.WriteHeader(http.StatusOK)
    json.NewEncoder(response).Encode(done)
}

func main() {
    router := mux.NewRouter()

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

