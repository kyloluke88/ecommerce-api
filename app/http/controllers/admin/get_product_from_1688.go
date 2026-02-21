package admin

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"api/app/models/company"
	"api/app/models/product"
	productattr "api/app/models/product_attr"
	productimages "api/app/models/product_images"
	productsku "api/app/models/product_sku"
	productskupackage "api/app/models/product_sku_package"
	productvideo "api/app/models/product_video"
	adminRequest "api/app/requests/admin"
	"api/pkg/response"
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

// SigninController 注册控制器
type ProductController struct {
	BaseAPIController
}

type CrawlResponse struct {
	Success bool   `json:"success"`
	HTML    string `json:"html"`
}

// 例：2件起批。
type PromotionModel struct {
	MixNumber            int    `json:"mixNumber"`
	PromotionDisplayName string `json:"promotionDisplayName"`
	PromotionLabel       string `json:"promotionLabel"`
}

type SKU struct {
	SpecAttrs     string `json:"specAttrs"`
	ProductId     uint64
	CanBookCount  int64           `json:"canBookCount"`
	SkuID         int64           `json:"skuId"`
	SaleCount     int64           `json:"saleCount"`
	PromotionSku  bool            `json:"promotionSku"`
	Price         decimal.Decimal `json:"price"`
	DiscountPrice decimal.Decimal `json:"discountPrice"`
}

type PieceWeightScale struct {
	PieceWeightScaleItems []PieceWeightScaleItem `json:"pieceWeightScaleInfo"`
	ColumnList            []ColumnInfo           `json:"columnList"`
}

type PieceWeightScaleItem struct {
	Volume float64 `json:"volume"`
	Title  string  `json:"sku1"`
	Length float64 `json:"length"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Weight int     `json:"weight"`
	SkuID  int64   `json:"skuId"`
}
type ColumnInfo struct {
	FID       int    `json:"fid"`
	Precision int    `json:"precision"`
	Name      string `json:"name"`
	Label     string `json:"label"`
}

type Video struct {
	CoverUrl string `json:"coverUrl"`
	VideoId  int64  `json:"videoId"`
	State    int    `json:"state"`
	Title    string `json:"title"`
	VideoUrl string `json:"videoUrl"`
}

type FeatureAttribute struct {
	Fid      int64  `json:"fid"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	IsActive bool
}

func (fa *FeatureAttribute) UnmarshalJSON(data []byte) error {
	// 默认值
	fa.IsActive = true
	// 用一个 alias，避免递归调用
	type Alias FeatureAttribute
	aux := (*Alias)(fa)
	// 正常反序列化
	return json.Unmarshal(data, aux)
}

