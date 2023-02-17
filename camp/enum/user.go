package enum

const (
	OuthType      int64 = 1
	PhoneType     int64 = 2
	PhonePassType int64 = 3

	AppleType int64 = 4
	JpushType int64 = 5

	//
	LoginTypeWxApp     int64 = 100
	LoginTypeWxMini    int64 = 101
	LoginTypeQQ        int64 = 11
	LoginTypeWeiBo     int64 = 12
	LoginTypeApp       int64 = 13
	LoginTypeJpush     int64 = 14
	LoginTypePhone     int64 = 2
	LoginTypePhonePass int64 = 3
)

const (
	UpdateUserName       = 1
	UpdateUserHead       = 2
	UpdateUserPhone      = 3
	UpdateUserPass       = 4
	UpdateUserIdentity   = 5
	UpdateUserModifyPass = 6
	UpdateUserSign       = 7
	UpdateUserSex        = 8
	UpdateUserAddress    = 9
	UpdateUserInfo       = 10
	UpdatePushInfo       = 11
)

const (
	CodePhoneLoginType = 1
	CodeBindPhoneType  = 2
	CodeRegisterType   = 3
	CodeForgetType     = 4
	CodeCancelType     = 5
)

// 第三方登录类型
const (
	OuthTypeWx int64 = 1 // 微信
)

const (
	GetProvince int64 = 1 //获取省份
	GetCity     int64 = 2 //获取城市
	GetDistrict int64 = 3 //获取地区
)

const (
	TEACHER_TYPE     int64 = 2
	DEFAULT_TYPE     int64 = 1
	STUDENT_TYPE     int64 = 3
	INSTITUTION_TYPE int64 = 4
)

const (
	CANCEL_DEFUALT_TYPE int64 = 0
	CANCEL_START_TYPE   int64 = 1
	CANCEL_END_TYPE     int64 = 2
)
