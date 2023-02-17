package enum

const(
	SchoolCodePrefix string ="qzSchoolCode?"
	ClassCodePrefix string ="qzClassCode?"
	InstitutionCodePrefix string ="qzInstitutionCode?"
	FootBallType string="足球"
)

const(
	ClassAdminType int64 =1
	ClassStudentType int64 =2
)

const (
	SCAN_SCHOOL_TYPE int64 =1
	SCAN_CLASS_TYPE int64 =2
	SCAN_INSTITUTION_TYPE int64 =3
)

// 商家身份
const (
	SYSTEM_ADMIN_GROUP   int64 = 1 // 系统商
	BRAND_ADMIN_GROUP    int64 = 2 // 品牌商
	OPERATOR_ADMIN_GROUP int64 = 3 // 运营商
	AGENT_ADMIN_GROUP    int64 = 4 // 代理商
)