func dump(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func (pc *ProductController) GetProductFrom1688(c *gin.Context) {

	req := adminRequest.GetProductFrom1688{}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.Error(c, err, "ShouldBindJSON ERR")
		return
	}

	crawlRes, err := getJsonFrom1688(c, req)

	htmlStr := crawlRes.HTML

	// html, err := os.ReadFile("1688_en.html")
	// if err != nil {
	// 	panic(err)
	// }
	// htmlStr := string(html)

	re := regexp.MustCompile(`window\.contextPath,\{"result"[\s\S]*?</script>`)

	jsonData := re.FindString(htmlStr)
	if jsonData == "" {
		panic("product json data not found")
	}

	// fmt.Printf("RESULT raw (len=%d) prefix: %q\n", len(jsonData), jsonData[:min(80, len(jsonData))])

	// 1 sku
	skuStr, err := extractJSONArray(jsonData, `"skuMap":`)

	if err != nil {
		panic(err)
	}
	var skus []SKU
	err = json.Unmarshal([]byte(skuStr), &skus)

	if err != nil {
		panic(err)
	}

	// 2 捆包
	packageStr, err := extractPieceWeightScale(jsonData)
	// fmt.Printf("RESULT raw (len=%d) prefix: %q\n", len(packageStr), packageStr[:min(80, len(packageStr))])

	if err != nil {
		panic(err)
	}

	var pws PieceWeightScale
	err = json.Unmarshal([]byte(packageStr), &pws)

	if err != nil {
		panic(err)
	}

	// dump(pws)

	// return

	// 3 主图
	imageListJSON, err := extractJSONArray(jsonData, `"imageList":`)
	if err != nil {
		panic(err)
	}

	var imageModels []productimages.ProductImage
	if err := json.Unmarshal([]byte(imageListJSON), &imageModels); err != nil {
		panic(err)
	}

	// dump(images)

	// 4 商品视频

	videoJSON, err := extractJSONObject(jsonData, `"video":`)
	if err != nil {
		panic(err)
	}

	var videoModel productvideo.ProductVideo
	if err := json.Unmarshal([]byte(videoJSON), &videoModel); err != nil {
		panic(err)
	}

	// dump(video)

	// 5 商品属性
	faJSON, err := extractJSONArray(jsonData, "featureAttributes")
	if err != nil {
		panic(err)
	}

	var attrs []FeatureAttribute
	if err := json.Unmarshal([]byte(faJSON), &attrs); err != nil {
		panic(err)
	}

	dump(attrs)

	// 6. 商品详情的url

	detailUrl, err := extract[string](jsonData, `"detailUrl":`, ValueString)

	if err != nil {
		panic(err)
	}

	dump(detailUrl)

	// 7 标题，货物所在地,最低价,最高价,价格字符串

	// re = regexp.MustCompile(`"location"\s*:\s*"[^"]*"`)
	// location := re.FindString(jsonData)

	location, err := extract[string](jsonData, `"location":`, ValueString)
	offerMaxPrice, err := extract[string](jsonData, `"offerMaxPrice":`, ValueString)
	offerMinPrice, err := extract[string](jsonData, `"offerMinPrice":`, ValueString)
	subject, err := extract[string](jsonData, `"subject":`, ValueString)
	offerPriceDisplay, err := extract[string](jsonData, `"offerPriceDisplay":`, ValueString)
	fmt.Printf("%q", offerPriceDisplay)

	// 创建公司
	companyName, err := extract[string](jsonData, `"companyName":`, ValueString)
	companyUrl, err := extract[string](jsonData, `"offerlistUrl":`, ValueString)
	companyId1688, err := extract[string](jsonData, `"sellerMemberId":`, ValueString)

	companyModel, err := company.FindOrCreateByCompanyID1688(companyId1688, companyName, companyUrl)
	if err != nil {
		response.Error(c, err, err.Error())
		return
	}

	// 创建商品
	maxPrice, _ := decimal.NewFromString(offerMaxPrice)
	minPrice, _ := decimal.NewFromString(offerMinPrice)

	skuPackageMap := make(map[int64]productskupackage.ProductSkuPackage, len(pws.PieceWeightScaleItems))

	for i := range pws.PieceWeightScaleItems {
		item := pws.PieceWeightScaleItems[i]
		skuPackageMap[item.SkuID] = item.ToModel()
	}

	// 创建sku
	productSkuModels := make([]productsku.ProductSku, 0, len(skus))
	for i := range skus {
		sku := skus[i].ToModel()

		if pkg, ok := skuPackageMap[skus[i].SkuID]; ok {
			sku.SkuPackage = &pkg // 注意：一般是指针
		}
		productSkuModels = append(productSkuModels, sku)
	}

	// product attrs
	productAttrsModels := make([]productattr.ProductAttr, 0, len(attrs))

	for i := range attrs {
		productAttrsModels = append(productAttrsModels, attrs[i].ToModel())
	}

	productId1688, err := extract[int64](jsonData, `"offerId":`, ValueInt)
	if err != nil {
		panic(err)
	}

	// product_id_1688 已存在则跳过，不做任何插入
	exists, existedProduct, err := product.ExistsByProductID1688(productId1688)
	if err != nil {
		response.Error(c, err, err.Error())
		return
	}
	if exists {
		response.Data(c, gin.H{
			"data":    "product already exists, skipped",
			"product": existedProduct.ID,
		})
		return
	}

	productModel := product.Product{
		CompanyID:     companyModel.ID,
		ProductId1688: productId1688,
		SourceUrl:     req.ProductUrl,
		Title:         subject,
		MaxPrice:      maxPrice,
		MinPrice:      minPrice,
		Location:      location,
		OnSale:        true,
		SaleCount:     0,
		Video:         &videoModel, // 可有可无
		Images:        imageModels,
		Skus:          productSkuModels,
		Attrs:         productAttrsModels,
	}

	if err := productModel.CreateWithRelations(); err != nil {
		response.Error(c, err, err.Error())
		return
	}

	createdProduct, err := product.GetByIDWithRelations(productModel.ID)
	if err != nil {
		response.Error(c, err, err.Error())
		return
	}

	response.Data(c, gin.H{
		"product": createdProduct,
	})

}
func extractPieceWeightScale(html string) (string, error) {
	key := `"pieceWeightScale":`
	idx := strings.Index(html, key)
	if idx == -1 {
		return "", fmt.Errorf("pieceWeightScale not found")
	}

	// 找到第一个 {
	start := strings.Index(html[idx+len(key):], "{")
	if start == -1 {
		return "", fmt.Errorf("opening { not found")
	}
	start = idx + len(key) + start

	depth := 0
	for i := start; i < len(html); i++ {
		switch html[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				// 只返回 {...}
				return html[start : i+1], nil
			}
		}
	}
	return "", fmt.Errorf("matching } not found")
}

// func extractStringValue(src, key string) (string, error) {

// 	start := strings.Index(src, key)
// 	if start == -1 {
// 		return "", fmt.Errorf("%s not found", key)
// 	}

