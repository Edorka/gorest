package books

import (
    "database/sql"
)

func List(source *sql.DB) (*sql.Rows, error) {
    const query = `SELECT id,isbn,title from books`
    selection, err := source.Query(query)
    return selection, err
}

func Get(source *sql.DB, expectedId int) (Book, error) {
    const query = `SELECT id,isbn,title from books where id = $1 `
    var result Book = Book{} 
    err := source.QueryRow(query, expectedId).Scan(&result.ID, &result.ISBN, &result.Title)
    return result, err
}

func Create(source *sql.DB, newBook Book) (Book, error) {
    const query = "INSERT INTO books(id, isbn, title) VALUES(?,?,?)"
    insertion, err := source.Prepare(query)
    if err != nil {
        panic(err.Error())
    }
    insertion.Exec(&newBook.ID, &newBook.ISBN, &newBook.Title)
    return newBook, err
}

func Update(source *sql.DB, id int, target Book) (Book, error) {
    const query = "UPDATE books SET isbn=?, title=?  where id = ?"
    update, err := source.Prepare(query)
    if err != nil {
        panic(err.Error())
    }
    update.Exec(&target.ISBN, &target.Title, id)
    return target, err
}


func Delete(source *sql.DB, id int) (bool, error) {
    const query = "DELETE FROM books where id = ?"
    update, err := source.Prepare(query)
    if err != nil {
        panic(err.Error())
    }
    update.Exec(id)
    const done = true
    return done, err
}
