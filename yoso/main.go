// app.go
package main

import (
	"fmt"
	"net/http"
	"shop-the-look/models"

	"github.com/gin-gonic/gin"
)

var products []models.Product      
var looks []models.Look
var nextLookID int = 1
var nextProductID int = 1

func main() {
	router := gin.Default()
	// Routes for product management
	router.POST("/product", AddProduct)        // POST: Add Product
	router.GET("/product", ListProducts)	   // GET: List all Products

	// Routes for look management
	router.POST("/look", CreateLook)           // POST: Create Look
	router.GET("/look", ListLooks)              // GET: List all Looks
	router.GET("/look/:look_id", GetLook)      // GET: Read Look
	router.PUT("/look/:look_id/product", AddProductToLook) // PUT: Add Product to Look
	router.DELETE("/look/:look_id/product/:product_id", RemoveProductFromLook) // DELETE: Remove Product from Look

	
	router.Run(":8080")
}


// API to add a product
func AddProduct(c *gin.Context) {
	var newProduct models.Product
	if err := c.BindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	newProduct.ProductID = nextProductID
	nextProductID++
	products = append(products, newProduct)

	c.JSON(http.StatusCreated, newProduct)
}


// API to list all products
func ListProducts(c *gin.Context) {
	c.JSON(http.StatusOK, products)
}


// API to create a Look using Product IDs
func CreateLook(c *gin.Context) {
	var request struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ProductIDs  []int  `json:"product_ids"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var selectedProducts []models.Product
	var totalPrice float64
	minAvailableQty := -1

	for _, productID := range request.ProductIDs {
		found := false
		for _, product := range products {
			if product.ProductID == productID {
				selectedProducts = append(selectedProducts, product)
				totalPrice += 0.8*product.PriceInINR

				// Track minimum availability across selected products
				if minAvailableQty == -1 || product.AvailableQty < minAvailableQty {
					minAvailableQty = product.AvailableQty
				}
				found = true
				break
			}
		}
		if !found {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Product with ID %d not found", productID)})
			return
		}
	}
	// Create a new Look
	newLook := models.Look{
		ID:           nextLookID,
		Name:         request.Name,
		Description:  request.Description,
		PriceInINR:   totalPrice,
		AvailableQty: minAvailableQty,
		Products:     selectedProducts,
		Discount:     0.2,
	}
	nextLookID++
	looks = append(looks, newLook)

	c.JSON(http.StatusCreated, newLook)
}

// API to list all products
func ListLooks(c *gin.Context) {
	c.JSON(http.StatusOK, looks)
}

// API to get a specific Look by ID
func GetLook(c *gin.Context) {
	lookID := c.Param("look_id")
	for _, look := range looks {
		if fmt.Sprintf("%d", look.ID) == lookID {
			c.JSON(http.StatusOK, look)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Look not found"})
}

// API to add a product to a Look
func AddProductToLook(c *gin.Context) {
	lookID := c.Param("look_id")
	var request struct {
		ProductID int `json:"product_id"`
	}

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Find Look
	var look *models.Look
	for i := range looks {
		if fmt.Sprintf("%d", looks[i].ID) == lookID {
			look = &looks[i]
			break
		}
	}
	if look == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Look not found"})
		return
	}

	// Find Product
	var product *models.Product
	for i := range products {
		if products[i].ProductID == request.ProductID {
			product = &products[i]
			break
		}
	}
	if product == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	// Add product to Look
	look.Products = append(look.Products, *product)
	look.PriceInINR += product.PriceInINR
	if product.AvailableQty < look.AvailableQty || look.AvailableQty == 0 {
		look.AvailableQty = product.AvailableQty
	}

	c.JSON(http.StatusOK, look)
}



// API to remove a product from a Look
func RemoveProductFromLook(c *gin.Context) {
	lookID := c.Param("look_id")
	productID := c.Param("product_id")

	// Find Look
	var look *models.Look
	for i := range looks {
		if fmt.Sprintf("%d", looks[i].ID) == lookID {
			look = &looks[i]
			break
		}
	}
	if look == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Look not found"})
		return
	}

	// Remove product by ID
	removed := false
	var updatedProducts []models.Product
	var totalPrice float64
	minAvailableQty := -1
	for _, product := range look.Products {
		if fmt.Sprintf("%d", product.ProductID) == productID {
			removed = true
			continue
		}
		updatedProducts = append(updatedProducts, product)
		totalPrice += product.PriceInINR
		if minAvailableQty == -1 || product.AvailableQty < minAvailableQty {
			minAvailableQty = product.AvailableQty
		}
	}

	if !removed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found in Look"})
		return
	}

	// Update Look details
	look.Products = updatedProducts
	look.PriceInINR = totalPrice
	look.AvailableQty = minAvailableQty

	c.JSON(http.StatusOK, look)
}