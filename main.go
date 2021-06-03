package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"log"
	"net/http"
	"runtime/pprof"
	"strconv"
)

//go:embed nextjs/dist
//go:embed nextjs/dist/_next
//go:embed nextjs/dist/_next/static/chunks/pages/*.js
//go:embed nextjs/dist/_next/static/*/*.js
var nextFS embed.FS

func main() {
	// Root at the `dist` folder generated by the Next.js app.
	distFS, err := fs.Sub(nextFS, "nextjs/dist")
	if err != nil {
		log.Fatal(err)
	}

	initial()

	app := fiber.New()

	// The static Next.js app will be served under `/`.
	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(distFS),
	}))

	// The Memory allocation stats API will be served under `/api`.
	app.Get("/api", adaptor.HTTPHandlerFunc(handleAPI))

	app.Get("/all", allPeople)

	app.Post("/add", addPeople)

	// Start HTTP server at :8080.
	log.Println("Starting HTTP server at http://localhost:8080 ...")
	log.Fatal(app.Listen(":8080"))
}

var (
	db       *sql.DB
	stmt     *sql.Stmt
	noPeople int64
	lid      int64
)

type People struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func initial() (err error) {
	log.Println("initial data ...")

	db, err = sql.Open("sqlite3", "./date_file.db?_journal_mode=WAL&_synchronous=NORMAL&mode=rwc&cache=shared&_busy_timeout=20000")
	if err != nil {
		log.Println("database error")
		return err
	}

	db.SetMaxOpenConns(1)

	db.Exec("PRAGMA page_size= 65535;")
	db.Exec("PRAGMA cache_size= 8000;")
	db.Exec("PRAGMA mmap_size = 30000000000;")

	statement, _ := db.Prepare(`
        CREATE TABLE IF NOT EXISTS people (
            id INTEGER PRIMARY KEY,
            name TEXT,
            surname TEXT
        );
    `)

	_, err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS people_id_index ON people (id);")
	if err != nil {
		log.Printf("not create index. : %v", err)
	} else {
		log.Println("create index.")
	}

	statement.Exec()


	var count int64
	_ = db.QueryRow("SELECT COUNT(*) FROM people").Scan(&count)

	noPeople = count

	log.Printf("no. people: %v", count)

	if count == 0 {
		statement, _ = db.Prepare("INSERT INTO people (name, surname) VALUES (?, ?)")
		statement.Exec("Nic1", "Robert1")
		statement.Exec("Nic2", "Robert2")
		statement.Exec("Nic3", "Robert3")
	} else {
		log.Println("not initial people data.")
	}

	return nil
}

func allPeople(c *fiber.Ctx) error{
	log.Println("all people ...")

	rows, _ := db.Query("SELECT id, name, surname FROM people ORDER BY id DESC LIMIT 10")

	defer rows.Close()

	var result []People
	for rows.Next() {
		item := People{}

		rows.Scan(&item.Id, &item.Name, &item.Surname)

		result = append(result, item)

		var _ = strconv.FormatInt(item.Id, 10) + ": " + item.Name + " " + item.Surname
	}

	jsonB, _ := json.Marshal(result)

	return c.Type("json", "utf-8").Status(fiber.StatusOK).Send(jsonB)
}

func addPeople(c *fiber.Ctx) error{


	var body People
	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	noPeople++

	log.Printf("add people no.: %v\n", noPeople)

	stmt, _ = db.Prepare("INSERT INTO people (name, surname) VALUES (?, ?)")

	var res sql.Result
	res, _ = stmt.Exec(
		body.Name + strconv.FormatInt(noPeople, 10),
		body.Surname + strconv.FormatInt(noPeople, 10),
	)

	lid, _ = res.LastInsertId()

	rows, _ := db.Query("SELECT id, name, surname FROM people WHERE id = ?", strconv.FormatInt(lid, 10))

	defer rows.Close()

	var item = People{}
	if rows.Next() {
		rows.Scan(&item.Id, &item.Name, &item.Surname)
	}

	log.Printf("add %#v\n", item)

	jsonB, _ := json.Marshal(item)

	return c.Type("json", "utf-8").Status(fiber.StatusOK).Send(jsonB)
}


func handleAPI(w http.ResponseWriter, _ *http.Request) {
	// Gather memory allocations profile.
	profile := pprof.Lookup("allocs")

	// Write profile (human readable, via debug: 1) to HTTP response.
	err := profile.WriteTo(w, 1)
	if err != nil {
		log.Printf("Error: Failed to write allocs profile: %v", err)
	}
}