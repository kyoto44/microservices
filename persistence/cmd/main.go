package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

const (
	portHTTP               = "8001"
	postgresConnect        = "postgresql://postgres:SuperStrongPassword@db:5432/postgres"
	insertStatement string = `INSERT INTO logs (uuid, user_agent, ip, country, timestamp) VALUES ($1, $2, $3, $4, $5)`
)

type Database struct {
	connect *pgx.Conn
	mutex   *sync.Mutex
}

type RandomData struct {
	ID          string       `json:"id"`
	Fingerprint *Fingerprint `json:"fingerprint"`
}

type Fingerprint struct {
	UserAgent string `json:"useragent"`
	IP        string `json:"ip"`
	Country   string `json:"country"`
}

var (
	db *Database
)

func (db *Database) CreateTable() error {
	_, err := db.connect.Exec(context.Background(),
		`create table if not exists logs
		(
			uuid character varying(36),
			user_agent character varying(256),
			ip character varying(16),
			country character varying(3),
			"timestamp" timestamp without time zone
		)`,
	)
	return err
}

func main() {

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	pgDB, err := pgx.Connect(context.Background(), postgresConnect)
	if err != nil {
		log.Fatal(err)
	}
	defer pgDB.Close(context.Background())

	db = &Database{
		connect: pgDB,
		mutex:   &sync.Mutex{},
	}

	err = db.CreateTable()
	if err != nil {
		log.Fatal(err)
	}

	app.Get("/internal/getrandom", func(c *fiber.Ctx) error {
		return c.JSON(GetRandomDBEntry())
	})

	app.Get("/internal/getall", func(c *fiber.Ctx) error {
		return c.JSON(GetAllDBEntries())
	})

	app.Get("/internal/deleteall", func(c *fiber.Ctx) error {
		err = db.DeleteAllEntries()
		if err != nil {
			return c.SendStatus(fasthttp.StatusInternalServerError)
		}
		return c.SendStatus(fasthttp.StatusOK)
	})

	app.Post("post", func(c *fiber.Ctx) error {
		var msgJSON RandomData
		err := json.Unmarshal(c.Body(), &msgJSON)
		if err != nil {
			log.Errorf("Can't deserialize data, reason: %s\n", err)
			return c.SendStatus(fasthttp.StatusBadRequest)
		}

		db.mutex.Lock()
		_, err = db.connect.Exec(context.Background(), insertStatement, msgJSON.ID, msgJSON.Fingerprint.UserAgent,
			msgJSON.Fingerprint.IP, msgJSON.Fingerprint.Country, time.Now().UTC())
		if err != nil {
			log.Error(err)
		}
		db.mutex.Unlock()

		return c.SendStatus(fasthttp.StatusOK)
	})

	log.Info("Persistence started successfully")
	app.Listen(":" + portHTTP)
}

func (db *Database) GetEntries() ([]byte, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	var buf []byte
	err := db.connect.QueryRow(
		context.Background(),
		"select jsonb_agg(jsonb_build_object('useragent', user_agent, 'ip', ip, 'country', country)) from logs",
	).Scan(&buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func GetAllDBEntries() *[]Fingerprint {
	rawData, err := db.GetEntries()
	if err != nil {
		log.Error(err)
	}

	var dbdata []Fingerprint
	err = json.Unmarshal(rawData, &dbdata)
	if err != nil {
		log.Error(err)
	}

	if len(dbdata) > 0 {
		return &dbdata
	}

	return &[]Fingerprint{}
}

func GetRandomDBEntry() *Fingerprint {
	rawData, err := db.GetEntries()
	if err != nil {
		log.Error(err)
	}

	var dbdata []Fingerprint
	err = json.Unmarshal(rawData, &dbdata)
	if err != nil {
		log.Error(err)
	}

	if len(dbdata) > 0 {
		return &dbdata[rand.Intn(len(dbdata))]
	}

	return &Fingerprint{}

}

func (db *Database) DeleteAllEntries() error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	_, err := db.connect.Exec(context.Background(), "delete from logs")
	if err != nil {
		return err
	}

	return nil
}
