package db

import (
	"os"
	"strconv"
)

// db 数据库
const (
	USER    = "USER"    //用户db名称
	BANNER    = "BANNER"    //bannerDb名称
	COURSE    = "COURSE"    //courseDb名称
	LIVE    = "LIVE"    //liveDb名称
	SCHOOL    = "SCHOOL"    //schoolDb名称
	VIDEO    = "VIDEO"    //videoDb名称
	ADMIN    = "ADMINUSER"    //adminUserDb名称
	GOOD    = "GOOD"    //goodDb名称
	ORDER    = "ORDER"    //orderDb名称
	COACH    = "COACH"    //coachDb名称
)
// SQLPlaceholderLimit 占位符查询限制

var (
	dbMap map[string]*Db
)

type DbConfig struct {
	Host     string
	Port     int
	User     string
	Pwd      string
	Database string
}


func init() {
	dbMap = make(map[string]*Db, 0)
}

func getDb(name string) *Db {
	db,ok:=dbMap["Test"]
	if ok{
		return db
	}
	db, ok = dbMap[name]
	if !ok {
		port, err := strconv.Atoi(os.Getenv("MYSQL_" + name + "_PORT"))
		if err != nil {
			panic(err)
		}
		config:=DbConfig{
			Host:     os.Getenv("MYSQL_" + name + "_HOST"),
			Port:     port,
			User:     os.Getenv("MYSQL_" + name + "_USER"),
			Pwd:      os.Getenv("MYSQL_" + name + "_PASSWORD"),
			Database: os.Getenv("MYSQL_" + name + "_DATABASE"),
		}
		db, err = ConnectDb(config)
		if err != nil {
			panic(err)
		}
		dbMap[name] = db
	}
	return db
}

func ConnectDb(config DbConfig) (*Db, error) {
	db := &Db{}
	err := db.Connect(config.Host, config.Port, config.User, config.Pwd, config.Database)
	return db, err
}

func GetUserDb() *Db {
	return getDb(USER)
}

func GetBannerDb() *Db {
	return getDb(BANNER)
}

func GetCourse() *Db {
	return getDb(COURSE)
}

func GetSchool() *Db {
	return getDb(SCHOOL)
}

func GetVideo() *Db {
	return getDb(VIDEO)
}
func GetAdminDb() *Db {
	return getDb(ADMIN)
}

func GetGoodDb() *Db {
	return getDb(GOOD)
}


func GetOrderDb() *Db {
	return getDb(ORDER)
}


func GetLiveDb() *Db {
	return getDb(LIVE)
}


func GetCoachDb() *Db {
	return getDb(COACH)
}

func SetTestDb(db *Db){
	dbMap["Test"] = db
}

func ClearTestDb(){
	db,ok:=dbMap["Test"]
	if ok{
		db.db.Close()
		dbMap["Test"]=nil
	}

}
