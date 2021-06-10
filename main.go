package main

import (
	"embed"
	"fmt"
	"github.com/dstotijn/golang-nextjs-portable/database"
	"github.com/dstotijn/golang-nextjs-portable/models"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"time"
)

var (
	//go:embed nextjs/dist
	//go:embed nextjs/dist/_next
	//go:embed nextjs/dist/_next/static/chunks/pages/*.js
	//go:embed nextjs/dist/_next/static/*/*.js
	nextFS embed.FS
	dbCon *gorm.DB
	noPeople int64
)

func main() {
	log.Println("initial ...")

	// Root at the `dist` folder generated by the Next.js app.
	distFS, err := fs.Sub(nextFS, "nextjs/dist")
	if err != nil {
		log.Panic(err)
	}

	err = initial()
	if err != nil {
		log.Panic(err)
	}

	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: customErrorHandler,
	})

	// The static Next.js app will be served under `/`.
	app.Use(filesystem.New(filesystem.Config{
		Root: http.FS(distFS),
		MaxAge: 3600,
	}))

	// Default request id config
	app.Use(requestid.New())

	// The Memory allocation stats API will be served under `/api`.
	app.Get("/api", adaptor.HTTPHandlerFunc(handleAPI))

	app.Get("/all", allPeople)

	app.Post("/add", addPeople)

	err = gracefulStop(app)
	if err != nil {
		log.Panic(err)
	}

	log.Println("Starting HTTP server at http://localhost:8080 ...")
	if err := app.Listen(":8080"); err != nil {
		log.Panic(err)
	}
}

func initial() (err error) {
	log.Println("initial data ...")

	var rooPath = "/db"

	if _, err := os.Stat(rooPath); os.IsNotExist(err) {
		rooPath = "./db"
		err = os.Mkdir(rooPath, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	dbCon, err = database.DbOpen(rooPath + "/data_file.db")
	if err != nil {
		return err
	}

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

func customErrorHandler(c *fiber.Ctx, err error) error {
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

	log.Printf("%v %v %v\n", people.Id, people.Name, people.Surname)

	return c.Type("json", "utf-8").Status(fiber.StatusOK).JSON(people)
}

func gracefulStop(app *fiber.App) (err error) {

	c := make(chan os.Signal, 1)   // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt) // When an interrupt is sent, notify the channel

	// Goroutine to monitor the channel and run app. Shutdown when an interrupt is recieved
	// This should cause app.Listen to return nil, then allowing the cleanup tasks to be run.
	go func() {
		sig := <-c
		fmt.Printf("caught sig: %+v >> waiting for 2 second to finish processing\n", sig)
		time.Sleep(2 * time.Second)

		fmt.Println("shutdown Fiber app")
		_ = app.Shutdown()
	}()

	return nil
}