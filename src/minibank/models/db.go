package models

import (
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocql/gocql"
)

var Database *sql.DB

var CassandraSession *gocql.Session
var CassandraEnabled bool

func init() {
	CassandraEnabled = cassandraEnabled()
}

func InitDB(dbDone_chan chan<- bool) {
	var err error
	Database, err = sql.Open("mysql", dbConn())
	if err != nil {
		log.Panic(err)
	}
	for {
		err = Database.Ping()
		if err != nil {
			log.Print(err)
			time.Sleep(1 * time.Second)
		} else {
			log.Print("Successfully connected to Database!")
			dbDone_chan <- true
			break
		}
	}
}

// dbConn looks up database connection string on environment
func dbConn() string {
	connectString := os.Getenv("DB_CONNECTION_STRING")
	if len(connectString) == 0 {
		connectString = "minibank:minibank@tcp(mysql)/minibank"
	}
	return connectString
}

func InitCassandra() error {
	cluster := gocql.NewCluster(getCassandraHost())
	pass := gocql.PasswordAuthenticator{"minibank", "minibank"}
	cluster.Authenticator = pass
	cluster.Consistency = gocql.One
	sess, err := cluster.CreateSession()
	for err != nil {
		time.Sleep(time.Second)
		sess, err = cluster.CreateSession()
	}
	log.Print("Got connected to Cassandra")
	CassandraSession = sess

	kstmt := "CREATE KEYSPACE IF NOT EXISTS minibank WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor' : 3}"
	err = CassandraSession.Query(kstmt).Exec()
	if err != nil {
		log.Fatal("Unable to create minibank keyspace")
		return err
	}
	tstmt := "CREATE TABLE IF NOT EXISTS minibank.sessions (session text, username text, expiration bigint, PRIMARY KEY(session))"
	err = CassandraSession.Query(tstmt).Exec()
	if err != nil {
		log.Fatal("Unable to create sessions table in minibank keyspace")
	}
	istmt := "CREATE INDEX IF NOT EXISTS ON minibank.sessions(username);"
	err = CassandraSession.Query(istmt).Exec()
	if err != nil {
		log.Fatal("Unable to create sessions index on username in minibank keyspace")
	}
	cluster.Keyspace = "minibank"
	return err
}

func cassandraEnabled() bool {
	enableCassandra := os.Getenv("ENABLE_CASSANDRA")
	if len(enableCassandra) != 0 {
		if strings.ToLower(enableCassandra) == "true" {
			return true
		}
	}
	return false
}

func getCassandraHost() string {
	chost := os.Getenv("CASSANDRA_HOST")
	if len(chost) != 0 {
		return chost
	}
	return "cassandra"
}
