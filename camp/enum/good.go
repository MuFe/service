package enum

// 栏目
const (
	GetAllCategory        int64 = 1 //获取全部类目数组
	GetCategoryFromSpu    int64 = 2 //获取库存所属类目或者渠道类目
	GetCategoryFromQuery  int64 = 3 //获取指定类目
	GetCategoryInfoFromId int64 = 4 //获取分类以及父分类

	FirstLevel  int64 = 1
	SecondLevel int64 = 2
	ThirdLevel  int64 = 3
)

// 操作
const (
	GoodAdd            int64 = 1 //添加
	GoodEditFromList   int64 = 2 //为列表中修改
	GoodEditFromDetail int64 = 3 //为详情中修改
	GoodEditDetail     int64 = 4 //单独更新detail
	GoodEditPhoto      int64 = 5 //单独更新轮播图
	GoodEditSku        int64 = 6 //为列表中修改sku
	GoodEditSort       int64 = 7 //为列表中修改排序
)

// 库存类型
const (
	SkuStockTypeBrandAdd                     = 1  //品牌商入库
	SkuStockTypeBrandOrderLess               = 2  //品牌商商城出库
	SkuStockTypeBrandOrderCancel             = 3  //品牌商商城取消订单入库
	SkuStockTypeBrandOrderRefund             = 4  //品牌商商城退货入库
)


// 配送类型:1-商城销售
const (
	DeliveryTypeShop      = 1
)
