package coachModel

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
	"strings"
	"testing"
)

func TestGetInstitutionCoach(t *testing.T) {
	mock:=db.GetTestDb()

	sqlSelectSql := "select uid,info from qz_institution_coach where 1=1 "
	defer db.ClearTestDb()
	// 模拟 select 报错
	mock.ExpectQuery(sqlSelectSql).WillReturnError(errors.New("select error"))
	_, err := GetInstitutionCoach(2,3)

	if err != nil &&!strings.Contains(err.Error(), "select error"){
		t.Fatalf("unexpected error:%s",err)
	}

	// 模拟 select 正常
	rows := sqlmock.NewRows(
		[]string{"2", "1", "3"},
	).AddRow([]driver.Value{123, 1, "AABBCCDDEEFF"}...)
	mock.ExpectQuery(sqlSelectSql).WillReturnRows(rows)
	res, err := GetInstitutionCoach(2,3)
xlog.Info(res)
	if err != nil {
		xlog.ErrorP(err)
		t.Fatalf("unexpected error:%s",err)
	}


	if res[0].Uid != int64(123) {
		t.Fatalf("unexpected id:%d",res[0].Uid)
	}

	if res[0].InstitutionId != int64(1) {
		t.Fatalf("unexpected Hostname:%d",res[0].InstitutionId)
	}


	if res[0].Info != "AABBCCDDEEFF" {
		t.Fatalf("unexpected UID:%s",res[0].Info)
	}
}
