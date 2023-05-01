package links

import (
	"context"
	"log"

	"github.com/captain-corgi/go-graphql-example/internal/hackernews/users"
	database "github.com/captain-corgi/go-graphql-example/pkg/db/mysql"
)

// #1
type Link struct {
	ID      string
	Title   string
	Address string
	User    *users.User
}

// #2
func (link *Link) Save(ctx context.Context) int64 {
	//#3
	stmt, err := database.Db.Prepare("INSERT INTO Links(Title,Address, UserID) VALUES(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	//#4
	res, err := stmt.ExecContext(ctx, link.Title, link.Address, link.User.ID)
	if err != nil {
		log.Fatal(err)
	}
	//#5
	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal("Error:", err.Error())
	}
	log.Print("Row inserted!")
	return id
}

func GetAll(ctx context.Context) []Link {
	stmt, err := database.Db.Prepare(`select L.id, L.title, L.address, L.UserID, U.Username 
	from Links L inner join Users U on L.UserID = U.ID`)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var (
		links    []Link
		username string
		id       string
	)
	for rows.Next() {
		var link Link
		err := rows.Scan(&link.ID, &link.Title, &link.Address, &id, &username)
		if err != nil {
			log.Fatal(err)
		}
		link.User = &users.User{
			ID:       id,
			Username: username,
		}
		links = append(links, link)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return links
}
