package orderModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
	"strings"
	"time"
)

type AddressData struct {
	Id int64
	Name string
	Phone string
	Address string
	Area string
	City string
	Province string
	IsDefault bool

}



func GetUserAddress(uid int64, isFirst bool, aID int64) ([]*AddressData,error) {
	var result *sql.Rows
	var err error
	var buf strings.Builder
	args := make([]interface{}, 1)
	buf.WriteString(`SELECT id,NAME,phone,province,city,area,address,is_default FROM qz_user_address WHERE `)
	if aID == 0 {
		buf.WriteString(" uid=? ")
		if isFirst {
			buf.WriteString(" ORDER BY is_default DESC,modify_time desc limit 1")
		} else {
			buf.WriteString(" ORDER BY modify_time desc")
		}
		args[0] = uid
	} else {
		args[0] = aID
		buf.WriteString(" id=? ")
	}
	result, err =db.GetGoodDb().Query(buf.String(), args...)
	list := make([]*AddressData, 0)
	if err != nil {
		return nil, err
	}
	var name, phone, province, city, area, address sql.NullString
	var isDefault sql.NullBool
	var id sql.NullInt64
	for result.Next() {
		err = result.Scan(&id, &name, &phone, &province, &city, &area, &address, &isDefault)
		if err == nil {
			list = append(list, &AddressData{
				Id:        id.Int64,
				Name:      name.String,
				Phone:     phone.String,
				Address:   address.String,
				Area:      area.String,
				City:      city.String,
				Province:  province.String,
				IsDefault: isDefault.Bool,
			})
		}
	}
	return list,nil
}



func EditUserAddress(id,uid int64,name,phone,province,city,area,address string,del,isDefault bool)error{
	if del{
		_, err := db.GetGoodDb().Exec("delete from qz_user_address where id=?", id)
		if err != nil {
			return xlog.Error(err)
		}
		return  nil
	}else{
		if isDefault {
			_, err := db.GetGoodDb().Exec("update qz_user_address set is_default=0 where uid=?", uid)
			if err != nil {
				return xlog.Error(err)
			}
		}
		var deInt int64
		if isDefault {
			deInt = 1
		} else {
			deInt = 0
		}
		var err error
		if id != 0 {
			_, err = db.GetGoodDb().Exec(
				"update qz_user_address set name=?,phone=?,province=?,city=?,area=?,address=?,is_default=?,modify_time=? where id=?",
				name, phone, province,city,area,address, deInt, time.Now().Unix(), id)
		} else {
			_, err = db.GetGoodDb().Exec(
				"insert into qz_user_address (name,phone,province,city,area,address,is_default,uid,modify_time) values (?,?,?,?,?,?,?,?,?)",
				name, phone, province,city,area,address, deInt, uid, time.Now().Unix())
		}
		if err != nil {
			return xlog.Error(err)
		}
		return  nil
	}
}
