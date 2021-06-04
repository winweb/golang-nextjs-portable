package main

import (
	"embed"
	"github.com/dstotijn/golang-nextjs-portable/database"
	"github.com/dstotijn/golang-nextjs-portable/models"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"gorm.io/gorm"
	"io/fs"
	"log"
	"net/http"
	"runtime/pprof"
	"strconv"
)

var (
	//go:embed nextjs/dist
	//go:embed nextjs/dist/_next
	//go:embed nextjs/dist/_next/static/chunks/pages/*.js
	//go:embed nextjs/dist/_next/static/*/*.js
	nextFS embed.FS
	dbCon *gorm.DB
	noPeople int64
	lid      int64
)

func main() {
	// Root at the `dist` folder generated by the Next.js app.
	distFS, err := fs.Sub(nextFS, "nextjs/dist")
	if err != nil {
		log.Fatal(err)
	}

	initial()

	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Default 500 status code
			code := fiber.StatusInternalServerError

			if e, ok := err.(*fiber.Error); ok {
				// Override status code if fiber.Error type
				code = e.Code
			}
			
			// Set Content-Type: text/plain; charset=utf-8
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

			log.Printf("error code: %v, %v", code, err.Error())

			// Return status code with error message
			return c.Status(code).SendString(err.Error())
		},
	})

	// The static Next.js app will be served under `/`.
	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(distFS),
		MaxAge: 3600,
	}))

	// The Memory allocation stats API will be served under `/api`.
	app.Get("/api", adaptor.HTTPHandlerFunc(handleAPI))

	app.Get("/all", allPeople)

	app.Post("/add", addPeople)

	// Start HTTP server at :8080.
	log.Println("Starting HTTP server at http://localhost:8080 ...")
	log.Fatal(app.Listen(":8080"))
}

func initial() (err error) {
	log.Println("initial data ...")

	dbCon, err = database.DbOpen("data_file.db")

	var count int64
	dbCon.Find(&models.People{}).Count(&count)

	noPeople = count

	log.Printf("no. people: %v", count)

	if count == 0 {
		dbCon.Create(&models.People{Name: "Nic1", Surname: "Robert1"})
		dbCon.Create(&models.People{Name: "Nic2", Surname: "Robert2"})
		dbCon.Create(&models.People{Name: "Nic3", Surname: "Robert3"})

	} else {
		log.Println("not initial people data.")
	}

	return nil
}

func allPeople(c *fiber.Ctx) error{
	log.Println("all people ...")

	var peoples []models.People

	if err := dbCon.Limit(10).Order("id desc").Find(&peoples).Error; err != nil {
		return err
	}

	return c.Type("json", "utf-8").Status(fiber.StatusOK).JSON(peoples)
}

func addPeople(c *fiber.Ctx) error{

	people := new(models.People)

	if err := c.BodyParser(people); err != nil {
		log.Printf("error: add people. %v\n", err.Error())
		return err
	}

	noPeople++

	log.Printf("add people no.: %v\n", noPeople)

	people.Name = people.Name + strconv.FormatInt(noPeople, 10)
	people.Surname = people.Surname + strconv.FormatInt(noPeople, 10)

	if err := dbCon.Create(&people).Error; err != nil {
		log.Printf("error: create people. %v\n", err.Error())
		return err
	}

	return c.Type("json", "utf-8").Status(fiber.StatusOK).JSON(people)
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