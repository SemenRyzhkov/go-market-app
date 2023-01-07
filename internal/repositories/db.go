package repositories

import (
	"database/sql"
	"log"
)

const (
	initUsersTableQuery = "" +
		"CREATE TABLE IF NOT EXISTS public.users (" +
		"id varchar(45) primary key, " +
		"login varchar(45) unique not null, " +
		"password varchar(45) not null" +
		")"
	initOrdersTableQuery = "" +
		"CREATE TABLE IF NOT EXISTS public.orders (" +
		"number int8 primary key, " +
		"status int2 not null, " +
		"accrual decimal, " +
		"uploaded_at timestamptz not null, " +
		"user_id varchar(45) references public.users (id)" +
		")"
	initWithdrawTableQuery = "" +
		"CREATE TABLE IF NOT EXISTS public.withdraw (" +
		"number int8 primary key, " +
		"sum decimal not null, " +
		"processed_at timestamptz not null, " +
		"user_id varchar(45) references public.users (id)" +
		")"
	setTimeZoneQuery = "" +
		"set timezone = 'Europe/Moscow'"
	createUserLoginIndex = "" +
		"CREATE INDEX IF NOT EXISTS user_login_index " +
		"ON public.users (login)"
)

var db *sql.DB

func InitDB(dbAddress string) (*sql.DB, error) {
	if db != nil {
		return db, nil
	}
	db, connectionErr := sql.Open("postgres", dbAddress)
	if connectionErr != nil {
		log.Println(connectionErr)
		return nil, connectionErr
	}
	createTableErr := createTableIfNotExists(db)
	if createTableErr != nil {
		log.Println(createTableErr)
		return nil, createTableErr
	}
	return db, nil
}

func createTableIfNotExists(db *sql.DB) error {
	_, createUserTableErr := db.Exec(initUsersTableQuery)
	if createUserTableErr != nil {
		return createUserTableErr
	}
	_, createOrderTableErr := db.Exec(initOrdersTableQuery)
	if createOrderTableErr != nil {
		return createUserTableErr
	}
	_, createWithdrawTableErr := db.Exec(initWithdrawTableQuery)
	if createWithdrawTableErr != nil {
		return createWithdrawTableErr
	}
	_, createIndexErr := db.Exec(createUserLoginIndex)
	if createIndexErr != nil {
		return createIndexErr
	}
	_, setTimeZoneErr := db.Exec(setTimeZoneQuery)
	if setTimeZoneErr != nil {
		return setTimeZoneErr
	}
	return nil
}
