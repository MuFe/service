package sequence

import (
	"database/sql"
	"errors"
	"fmt"
	"mufe_service/camp/db"
	"mufe_service/camp/utils"
	"strconv"
	"strings"

	"mufe_service/camp/xlog"
)

var (
	CategoryNumber = sequencePrefixFirst{prefix: "", seq: "category_number", length: 8}
	SpuNo             = sequencePrefixFirst{prefix: "S", seq: "spu_no", length: 8}
	SkuNo             = sequencePrefixFirst{prefix: "G", seq: "sku_no", length: 8}
	UserNo            = sequencePrefix{prefix: "U", seq: "user_no", length: 8}
	UserInviteCode    = sequenceUserCodePrefixFirst{prefix: "", seq: "user_invite_code", length: 8}
)

type sequencePrefix struct {
	seq    string
	prefix string
	length int64
}

func (t *sequencePrefix) NewNo() (string, error) {
	if t.prefix == "" {
		return "", xlog.Error("没有设置编号前缀")
	}
	no, err := newNo(t.seq, t.length)
	if err != nil {
		return "", xlog.Error(err)
	}
	return t.prefix + no, nil
}

type sequencePrefixFirst struct {
	seq    string
	prefix string
	length int64
}

func (t *sequencePrefixFirst) NewNo(name string) (string, error) {
	no, err := newNo(t.seq, t.length)
	if err != nil {
		return "", err
	}
	return t.prefix + utils.ConvertFirst(name) + no, nil
}

type sequenceMachinePrefix struct {
	seq    string
	prefix string
	length int64
}

func (t *sequenceMachinePrefix) NewNo(name string) (string, error) {
	if t.prefix == "" {
		return "", xlog.Error("没有设置编号前缀")
	}
	if strings.Index(name, t.prefix) != -1 {
		return name, nil
	}
	no, err := newNo(t.seq, t.length)
	if err != nil {
		return "", xlog.Error(err)
	}
	return t.prefix + no + utils.GetRandomInt(8), nil
}

type sequenceUserCodePrefixFirst struct {
	seq    string
	prefix string
	length int64
}

func (t *sequenceUserCodePrefixFirst) NewNo() (string, error) {
	no, err := newNo(t.seq, t.length)
	if err != nil {
		return "", err
	}
	noInt, _ := strconv.ParseInt(no, 10, 64)
	return utils.DecimalToAny(noInt, 62), nil
}

// 获取编号
func newNo(seqName string, length int64) (string, error) {
	no := ""
	n := 1
	if err := db.GetUserDb().WithTransaction(func(tx *db.Tx) error {
		for {
			r, err := tx.Exec(`update qz_seq set seq=seq+? where name = ?`, n, seqName)
			if err != nil {
				return err
			}
			affected, err := r.RowsAffected()
			if err != nil {
				return err
			}
			if affected == 0 {
				return errors.New("no seq name " + seqName)
			}

			var seq32 sql.NullInt32
			err = tx.QueryRow("select seq from qz_seq where name = ?", seqName).Scan(&seq32)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			seq := strconv.Itoa(int(seq32.Int32))
			index := strings.Index(seq, "4")
			if index == -1 {
				no = fmt.Sprintf("%0"+strconv.Itoa(int(length))+"s", seq)
				break
			} else {
				pow := n * (len(seq) - index - 1)
				for i := 0; i < pow; i++ {
					n = n * 10
				}
			}
		}
		return nil
	}); err != nil {
		return "", err
	}
	return no, nil
}
