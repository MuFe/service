package db

import (

	"github.com/DATA-DOG/go-sqlmock"
)

func GetTestDb()sqlmock.Sqlmock{
	testDb, mock, err := sqlmock.New()
	tDb := &Db{

	}
	tDb.ConnectTest(testDb)
	if err != nil {

	}
	SetTestDb(tDb)
	return mock
}


