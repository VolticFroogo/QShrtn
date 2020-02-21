package main

// The way I imported the SQL database into MongoDB.
// This should never need to be used again, but in case it does, here is the function.

// import (
// 	"context"
// 	"database/sql"
// 	"log"
//
// 	"github.com/VolticFroogo/QShrtn/db"
// 	"github.com/VolticFroogo/QShrtn/model"
// 	_ "github.com/go-sql-driver/mysql" // MySQL Driver.
// 	"go.mongodb.org/mongo-driver/bson"
// )
//
// func importFromSQL() {
// 	ctx := context.Background()
//
// 	// Open the SQL connection.
// 	SQL, err := sql.Open("mysql", "root:development@tcp(localhost:3306)/qshrtn")
// 	if err != nil {
// 		log.Print(err)
// 		return
// 	}
//
// 	rows, err := SQL.Query("SELECT * FROM qshrtn.redirect")
// 	if err != nil {
// 		log.Print(err)
// 		return
// 	}
//
// 	var items []interface{}
//
// 	for rows.Next() {
// 		var redirect model.Redirect
//
// 		err = rows.Scan(&redirect.ID, &redirect.URL)
// 		if err != nil {
// 			log.Print(err)
// 			return
// 		}
//
// 		items = append(items, bson.M{
// 			"_id": redirect.ID,
// 			"url": redirect.URL,
// 		})
// 	}
//
// 	_, err = db.Redirect.InsertMany(ctx, items)
// 	if err != nil {
// 		log.Print(err)
// 		return
// 	}
// }
