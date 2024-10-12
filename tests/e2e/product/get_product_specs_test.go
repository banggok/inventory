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
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("GetProduct E2E Tests", func() {
	var productHandler *handler.ProductHandler
	var createdProductID uint
	var database *gorm.DB
	var sqlDB *sql.DB

	ginkgo.BeforeEach(func() {
		// Initialize test environment
		database, sqlDB = db.InitDB(true)
		TruncateTables(database) // Clean up before each test

		productRepo := repository.NewPostgresProductRepository(database)
		productUsecase := usecase.NewProductUsecase(productRepo)
		productHandler = handler.NewProductHandler(productUsecase)

		// Create a test product
		reqBody := map[string]string{"name": "Test Product"}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		productHandler.CreateProduct(c)
		var response map[string]interface{}
		_ = json.NewDecoder(w.Body).Decode(&response)
		createdProductID = uint(response["id"].(float64)) // Capture created product ID for future tests
	})

	ginkgo.AfterEach(func() {
		TruncateTables(database) // Clean up after each test
		sqlDB.Close()
	})

	ginkgo.Context("GET /products/:id", func() {
		ginkgo.It("should retrieve an existing product successfully", func() {
			// Create a request to get the created product by ID
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(createdProductID))}}
			c.Request = httptest.NewRequest("GET", "/api/v1/products/"+strconv.Itoa(int(createdProductID)), nil)

			// Call the GetProduct handler
			productHandler.GetProduct(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusOK))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["name"]).To(gomega.Equal("Test Product"))
			gomega.Expect(response["id"]).To(gomega.Equal(float64(createdProductID)))
		})

		ginkgo.It("should return 404 if the product does not exist", func() {
			// Try to get a non-existent product
			nonExistentProductID := createdProductID + 999
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(nonExistentProductID))}}
			c.Request = httptest.NewRequest("GET", "/api/v1/products/"+strconv.Itoa(int(nonExistentProductID)), nil)

			// Call the GetProduct handler
			productHandler.GetProduct(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusNotFound))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.Equal("product not found"))
		})

		ginkgo.It("should return 400 for an invalid product ID", func() {
			// Try to get a product with an invalid ID
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: "invalid-id"}}
			c.Request = httptest.NewRequest("GET", "/api/v1/products/invalid-id", nil)

			// Call the GetProduct handler
			productHandler.GetProduct(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusBadRequest))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.Equal("invalid product ID"))
		})

		ginkgo.It("should return 500 if there is an internal error during product retrieval", func() {
			// Simulate a failure in GetProductByID
			mockUsecase := new(MockProductUsecase)
			productHandler := handler.NewProductHandler(mockUsecase)
			mockUsecase.On("GetProductByID", createdProductID).Return(nil, errors.New("database connection error"))

			// Create a request with a valid product ID
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(createdProductID))}}
			c.Request = httptest.NewRequest("GET", "/api/v1/products/"+strconv.Itoa(int(createdProductID)), nil)

			// Call the GetProduct handler
			productHandler.GetProduct(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusInternalServerError))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.Equal("failed to retrieve product"))
		})
	})
})
