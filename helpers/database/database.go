package database

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

var ApiKey = os.Getenv("API_KEY")

// var Api = "http://ariesautoapi.curtmfg.com/part?key=" + apiKey
var Api = "http://localhost:8080/"

func InitMongo() error {
	var err error
	if MongoSession == nil {
		connectionString := MongoConnectionString()
		MongoSession, err = mgo.DialWithInfo(connectionString)
		if err == nil {
			MongoDatabase = connectionString.Database
		}
	}
	return err
}

var (
	MongoDatabase string
	MongoSession  *mgo.Session
)

func MongoConnectionString() *mgo.DialInfo {
	var info mgo.DialInfo

	addresses := []string{"127.0.0.1"}
	if hostString := os.Getenv("MONGO_URL"); hostString != "" {
		addresses = strings.Split(hostString, ",")
	}
	info.Addrs = addresses

	info.Username = os.Getenv("MONGO_USERNAME")
	info.Password = os.Getenv("MONGO_PASSWORD")
	info.Database = os.Getenv("MONGO_DATABASE")
	info.Timeout = time.Second * 2
	if info.Database == "" {
		info.Database = "DataMigration"
	}
	info.Source = "admin"

	return &info
}

func NewDBConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("CURT_DEV_NAME")
		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}

	return "root:@tcp(127.0.0.1:3306)/CurtDev2?parseTime=true&loc=America%2FChicago"
}

func OldDBConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("CURT_DEV_NAME")
		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}

	return "root:@tcp(127.0.0.1:3306)/CurtAriesDev?parseTime=true&loc=America%2FChicago"
}

func VcdbConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("VCDB_NAME")
		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}

	return "root:@tcp(127.0.0.1:3306)/vcdb?parseTime=true&loc=America%2FChicago"
}

func PcdbConnectionString() string {
	if addr := os.Getenv("DATABASE_HOST"); addr != "" {
		proto := os.Getenv("DATABASE_PROTOCOL")
		user := os.Getenv("DATABASE_USERNAME")
		pass := os.Getenv("DATABASE_PASSWORD")
		db := os.Getenv("PCDB_NAME")
		return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s", user, pass, proto, addr, db, "America%2FChicago")
	}

	return "root:@tcp(127.0.0.1:3306)/pcdb?parseTime=true&loc=America%2FChicago"
}
