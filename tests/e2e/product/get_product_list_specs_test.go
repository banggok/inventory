package product_e2e_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"inventory_management/api/handler"
	"inventory_management/internal/repository"
	"inventory_management/internal/usecase"
	"inventory_management/pkg/db"
	"net/http"
	"net/http/httptest"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("GetProductList E2E Tests", func() {
	var productHandler *handler.ProductHandler
	var database *gorm.DB
	var sqlDB *sql.DB

	ginkgo.BeforeEach(func() {
		// Initialize test environment
		database, sqlDB = db.InitDB(true)
		TruncateTables(database) // Clean up before each test

		productRepo := repository.NewPostgresProductRepository(database)
		productUsecase := usecase.NewProductUsecase(productRepo)
		productHandler = handler.NewProductHandler(productUsecase)

		// Create some test products
		createTestProducts(productHandler, 5) // Helper function to create test products
	})

	ginkgo.AfterEach(func() {
		TruncateTables(database) // Clean up after each test
		sqlDB.Close()
	})
	ginkgo.Context("GET /products", func() {

		// 1. Successful retrieval of product list with valid query parameters
		ginkgo.It("should return a list of products successfully with valid query parameters", func() {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request with valid query parameters (sortBy, sortDirection, limit, offset)
			c.Request = httptest.NewRequest("GET", "/api/v1/products?sortBy=name&sortDirection=asc&limit=10&offset=0", nil)

			// Call the GetProductList handler
			productHandler.GetProductList(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusOK))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			// Convert "total" from float64 to int before comparing
			total := int(response["total"].(float64))
			gomega.Expect(total).To(gomega.Equal(5)) // Now comparing int to int

			// Check that the "products" array has 5 items
			gomega.Expect(response["products"]).To(gomega.HaveLen(5)) // 5 products created in setup
		})

		// 2. Scenario for an empty product list
		ginkgo.It("should return an empty list of products with a 200 status if no products exist", func() {
			// Clean up the database to ensure no products exist
			TruncateTables(database)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up request with valid query parameters (sortBy, sortDirection, limit, offset)
			c.Request = httptest.NewRequest("GET", "/api/v1/products?sortBy=name&sortDirection=asc&limit=10&offset=0", nil)

			// Call the GetProductList handler
			productHandler.GetProductList(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusOK))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			// Check that the "products" list is empty
			gomega.Expect(response["products"]).To(gomega.HaveLen(0))

			// Check that "total" is 0
			total := int(response["total"].(float64))
			gomega.Expect(total).To(gomega.Equal(0))
		})

		// 3. Internal server error
		ginkgo.It("should return 500 if there is an internal server error", func() {
			// Mock use case to return an internal server error
			mockUsecase := new(MockProductUsecase)
			productHandler := handler.NewProductHandler(mockUsecase)

			// Simulate an error when calling ListProducts after validation passes
			mockUsecase.On("ListProducts", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
				Return(nil, errors.New("internal server error"))

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Use valid query parameters to bypass validation
			c.Request = httptest.NewRequest("GET", "/api/v1/products?sortBy=name&sortDirection=asc&limit=10&offset=0", nil)

			// Call the GetProductList handler
			productHandler.GetProductList(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusInternalServerError)) // Expecting 500 Internal Server Error
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.Equal("failed to retrieve product"))
		})

	})

})

func createTestProducts(handler *handler.ProductHandler, count int) {
	for i := 1; i <= count; i++ {
		reqBody := map[string]string{"name": "Product " + strconv.Itoa(i)}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.CreateProduct(c)
	}
}
