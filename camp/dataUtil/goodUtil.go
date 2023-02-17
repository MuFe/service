package dataUtil

import (
	pb "mufe_service/jsonRpc"
	"mufe_service/model/good"
)

type GoodDetail struct {
	Info Good `json:"info"`
}
type Good struct {
	SpuId        int64           `json:"spu_id"`
	SpuNumber    string          `json:"spu_number"`
	SpuName      string          `json:"spu_name"`
	End          int64           `json:"end"`
	Stock        int64           `json:"stock"`
	DeliveryInfo SkuDeliveryInfo `json:"delivery"`
	Status       int64           `json:"status"`
	Location     string          `json:"location"`
	Detail       string          `json:"detail"`
	CategoryId   int64           `json:"category_id"`
	SaleTime     int64           `json:"sale_time"`
	CommentNum   int64           `json:"comment_num"`
	List         []ShopSkuInfo   `json:"list"`
	Photos       []PhotoInfo     `json:"photos"`
	IsLike       bool            `json:"is_like"`
	IsMember     bool            `json:"is_member"`
	ShopName     string          `json:"shop_name"`
	ShopId       int64           `json:"shop_id"`
	ShopPhoto    string          `json:"shop_photo"`
	ShopScore    float64         `json:"shop_score"`
}

type SkuDeliveryInfo struct {
	DefaultPrice int64 `json:"price"`
	DefaultNum   int64 `json:"num"`
	AddPrice     int64 `json:"add_price"`
	AddNum       int64 `json:"add_num"`
}

type ShopSkuInfo struct {
	SkuId     int64       `json:"sku_id"`
	SkuNumber string      `json:"sku_number"`
	SkuName   string      `json:"sku_name"`
	Price     int64       `json:"price"`
	Stock     int64       `json:"stock"`
	Photo     string      `json:"photo"`
	Options   []SkuOption `json:"options"`
}

type SkuOption struct {
	OptionId      int64  `json:"id"`
	OptionValue   string `json:"value"`
	OptionValueId int64  `json:"o_id"`
}

type PhotoInfo struct {
	Key string `json:"key"`
	Url string `json:"url"`
}

type AdminParentGoodCategoryInfo struct {
	Id       int64                   `json:"id"`
	Name     string                  `json:"name"`
	Category []AdminGoodCategoryInfo `json:"category"`
}

type AdminGoodCategoryInfo struct {
	Id    int64                    `json:"id"`
	Name  string                   `json:"name"`
	PName string                   `json:"p_name"`
	PId   int64                    `json:"p_id"`
	Spec  []GoodSpecificationsInfo `json:"spec"`
}

type GoodSpecificationsInfo struct {
	Id     int64                     `json:"id"`
	Name   string                    `json:"name"`
	Values []GoodSpecificationsValue `json:"values"`
}

type GoodSpecificationsValue struct {
	Id      int64  `json:"id"`
	Content string `json:"content"`
}

type DeliveryInfo struct {
	Type       int64 `json:"type"`
	TemplateId int64 `json:"template_id"`
}



func ParseCategoryInfoList(list []*pb.GoodCategory) []AdminParentGoodCategoryInfo {
	resultList := make([]AdminParentGoodCategoryInfo, 0)
	for _, info := range list {
		result := AdminParentGoodCategoryInfo{}
		result.Id = info.Id
		result.Name = info.Name
		for _, category := range info.Children {
			specs := make([]GoodSpecificationsInfo, 0)
			for _, value := range category.Specs {
				specs = append(specs, parseSpec(value))
			}
			result.Category = append(result.Category, AdminGoodCategoryInfo{
				Id:    category.Id,
				Name:  category.Name,
				PId:   info.Id,
				PName: info.Name,
				Spec:  specs,
			})
		}
		resultList = append(resultList, result)
	}
	return resultList
}

