package goodmodel

import (
	"database/sql"
	"strings"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/utils"

	"mufe_service/camp/xlog"
)

// SpuPhoto 商品图片
type SpuPhoto struct {
	ID     int64
	SpuID  int64
	Key    string
	Prefix string
}

// GetSpuPhotoBySpuIDs 获取商品图片
func GetSpuPhotoBySpuIDs(spuIDs []int64) (result []SpuPhoto, err error) {
	if len(spuIDs) == 0 {
		return result, nil
	}
	var buf strings.Builder
	buf.WriteString(`select p.id,p.spu_id,p.prefix,p.key from qz_good_spu_photo p`)
	utils.MysqlStringInUtils(&buf, spuIDs, " where p.spu_id ")
	buf.WriteString(` and p.status = ?  order by is_main desc,sort desc,id asc`)
	rows, err := db.GetGoodDb().Query(buf.String(), enum.StatusNormal)
	if err != nil {
		return result, xlog.Error(err)
	}
	for rows.Next() {
		var id, spuID sql.NullInt64
		var prefix, key sql.NullString
		err := rows.Scan(&id, &spuID, &prefix, &key)
		if err != nil {
			return result, xlog.Error(err)
		}
		result = append(result, SpuPhoto{
			ID:     id.Int64,
			SpuID:  spuID.Int64,
			Key:    key.String,
			Prefix: prefix.String,
		})
	}
	if err := rows.Err(); err != nil {
		return result, xlog.Error(err)
	}

	return result, nil
}

//GetSpuPhotoByCategoryIDs 获取指定商品类别下的某一商品的图片
func GetSpuPhotoByCategoryIDs(ids []int64) (result map[int64]SpuPhoto, err error) {
	result = make(map[int64]SpuPhoto, 0)
	if len(ids) == 0 {
		return
	}
	ids = utils.RemoveDupsInt64(ids)
	args := []interface{}{}
	ph := []string{}
	for _, v := range ids {
		args = append(args, v)
		ph = append(ph, "?")
	}

	rows, err := db.GetGoodDb().Query(`select c.category_id,p.key,p.prefix from qz_good_category c 
left join qz_good_category_record r on r.category_id = c.category_id 
left join qz_good_spu_photo p on r.spu_id = p.spu_id 
where c.category_id in (`+strings.Join(ph, ",")+`) group by c.category_id`, args...)
	if err != nil {
		return result, xlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var categoryID sql.NullInt64
		var key, prefix sql.NullString
		err := rows.Scan(&categoryID, &key, &prefix)
		if err != nil {
			return result, xlog.Error(err)
		}
		result[categoryID.Int64] = SpuPhoto{
			Key:    key.String,
			Prefix: prefix.String,
		}
	}
	err = rows.Err()
	if err != nil {
		return result, xlog.Error(err)
	}
	return result, nil
}
