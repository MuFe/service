package goodmodel

import (
	"database/sql"

	"strings"
	"mufe_service/camp/db"
	"mufe_service/camp/enum"
	"mufe_service/camp/sequence"
	"mufe_service/camp/utils"
	pb "mufe_service/jsonRpc"

	"mufe_service/camp/xlog"
)

// GoodCategory 商品类别
type GoodCategory struct {
	ID          int64
	Name        string
	ParentID    int64
	TopParentID int64
	Icon        string
	Photo       string
	Sort        int64
	Status      int64
	Level       int64
	ShopCount   int64
	CategoryNum int64 //三级分类个数
}

type GoodCategoryInfo struct {
	ID          int64
	ParendId    int64
	Level       int64
	Name        string
	ShopCount   int64
	LocalCount  int64
	WholesCount int64
	Spec        []*SpecValue
	Children    []*GoodCategoryInfo
}

type SpecInfo struct {
	Id         int64
	Name       string
	CategoryId int64
	Values     []SpecValue
}

type SpecValue struct {
	Id    int64
	Name  string
	Value string
}

// GetGoodCategorySpuCountByCategoryIds 获取指定类别的商品数
func GetGoodCategorySpuCountByCategoryIds(ids []int64) (map[int64]*GoodCategory, error) {
	var allBuf strings.Builder
	allBuf.WriteString(`select tb.category_id,tb.level,tb.category_name,tb1.category_id,tb2.category_id from qz_good_category tb 
left join qz_good_category tb1 on tb1.category_id=tb.parent_id
left join qz_good_category tb2 on tb2.category_id=tb1.parent_id`)
	goodCategoryResult, err := db.GetGoodDb().Query(allBuf.String())

	if err != nil {
		return nil, err
	}
	result := make(map[int64]*GoodCategory)
	secondResult := make(map[int64]*GoodCategory)
	topResult := make(map[int64]*GoodCategory)
	var categoryId, pCategoryId, tCategoryId, level sql.NullInt64
	var categoryName string
	for goodCategoryResult.Next() {
		_ = goodCategoryResult.Scan(&categoryId, &level, &categoryName, &pCategoryId, &tCategoryId)
		temp := &GoodCategory{
			ID:          categoryId.Int64,
			ParentID:    pCategoryId.Int64,
			TopParentID: tCategoryId.Int64,
			Name:        categoryName,
			Level:       level.Int64,
		}
		result[categoryId.Int64] = temp
	}
	for _, info := range result {
		if info.Level == 3 {
			temp, ok := result[info.ParentID]
			if ok {
				secondResult[info.ID] = temp
			}
			temp, ok = result[info.TopParentID]
			if ok {
				topResult[info.ID] = temp
			}
		}
	}
	var buf strings.Builder
	buf.WriteString(`select count(tb2.spu_id),tb.category_id,tb3.type from qz_good_category tb 
inner join qz_good_category_record tb1 on tb1.category_id=tb.category_id
inner join qz_good_spu tb2 on tb2.spu_id=tb1.spu_id
inner join qz_good_delivery tb3 on tb3.spu_id=tb1.spu_id
where tb2.status=? and tb2.sale_time>0`)
	buf.WriteString(" GROUP BY tb.category_id,tb3.type")
	rows, err := db.GetGoodDb().Query(buf.String(), enum.StatusNormal)
	if err != nil {
		return result, xlog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, count, typeInt sql.NullInt64
		err := rows.Scan(&count, &id, &typeInt)
		if err != nil {
			return nil, xlog.Error(err)
		}
		if v, ok := result[id.Int64]; ok {
			v.ShopCount = count.Int64
			if s, ok := secondResult[v.ID]; ok {
				s.ShopCount += count.Int64
				s.CategoryNum++
			}
			if s, ok := topResult[v.ID]; ok {
				s.ShopCount += count.Int64
				s.CategoryNum++
			}
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, xlog.Error(err)
	}
	return result, nil
}

func GetAllCategory() ([]*GoodCategoryInfo, error) {
	// 所有分类
	result, err := db.GetGoodDb().Query(`SELECT
	tb.category_id,
	tb.category_name,
	tb.level,
	tb.parent_id
FROM
	qz_good_category tb
WHERE
	tb.STATUS = 1 
ORDER BY
	tb.category_id ASC`)
	categoryIds := make([]int64, 0)
	category := make([]*GoodCategoryInfo, 0)
	if err != nil {
		return nil, err
	}
	defer result.Close()
	for result.Next() {
		var name sql.NullString
		var id, level, pId sql.NullInt64
		err = result.Scan(&id, &name, &level, &pId)
		if err != nil {
			return nil, xlog.Error(err)
		}
		if level.Int64 == enum.ThirdLevel {
			categoryIds = append(categoryIds, id.Int64)
		}
		category = append(category, &GoodCategoryInfo{
			ID:       id.Int64,
			ParendId: pId.Int64,
			Name:     name.String,
			Level:    level.Int64,
		})
	}
	countResult, err := GetGoodCategorySpuCountByCategoryIds([]int64{})
	if err == nil {
		for _, info := range category {
			temp, ok := countResult[info.ID]
			if ok {
				info.ShopCount = temp.ShopCount
			}
		}
	}

	err = result.Err()
	if err != nil {
		return nil, xlog.Error(err)
	}
	list := makeCategoryTree(&GoodCategoryInfo{}, category)
	return list, nil
}

func makeCategoryTree(parentCategory *GoodCategoryInfo, categoryList []*GoodCategoryInfo) []*GoodCategoryInfo {
	var result = make([]*GoodCategoryInfo, 0)
	for _, c := range categoryList {
		if c.ParendId == parentCategory.ID {
			c.Children = makeCategoryTree(c, categoryList)
			result = append(result, c)
		}
	}
	return result
}

//获取三级分类以及父分类信息
func GetCategoryFromId(categoryIdList, secondIdList, parentIdList []int64) (map[int64]*GoodCategoryInfo, error) {
	var buf strings.Builder
	buf.WriteString(`SELECT
	tb1.category_id,
	tb1.category_name,
	tb.category_id,
	tb.category_name,
	tb3.id,
	tb3.specification_name 
FROM
	qz_good_category tb
	left JOIN qz_good_category tb1 ON tb1.category_id = tb.parent_id
	left JOIN qz_good_category tb2 ON tb2.category_id = tb1.parent_id
	LEFT JOIN qz_good_specification tb3 ON tb3.category_id = tb.category_id
	WHERE
	tb.LEVEL = 3
	AND tb.STATUS = 1 
	AND tb1.STATUS = 1 
	AND tb2.STATUS = 1`)

	if len(utils.RemoveDupsInt64WithZero(categoryIdList)) > 0 && len(utils.RemoveDupsInt64WithZero(parentIdList)) > 0 {
		utils.MysqlStringInUtils(&buf, categoryIdList, " AND (tb.category_id ")
		utils.MysqlStringInUtils(&buf, parentIdList, " or tb2.category_id ")
		buf.WriteString(")")
	} else {
		if len(categoryIdList) > 0 {
			utils.MysqlStringInUtils(&buf, categoryIdList, " AND tb.category_id ")
		}
		if len(parentIdList) > 0 {
			utils.MysqlStringInUtils(&buf, parentIdList, " AND tb2.category_id ")
		}
		if len(secondIdList) > 0 {
			utils.MysqlStringInUtils(&buf, secondIdList, " AND tb1.category_id ")
		}
	}

	buf.WriteString(` ORDER BY
	tb.category_id ASC,
	tb1.category_id ASC,
	tb2.category_id ASC,
	tb3.id ASC`)
	result, err := db.GetGoodDb().Query(buf.String())
	if err == nil {
		var name, cName, csName sql.NullString
		var id, cSid, cId sql.NullInt64
		pMap := make(map[int64]*GoodCategoryInfo)
		cMap := make(map[int64]*GoodCategoryInfo)
		sMap := make(map[int64]*SpecValue)
		for result.Next() {
			err = result.Scan(&id, &name, &cId, &cName, &cSid, &csName)
			if err == nil {
				info, ok := pMap[id.Int64]
				if !ok {
					info = &GoodCategoryInfo{ID: id.Int64, Name: name.String, Children: make([]*GoodCategoryInfo, 0)}
					pMap[id.Int64] = info
				}
				cInfo, ok := cMap[cId.Int64]
				if !ok {
					cInfo = &GoodCategoryInfo{ID: cId.Int64, Name: cName.String, Spec: make([]*SpecValue, 0)}
					info.Children = append(info.Children, cInfo)
					cMap[cId.Int64] = cInfo
				}
				if cSid.Int64 != 0 {
					spec, ok := sMap[cSid.Int64]
					if !ok {
						spec = &SpecValue{Id: cSid.Int64, Name: csName.String}
						cInfo.Spec = append(cInfo.Spec, spec)
						sMap[cSid.Int64] = spec
					}
				}
			} else {
				xlog.ErrorP(err)
			}
		}
		return pMap, nil
	} else {
		return nil, err
	}
}

//获取spu相关的分类
func GetCategoryInfoFromSpu(channel, businessId, businessGroupId int64) ([]GoodCategoryInfo, error) {
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString(`SELECT
	tb.category_id,tb.category_name
FROM
	qz_good_category tb
	INNER JOIN qz_good_category_record tb1 ON tb1.category_id = tb.category_id
	INNER JOIN qz_good_spu tb2 ON tb2.spu_id = tb1.spu_id`)
	if businessId != 0 {
		//获取商家相关要加上spu关联
		buf.WriteString(`  INNER JOIN qz_good_sku tb3 on tb3.spu_id=tb2.spu_id
                              INNER JOIN qz_good_sku_stock tb5 on tb5.sku_id=tb3.sku_id and tb5.business_group_id=? and tb5.business_id=?`)
		args = append(args, businessGroupId, businessId)
	}
	buf.WriteString(" where tb2.status < 4 ")
	buf.WriteString(` GROUP BY tb.category_id
	order by tb.category_id asc`)
	result, err := db.GetGoodDb().Query(buf.String(), args...)
	list := make([]GoodCategoryInfo, 0)
	if err == nil {
		var name string
		var id int64
		for result.Next() {
			err = result.Scan(&id, &name)
			if err == nil {
				info := GoodCategoryInfo{ID: id, Name: name}
				list = append(list, info)
			} else {
				xlog.ErrorP(err)
			}
		}
		return list, nil
	} else {
		xlog.ErrorP(err)
		return nil, err
	}
}

//获取指定条件的分类类别
func GetCategoryInfoList(level int64, categoryIdList []int64) ([]GoodCategoryInfo, error) {
	var buf strings.Builder
	args := make([]interface{}, 0)
	buf.WriteString(`SELECT
	tb.category_id,tb.category_name
FROM
	qz_good_category tb`)
	buf.WriteString(" where 1=1 ")
	if len(categoryIdList) > 0 {
		utils.MysqlStringInUtils(&buf, categoryIdList, " AND tb.category_id ")
	}
	if level != 0 {
		//获取某个级别
		buf.WriteString(" AND tb.level=?")
		args = append(args, level)
	}
	buf.WriteString(` GROUP BY tb.category_id
	order by tb.category_id asc`)
	result, err := db.GetGoodDb().Query(buf.String(), args...)
	list := make([]GoodCategoryInfo, 0)
	if err == nil {
		var name string
		var id int64
		for result.Next() {
			err = result.Scan(&id, &name)
			if err == nil {
				info := GoodCategoryInfo{ID: id, Name: name}
				list = append(list, info)
			} else {
				xlog.ErrorP(err)
			}
		}
		return list, nil
	} else {
		xlog.ErrorP(err)
		return nil, err
	}
}

//修改分类
func EditGoodCategory(pid int64, name string, id int64, level int64, info []*pb.GoodSpecificationsInfo) error {
	if id == 0 && level == 3 && len(info) == 0 {
		return xlog.Error("新增3级分类时，规格不能为空")
	}
	err := db.GetGoodDb().WithTransaction(func(tx *db.Tx) error {
		// 删除规格
		if id != 0 {
			// 获取分类所有规格
			var updateSpecIdMap = make(map[int64]struct{}, 0)
			var deleteSpecIdMap = make(map[int64]struct{}, 0)
			var haveNew, deleteAll = false, true
			for _, v := range info {
				if v.Id != 0 {
					updateSpecIdMap[v.Id] = struct{}{}
				} else {
					haveNew = true
				}
			}
			rows, err := tx.Query(`select gs.id,gs.specification_name,gsp.spec_id from qz_good_specification gs 
			left join qz_good_specification_option gsp  on gsp.spec_id = gs.id
			where gs.category_id = ?`, id)
			if err != nil {
				return xlog.Error(err)
			}
			defer rows.Close()
			for rows.Next() {
				var id, usedId sql.NullInt64
				var name sql.NullString
				err := rows.Scan(&id, &name, &usedId)
				if err != nil {
					return xlog.Error(err)
				}
				_, ok := updateSpecIdMap[id.Int64]
				if !ok {
					if usedId.Int64 != 0 {
						return xlog.Errorf("规格“%s”已被使用，不能删除", name.String)
					} else {
						deleteSpecIdMap[id.Int64] = struct{}{}
					}
				} else {
					deleteAll = false
				}
			}
			err = rows.Err()
			if err != nil {
				return xlog.Error(err)
			}
			// 删除规格
			if !haveNew && deleteAll {
				return xlog.Error("必须保留一个规格")
			}
			args := []interface{}{}
			ph := []string{}
			for k := range deleteSpecIdMap {
				args = append(args, k)
				ph = append(ph, "?")
			}
			if len(args) > 0 {
				_, err = tx.Exec("delete from qz_good_specification where id in ("+strings.Join(ph, ",")+")", args...)
				if err != nil {
					return xlog.Error(err)
				}
			}
		}

		var err error
		if id == 0 {
			category_number, err := sequence.CategoryNumber.NewNo(name)
			if err != nil {
				return xlog.Error(err)
			}
			result, err := tx.Exec("insert into qz_good_category (category_name,category_number,parent_id,level) values(?,?,?,?)", name, category_number, pid, level)
			if err != nil {
				return xlog.Error(err)
			}
			idInt, _ := result.LastInsertId()
			id = int64(idInt)
		} else {
			_, err = tx.Exec("update  qz_good_category set category_name=?,parent_id=? where category_id=?", name, pid, id)
			if err != nil {
				return xlog.Error(err)
			}
		}

		for _, temp := range info {
			if temp.Id == 0 {
				_, err = tx.Exec("insert into qz_good_specification (specification_name,category_id,status) values(?,?,1)", temp.Name, id)
			} else {
				_, err = tx.Exec("update qz_good_specification set specification_name=? where id=?", temp.Name, temp.Id)
			}
			if err != nil {
				return xlog.Error(err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

//修改分类状态
func EditGoodCategoryStatus(idsList []int64, status int64) error {
	if len(idsList) == 0 {
		return nil
	}
	var buf strings.Builder
	buf.WriteString(`update qz_good_category set status = ?`)
	utils.MysqlStringInUtils(&buf, idsList, " where category_id ")
	_, err := db.GetGoodDb().Exec(buf.String(), status)
	if err != nil {
		return xlog.Error(err)
	}
	return nil
}

func GetGoodCategoryName(id int64) (string,error){
	categoryName := ""
	err := db.GetGoodDb().QueryRow("select category_name from qz_good_category where category_id=?", id).Scan(&categoryName)
	if err != nil {
		return categoryName, xlog.Error(err)
	}
	return categoryName,nil
}
