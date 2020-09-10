package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Siswa struct {
	Id     int
	Nama   string
	Jk     string
	Alamat string
}

var tmpl, err = template.ParseGlob("views/*")

func connectDb() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "gocrud"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		db := connectDb()
		selDb, err := db.Query("SELECT * FROM siswa ORDER BY id DESC")
		if err != nil {
			panic(err.Error())
		}

		siswa := Siswa{}
		res := []Siswa{}

		for selDb.Next() {
			var id int
			var nama, jk, alamat string
			err = selDb.Scan(&id, &nama, &jk, &alamat)
			if err != nil {
				panic(err.Error())
			}

			siswa.Id = id
			siswa.Nama = nama
			siswa.Jk = jk
			siswa.Alamat = alamat
			res = append(res, siswa)
		}

		tmpl.ExecuteTemplate(w, "Index", res)
		defer db.Close()
	})

	http.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "New", nil)
	})

	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		db := connectDb()
		if r.Method == "POST" {
			nama := r.FormValue("nama")
			jk := r.FormValue("jk")
			alamat := r.FormValue("alamat")

			_, err := db.Exec("INSERT INTO siswa VALUES (?, ?, ?, ?)", nil, nama, jk, alamat)
			if err != nil {
				panic(err.Error())
			}

			log.Println("Sukses")
		}

		defer db.Close()
		http.Redirect(w, r, "/", 301)
	})

	http.HandleFunc("/edit", func(w http.ResponseWriter, r *http.Request) {
		db := connectDb()
		id := r.URL.Query().Get("id")
		data, err := db.Query("SELECT * FROM siswa WHERE id = ? ", id)
		if err != nil {
			panic(err.Error())
		}

		siswa := Siswa{}
		for data.Next() {
			var id int
			var nama, jk, alamat string
			err = data.Scan(&id, &nama, &jk, &alamat)
			if err != nil {
				panic(err.Error())
			}

			siswa.Id = id
			siswa.Nama = nama
			siswa.Jk = jk
			siswa.Alamat = alamat
		}

		tmpl.ExecuteTemplate(w, "Edit", siswa)
		defer db.Close()
	})

	http.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		db := connectDb()

		if r.Method == "POST" {
			id := r.FormValue("id")
			nama := r.FormValue("nama")
			jk := r.FormValue("jk")
			alamat := r.FormValue("alamat")

			_, err := db.Exec("UPDATE siswa SET nama = ?, jk = ?, alamat = ? WHERE id = ?", nama, jk, alamat, id)

			if err != nil {
				panic(err.Error())
			}

			log.Println("Sukses update")
		}

		defer db.Close()
		http.Redirect(w, r, "/", 301)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		db := connectDb()

		var id string = r.URL.Query().Get("id")
		_, err := db.Exec("DELETE FROM siswa WHERE id = ? ", id)
		if err != nil {
			panic(err.Error())
		}

		log.Println("Sukses hapus")

		defer db.Close()
		http.Redirect(w, r, "/", 301)
	})

	fmt.Println("server started at localhost:8000")
	http.ListenAndServe(":8000", nil)
}