func parseSpec(info *pb.GoodSpecificationsInfo) GoodSpecificationsInfo {
	result := GoodSpecificationsInfo{}
	result.Name = info.Name
	result.Id = info.Id
	list := make([]GoodSpecificationsValue, 0, len(info.List))
	for k := range info.List {
		list = append(list, parseSpecValue(info.List[k]))
	}
	result.Values = list
	return result
}

func parseSpecValue(info *pb.GoodSpecificationsValue) GoodSpecificationsValue {
	result := GoodSpecificationsValue{}
	result.Content = info.Content
	result.Id = info.Id
	return result
}

func MakeCategory(categoryInfo []*goodmodel.GoodCategoryInfo) (result []*pb.GoodCategory) {
	result = make([]*pb.GoodCategory, 0)
	if len(categoryInfo) == 0 {
		return
	}

	for _, c := range categoryInfo {
		cinfo := &pb.GoodCategory{
			Id:           c.ID,
			Name:         c.Name,
			Level:        c.Level,
			WholesSpcNum: c.WholesCount,
			ShopSpcNum:   c.ShopCount,
			TmSpcNum:     c.LocalCount,
			Children:     MakeCategory(c.Children),
		}
		result = append(result, cinfo)
	}
	return
}

func ParseShopGoodDetail(info *pb.GoodDetailResponse) GoodDetail {
	result := GoodDetail{}
	result.Info = ParseListData(info)
	result.Info.List = make([]ShopSkuInfo, 0)
	for _, temp := range info.List {
		tempInfo := parseShopSkuInfos(temp)
		if len(result.Info.Photos) > 0 {
			tempInfo.Photo = result.Info.Photos[0].Url
		}

		result.Info.List = append(result.Info.List, tempInfo)
		result.Info.Stock += temp.Stock
	}
	return result
}

func ParseListData(info *pb.GoodDetailResponse) Good {
	data := Good{}
	data.SpuId = info.SpuId
	data.SpuNumber = info.SpuNumber
	data.SpuName = info.SpuName
	data.Location = info.Location
	data.Detail = info.Detail
	data.CategoryId = info.CategoryId
	data.Status = info.Status
	data.CommentNum = info.CommentNum
	data.Photos = make([]PhotoInfo, 0)
	for _, temp := range info.Photos {
		data.Photos = append(data.Photos, parsePhotos(temp))
	}
	return data
}

func parsePhotos(info *pb.PhotoInfo) PhotoInfo {
	result := PhotoInfo{}
	result.Key = info.Key
	result.Url = info.Url
	return result
}

func parseShopSkuInfos(info *pb.Sku) ShopSkuInfo {
	result := ShopSkuInfo{}
	result.SkuId = info.SkuId
	result.SkuName = info.SkuName
	result.SkuNumber = info.SkuNumber
	result.Price = info.Price
	result.Stock = info.Stock
	result.Options = make([]SkuOption, 0)
	for _, temp := range info.Options {
		result.Options = append(result.Options, parseSkuOption(temp))
	}

	return result
}

func parseSkuOption(info *pb.SkuOption) SkuOption {
	result := SkuOption{}
	result.OptionId = info.OptionId
	result.OptionValue = info.OptionValue
	result.OptionValueId = info.OptionValueId
	return result
}

func ParseSkuDelivery(infos []*pb.DeliveryData) SkuDeliveryInfo {
	if len(infos) == 0 {
		return SkuDeliveryInfo{
			DefaultNum:   0,
			DefaultPrice: 0,
			AddNum:       0,
			AddPrice:     0,
		}
	} else {
		return parseSkuDelivery(infos[0])
	}
}

func parseSkuDelivery(temp *pb.DeliveryData) SkuDeliveryInfo {
	info := SkuDeliveryInfo{}
	info.DefaultNum = temp.DefaultNum
	info.DefaultPrice = temp.DefaultPrice
	info.AddNum = temp.IncreaseNum
	info.AddPrice = temp.IncreasePrice
	return info
}
