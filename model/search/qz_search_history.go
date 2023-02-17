package searchModel

import (
	"database/sql"
	"mufe_service/camp/db"
	"mufe_service/camp/xlog"
	"time"
)

type SearchData struct {
	Id      int64
	Number    int64
	TodayNumber    int64
	Content string
}

func AddSearch(uid int64, content string) error {
	return db.GetCourse().WithTransaction(func(tx *db.Tx) error {
		var id, modifyTime, number,todayNumber sql.NullInt64
		err := tx.QueryRow("select id,modify_time,number,today_number from qz_search_history where uid=? and content=?", uid, content).Scan(&id, &modifyTime, &number,&todayNumber)
		if err != nil && err != sql.ErrNoRows {
			return xlog.Error(err)
		} else if err == sql.ErrNoRows {
			_, err = tx.Exec("insert into qz_search_history (number,today_number,modify_time,uid,content) values (?,?,?,?,?)", 1,1, time.Now().Unix(), uid, content)
			return err
		} else {
			n := number.Int64
			n1:=todayNumber.Int64
			n1=n1+1
			if (modifyTime.Int64/86400)*86400 == (time.Now().Unix()/86400)*86400 {
				n = n + 1
			} else {
				n = 1
			}
			_, err = tx.Exec("update qz_search_history set number=?,today_number=?,modify_time=? where id=?",n1, n, time.Now().Unix(), id.Int64)
			return err
		}
	})
}


//返回搜索记录和热门搜索记录
func GetSearchHistory(uid int64)([]SearchData,[]SearchData){
	list:=make([]SearchData,0)
	hot:=make([]SearchData,0)
	if uid!=0{
		result,err:=db.GetCourse().Query("select id,content,number,today_number from qz_search_history where uid=? order by modify_time desc",uid)
		if err==nil{
			var id,number,todayNumber sql.NullInt64
			var content sql.NullString
			for result.Next(){
				err=result.Scan(&id,&content,&number,&todayNumber)
				if err==nil{
					list=append(list,SearchData{
						Id:id.Int64,
						Number:number.Int64,
						TodayNumber:todayNumber.Int64,
						Content:content.String,
					})
				}
			}
		}
	}
	result,err:=db.GetCourse().Query("select content,sum(number) as total,sum(today_number) as total from qz_search_history group by content ORDER BY total desc")
	if err==nil{
		var number,todayNumber sql.NullInt64
		var content sql.NullString
		for result.Next(){
			err=result.Scan(&content,&number,&todayNumber)
			if err==nil{
				hot=append(hot,SearchData{
					Number:number.Int64,
					TodayNumber:todayNumber.Int64,
					Content:content.String,
				})
			}
		}
	}
	return list,hot
}

func GetSearchHint(content string)[]SearchData{
	list:=make([]SearchData,0)
	result,err:=db.GetCourse().Query("select content from qz_search_history where content like (?) group by content","%"+content+"%")
	if err==nil{
		var content sql.NullString
		for result.Next(){
			err=result.Scan(&content)
			if err==nil{
				list=append(list,SearchData{
					Content:content.String,
				})
			}
		}
	}
	return list
}
