package client

import (
	"api/app/models/product"
	"api/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	BaseAPIController
}

type homeProductItem struct {
	ID         uint64  `json:"id"`
	Title      string  `json:"title"`
	Rating     float64 `json:"rating"`
	MinPrice   string  `json:"min_price"`
	MaxPrice   string  `json:"max_price"`
	PriceRange string  `json:"price_range"`
	ImageURL   string  `json:"image_url"`
}

func (pc *ProductController) HomeProducts(c *gin.Context) {
	newLimit := parseLimitQuery(c.Query("new_limit"), 12, 60)
	recommendedLimit := parseLimitQuery(c.Query("recommended_limit"), 15, 60)

	newProducts, recommendedProducts, err := product.ListHomeSections(newLimit, recommendedLimit)
	if err != nil {
		response.Error(c, err, "query home products failed")
		return
	}

	response.Data(c, gin.H{
		"new_arrivals":         mapHomeProducts(newProducts),
		"recommended_products": mapHomeProducts(recommendedProducts),
	})
}

func mapHomeProducts(products []product.Product) []homeProductItem {
	items := make([]homeProductItem, 0, len(products))
	for i := range products {
		imageURL := ""
		if len(products[i].Images) > 0 {
			imageURL = products[i].Images[0].FullPathImageURI
		}

		minPrice := products[i].MinPrice.String()
		maxPrice := products[i].MaxPrice.String()

		items = append(items, homeProductItem{
			ID:         products[i].ID,
			Title:      products[i].Title,
			Rating:     products[i].Rating,
			MinPrice:   minPrice,
			MaxPrice:   maxPrice,
			PriceRange: minPrice + " ~ " + maxPrice,
			ImageURL:   imageURL,
		})
	}

	return items
}

func parseLimitQuery(raw string, defaultValue int, maxValue int) int {
	if raw == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return defaultValue
	}

	if value > maxValue {
		return maxValue
	}

	return value
}
