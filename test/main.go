package main

import (
	"fmt"
	"github.com/thaianhsoft/go-ghost/edge"
	"github.com/thaianhsoft/go-ghost/mockdb"
	"github.com/thaianhsoft/go-ghost/test/book"
	"github.com/thaianhsoft/go-ghost/test/note"
	"github.com/thaianhsoft/go-ghost/test/student"
)

func main() {
	db := mockdb.OpenDB()
	createQuery := edge.DefaultManager().Migrate(&book.Book{}, &student.Student{}, &note.Note{})
	if _, err := db.Query(*createQuery); err == nil {
		fmt.Println(*createQuery)
		fmt.Println("create successfully !")
	}
}
