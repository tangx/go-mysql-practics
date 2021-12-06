package main

import (
	"database/sql"
	"encoding/csv"
	"os"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/tangx/mysql-go-pratics/pkg/dbconn"
)

type News struct {
	id      int
	news_id int
	title   string
	link    string
}

var (
	newsTableSql = `
create table news(
	id int not null primary key auto_increment,
	news_id int not null,
	title varchar(128) not null,
	link varchar(256)
) default charset utf8;
`

	insertRecordSql = `
insert into news(news_id, title, link)
	values(?, ?, ?);
`
)

var (
	db      *sql.DB
	csvFile = "data.csv"
)

func init() {
	db = dbconn.InitDB()
}

func main() {
	// createTabel()

	for _, record := range readCsvFile() {
		if len(record) != 3 {
			logrus.Errorf("invalid news %s", strings.Join(record, ", "))
			continue
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			logrus.Errorf("invalid news id: %s, err: %v", record[0], err)
		}

		ret, err := db.Exec(insertRecordSql, id, record[1], record[2])
		if err != nil {
			logrus.Errorf("insert news failed: %v", err)
			continue
		}

		n, err := ret.RowsAffected()
		if err != nil {
			logrus.Errorf("Insert RowsAffect failed: %v", err)
			continue
		}

		logrus.Infof("Insert RowsAffectd: %d", n)
	}
}

func createTabel() {
	ret, err := db.Exec(newsTableSql)
	if err != nil {
		logrus.Fatalf("create table faield: %v", err)
	}

	n, err := ret.LastInsertId()
	if err != nil {
		logrus.Errorf("rows affected faield: %v", err)
		return
	}

	logrus.Infof("last insertid %d", n)
}
func readCsvFile() [][]string {
	f, err := os.Open(csvFile)
	if err != nil {
		logrus.Fatalf("open %s failed: %v", csvFile, err)
	}
	defer f.Close()

	csvobj := csv.NewReader(f)
	records, err := csvobj.ReadAll()
	if err != nil {
		logrus.Fatalf("read csvobj faield: %v", err)
	}

	return records
}
