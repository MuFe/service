package enum

const SQLPlaceholderLimit = 10000
// 时间格式
const (
	TimeFormatMonthTimeNoSymbol = "0102150405"
	TimeFormatTimeNoSymbol      = "20060102150405"
	TimeFormatMonth             = "2006-01"
	TimeFormatDate              = "2006-01-02"
	TimeFormatTime              = "2006-01-02 15:04:05"
	TimeFormatTimeCh            = "2006年01月02日 15:04:05"
	TimeFormatMinuteCh          = "2006年01月02日 15:04"
	TimeFormatMonthDay          = "0102"
	TimeFormatHour          = "15"
)

// 状态
const (
	StatusAll         int64 = -1 // 全部
	StatusVerify      int64 = 0  // 待审核
	StatusNormal      int64 = 1  // 正常
	StatusVerifyFail  int64 = 2  // 审核失败
	StatusClose       int64 = 3  // 关闭
	StatusDelete      int64 = 4  // 删除
	StatusExpired     int64 = 5  // 过期
	StatusDraft       int64 = 5  // 草稿
	StatusExceptional int64 = 6  // 异常
)

// 上下架
const (
	StatusAllSale int64 = 0 //全部
	StatusSale    int64 = 1 //已上架
	StatusUnSale  int64 = 2 //已下架
)

const (
	AddTag  int64 = 1
	EditTagCover int64 = 2
)


const (
	CourseTagType int64=1
	ChapterTagType int64=2
	VideoTagType int64=3
	HomeWorkTagType int64=4
)


const (
	SortMax      int64 = 1
	SortUp       int64 = 2
	SortDown     int64 = 3
	SortMaxValue int64 = 99
)

const(
	FreeCanSee int64 =1
	FreeSchoolNoSee int64 =2
	PriceNoSee int64 =3
)

const (
	WebIcp int64=1
	WebAndroid int64=2
	WebIos int64=3
	WebShareTitle int64=4
	WebShareContent int64=5
	WebShareIcon int64=6
	WebShareUrl int64=7
)

const(
	HomeWork_Edit_Cover int64=1
	HomeWork_Edit_Content int64=2
	HomeWork_Edit_Video int64=3
)

const(
	EDIT_COVER   int64 =1
	EDIT_VIDEO   int64 =2
	EDIT_PLAN    int64 =4
	EDIT_CONTENT int64 =3
)

const(
	RECOMMEND_LEVEL int64=1
	RECOMMEND_SOURCE int64=2
	RECOMMEND_COURSE int64=3
	RECOMMEND_ITEM int64=4
)

const(
	RECOMMEND_MAIN int64 =1
)
