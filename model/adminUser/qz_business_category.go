package adminUserModel

import (
	"database/sql"
	"strconv"
	"strings"
	"mufe_service/camp/db"

	"mufe_service/camp/xlog"
)

// BusinessCategory 商家的商品类别（商品的一级类别）
type BusinessCategory struct {
	ID         int64
	BusinessID int64
	CategoryID int64
	Status     int64
}

// GetBusinessCategoryByBusinessID 获取商家注册时的商品一级类别
func GetBusinessCategoryByBusinessID(ids []int64) (result []BusinessCategory, err error) {
	var buf strings.Builder
	buf.WriteString(`select bc.id,bc.business_id,bc.category_id,bc.status from qz_business_category bc where bc.status=1`)
	for k, id := range ids {
		if k == 0 {
			buf.WriteString(" and bc.business_id in (")
		}

		buf.WriteString(strconv.FormatInt(id, 10))
		if k < len(ids)-1 {
			buf.WriteString(",")
		} else {
			buf.WriteString(") ")
		}
	}
	buf.WriteString(" order by bc.id asc")
	rows, err := db.GetAdminDb().Query(buf.String())
	if err != nil {
		return result, xlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, bid, cid, status sql.NullInt64
		err := rows.Scan(&id, &bid, &cid, &status)
		if err != nil {
			return result, xlog.Error(err)
		}
		result = append(result, BusinessCategory{
			ID:         id.Int64,
			BusinessID: bid.Int64,
			CategoryID: cid.Int64,
			Status:     status.Int64,
		})
	}
	err = rows.Err()
	if err != nil {
		return result, xlog.Error(err)
	}

	return result, nil
}
