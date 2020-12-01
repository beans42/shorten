package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gofiber/fiber"
)

const (
	port         = ":80"
	databaseFile = "./database.json"
)

var urls map[string]string = make(map[string]string)

func randomHexString() (out string) {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	for _, v := range bytes {
		out += fmt.Sprintf("%02x", v)
	}
	return
}

func getLink(c *fiber.Ctx) error { // /:id
	text, ok := urls[c.Params("id")]
	if ok {
		return c.Redirect(text)
	}
	return c.SendString("invalid id")
}

func newLink(c *fiber.Ctx) error {
	var shortened string
	if len(c.FormValue("pref_id")) > 0 {
		shortened = c.FormValue("pref_id")
	} else {
		shortened = randomHexString()
	}

	urls[shortened] = c.FormValue("url")
	bytes, _ := json.Marshal(urls)
	defer ioutil.WriteFile(databaseFile, bytes, 0644)
	return c.SendString("shortened url: " + c.BaseURL() + "/" + shortened)
}

func main() {
	app := fiber.New()

	content, err := ioutil.ReadFile(databaseFile)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, &urls)
	if err != nil {
		log.Fatal(err)
	}

	app.Static("/", "./static")
	app.Get("/:id", getLink)
	app.Post("/", newLink)

	app.Listen(port)
}
