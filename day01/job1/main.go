package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tangx/mysql-go-pratics/pkg/dbconn"
)

var db *sql.DB

const (
	createTableSql = `
create table userinfo (
	id int not null auto_increment primary key,
	name varchar(32) not null,
	password varchar(64) not null,
	gender enum('male','female') not null,
	email varchar(64) null,
	amount decimal(10,2) not null default 0,
	ctime datetime
) default charset=utf8;
`

	insertUserSql = `
insert into userinfo(name, password, gender, email, amount, ctime)
	values(?, ?, ?, ?, ?, ?);
`

	setGenderToMaleIDgtNSql = `
update userinfo set gender='male' where id > ?;
`

	selectUserAmountGtNSql = `
select name,amount from userinfo where amount > ?;
`

	increaseAmountSql = `
	update userinfo set amount=amount + ? ;
`

	deleteAllMaleSql = `
	delete from userinfo where gender='male';	
`
)

type User struct {
	id       int
	name     string
	password string
	gender   string
	email    string
	amount   decimal.Decimal
	ctime    time.Time
}

func main() {
	// createTable()
	// insertUser()
	// setGenderToMaleIDgtN(3)
	// selectUserAmountGtN(1000)
	// increaseAmount(1000)
	deleteAllMale()
}

// 06. 删除所有男性角色
func deleteAllMale() {

	ret, err := db.Exec(deleteAllMaleSql)
	if err != nil {
		logrus.Errorf("delete male failed: %v", err)
		return
	}

	nn, err := ret.RowsAffected()
	if err != nil {
		logrus.Errorf("delete male RowAffected failed: %v", err)
		return
	}

	logrus.Infof("delete male RowAffected %d", nn)
}

// 05. 所有人工资涨 1000
func increaseAmount(n int) {
	ret, err := db.Exec(increaseAmountSql, n)
	if err != nil {
		logrus.Errorf("increase amount failed: %v", err)
		return
	}

	nn, err := ret.RowsAffected()
	if err != nil {
		logrus.Errorf("increase amount RowAffected failed: %v", err)
		return
	}

	logrus.Infof("increase amount RowAffected %d", nn)
}

// 04. 查询所有工资大于 1000 的人
func selectUserAmountGtN(n int) {
	rows, err := db.Query(selectUserAmountGtNSql, n)
	if err != nil {
		logrus.Fatalf("select user amount > 1000 failed: %v", err)
	}
	// 关闭 rows 释放连接池
	defer rows.Close()

	// 循环读取结果
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.name, &user.amount)
		if err != nil {
			logrus.Fatalf("select user scan failed: %v", err)
		}

		logrus.Infof("select users: %s, %s", user.name, user.amount)
	}
}

// 03. 将所有 id > 3 的人设置为 male
func setGenderToMaleIDgtN(n int) {
	ret, err := db.Exec(setGenderToMaleIDgtNSql, 3)
	if err != nil {
		logrus.Errorf("set gender to male if id > 3 faield: %v", err)
	}

	nn, err := ret.RowsAffected()
	if err != nil {
		logrus.Errorf("update gender to male RowsAffect failed: %v", err)
	}

	logrus.Infof("update gender to male RowAffectd %d", nn)
}

// 02. 插入用户
func insertUser() {
	for _, user := range []User{
		{
			name:     "zhangsan",
			password: "zhangsan123",
			gender:   "male",
			email:    "zhangsan@example.com",
			amount:   decimal.NewFromFloat(123.43),
			ctime:    time.Now(),
		},
		{
			name:     "muwanqing",
			password: "muwanqing123",
			gender:   "female",
			email:    "muwanqing@example.com",
			amount:   decimal.NewFromFloat(2131231.43),
			ctime:    time.Now(),
		},
		{
			name:     "ningrongrong",
			password: "ningrongrong123",
			gender:   "female",
			email:    "ningrongrong@example.com",
			amount:   decimal.NewFromInt(0),
			ctime:    time.Now(),
		},
		{
			name:     "mahongjun",
			password: "ningrongrong123",
			gender:   "female",
			email:    "mahongjun@example.com",
			amount:   decimal.NewFromFloat(211231.43),
			ctime:    time.Now(),
		},
		{
			name:     "zhugeliang",
			password: "zhuge123",
			gender:   "male",
			email:    "zhuge@example.com",
			amount:   decimal.NewFromFloat(1234.43),
			ctime:    time.Now(),
		},
		{
			name:     "zhangfei",
			password: "zhagnfei123",
			gender:   "male",
			email:    "zhangfei@example.com",
			amount:   decimal.NewFromFloat(0.43),
			ctime:    time.Now(),
		},
	} {
		ret, err := db.Exec(insertUserSql,
			user.name, user.password, user.gender,
			user.email, user.amount, user.ctime,
		)
		if err != nil {
			logrus.Errorf("ineruser failed: %v", err)
			continue
		}

		n, err := ret.RowsAffected()
		if err != nil {
			logrus.Errorf("insert user rowsAffect failed: %v", err)
		}

		logrus.Infof("insert user %s suceess: rows affected %d", user.name, n)
	}
}

// 01. 创建表
func createTable() {
	ret, err := db.Exec(createTableSql)
	if err != nil {
		logrus.Fatalf("create table failed: %v", err)
	}

	n, err := ret.RowsAffected()
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Printf("create table affectd rows: %d", n)
}

func initDB() *sql.DB {
	dsn := `root:Mysql12345@tcp(127.0.0.1:3306)/mysqlpractics`
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logrus.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}

func init() {
	db = dbconn.InitDB()
}
