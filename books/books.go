package books
import (
    "fmt"
)
// Book

func init() {
    fmt.Print("books")
}

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

type Books struct {
    items[] Book `json:"books"`
}
