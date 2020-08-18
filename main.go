package main

import (
	"fmt"
	phonedb "normalizer/db"
	"regexp"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgresql"
	dbname   = "phone_normalizer"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)

	must(phonedb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(phonedb.Migrate("postgres", psqlInfo))

	db, err := phonedb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	must(err)
	defer db.Close()

	err = db.Seed()
	must(err)

	phones, err := db.AllPhones()
	must(err)

	for _, p := range phones {
		fmt.Printf("Working on %+v\n", p)
		number := normalize(p.Number)
		if number != p.Number {
			existing, err := db.Findphone(number)
			must(err)
			if existing != nil {
				must(db.Deletephone(p.ID))
				//delete
			} else {
				p.Number = number
				must(db.Updatephone(&p))
				//update
			}
		} else {
			fmt.Println("no change")
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}


func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}

/*
func normalize(phone string) string {
	var buf bytes.Buffer
	for _, ch := range phone {
		if ch >= '0' && ch <= '9' {
			buf.WriteRune(ch)
		}
	}
	return buf.String()
}
*/