// 	i := start + len(key)
// 	for i < len(src) && src[i] != '"' {
// 		i++
// 	}
// 	if i >= len(src) {
// 		return "", fmt.Errorf("invalid value")
// 	}

// 	j := i + 1
// 	for j < len(src) {
// 		if src[j] == '"' && src[j-1] != '\\' {
// 			raw := src[i : j+1]
// 			return strconv.Unquote(raw)
// 		}
// 		j++
// 	}

// 	return "", fmt.Errorf("unterminated string")
// }

func extract[T any](src, key string, vt ValueType) (T, error) {
	v, err := extractValue(src, key, vt)
	if err != nil {
		var zero T
		return zero, err
	}
	return v.(T), nil
}

type ValueType int

const (
	ValueString ValueType = iota
	ValueInt
)

func extractValue(src, key string, vt ValueType) (any, error) {
	start := strings.Index(src, key)
	if start == -1 {
		return nil, fmt.Errorf("%s not found", key)
	}

	i := start + len(key)

	// 跳过 : 和空格
	for i < len(src) && (src[i] == ':' || src[i] == ' ') {
		i++
	}

	if i >= len(src) {
		return nil, fmt.Errorf("invalid value")
	}

	switch vt {

	case ValueString:
		if src[i] != '"' {
			return nil, fmt.Errorf("expected string value")
		}

		j := i + 1
		for j < len(src) {
			if src[j] == '"' && src[j-1] != '\\' {
				raw := src[i : j+1]
				return strconv.Unquote(raw)
			}
			j++
		}
		return nil, fmt.Errorf("unterminated string")

	case ValueInt:
		j := i
		for j < len(src) && (src[j] == '-' || (src[j] >= '0' && src[j] <= '9')) {
			j++
		}

		if i == j {
			return nil, fmt.Errorf("expected int value")
		}

		return strconv.ParseInt(src[i:j], 10, 64)

	default:
		return nil, fmt.Errorf("unknown value type")
	}
}

func extractJSONObject(html, key string) (string, error) {
	idx := strings.Index(html, key)
	if idx == -1 {
		return "", fmt.Errorf("%s not found", key)
	}

	// 找到第一个 {
	start := strings.Index(html[idx+len(key):], "{")
	if start == -1 {
		return "", fmt.Errorf("opening { not found")
	}
	start = idx + len(key) + start

	depth := 0
	for i := start; i < len(html); i++ {
		switch html[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				// 只返回 {...}
				return html[start : i+1], nil
			}
		}
	}

	return "", fmt.Errorf("matching } not found")
}

func extractJSONArray(html, key string) (string, error) {

	idx := strings.Index(html, key)
	if idx == -1 {
		return "", fmt.Errorf("%s not found", key)
	}

	start := strings.Index(html[idx+len(key):], "[")
	if start == -1 {
		return "", fmt.Errorf("opening [ not found")
	}
	start = idx + len(key) + start

	depth := 0
	for i := start; i < len(html); i++ {
		switch html[i] {
		case '[':
			depth++
		case ']':
			depth--
			if depth == 0 {
				return html[start : i+1], nil
			}
		}
	}

	return "", fmt.Errorf("matching ] not found")
}

func (sku *SKU) ToModel() productsku.ProductSku {
	return productsku.ProductSku{
		SkuID1688:     sku.SkuID,
		Price:         sku.Price,
		PromotionSku:  sku.PromotionSku,
		Stock:         sku.CanBookCount,
		DiscountPrice: sku.DiscountPrice,
		SaleCount:     sku.SaleCount,
		Title:         sku.SpecAttrs,
	}
}

func (pwsi *PieceWeightScaleItem) ToModel() productskupackage.ProductSkuPackage {
	return productskupackage.ProductSkuPackage{
		Volume:    pwsi.Volume,
		Title:     pwsi.Title,
		Length:    pwsi.Length,
		Width:     pwsi.Width,
		Height:    pwsi.Height,
		Weight:    pwsi.Weight,
		SkuID1688: pwsi.SkuID,
	}
}

func (fa *FeatureAttribute) ToModel() productattr.ProductAttr {
	return productattr.ProductAttr{
		Fid:      fa.Fid,
		Name:     fa.Name,
		Value:    fa.Value,
		IsActive: fa.IsActive,
	}
}

func getJsonFrom1688(c *gin.Context, req adminRequest.GetProductFrom1688) (CrawlResponse, error) {
	reqBody, _ := json.Marshal(req)

	resp, err := http.Post(
		"http://crawler:4000/crawl",
		"application/json",
		bytes.NewBuffer(reqBody),
	)

	if err != nil {
		response.Error(c, err, "crawler service ERROR")
		return CrawlResponse{}, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result CrawlResponse
	if err := json.Unmarshal(body, &result); err != nil {
		response.Error(c, err, "crawler data unmarshal Error")
		return CrawlResponse{}, err
	}

	return result, nil
}
