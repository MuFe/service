package adminUserModel

import (
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
)

type Commitment struct {
	Id     int64
	Name   string
	Status int64
}

//获取承诺服务列表
func GetCommitment() (result []*Commitment, err error) {
	rows, err := db.GetAdminDb().Query("select id,name,status from qz_good_commitment")
	if err != nil {
		return nil, xlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, status int64
		var name string
		err := rows.Scan(&id, &name, &status)
		if err != nil {
			return nil, xlog.Error(err)
		}
		result = append(result, &Commitment{Id: id, Name: name, Status: status})
	}
	err = rows.Err()
	if err != nil {
		return result, xlog.Error(err)
	}
	return result, nil
}
