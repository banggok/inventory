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

	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("CreateProduct E2E Tests", func() {
	var productHandler *handler.ProductHandler
	var database *gorm.DB
	var sqlDB *sql.DB

	ginkgo.BeforeEach(func() {
		// Use static configuration for the test environment
		database, sqlDB = db.InitDB(true)
		TruncateTables(database) // Ensure tables are clean before each test

		productRepo := repository.NewPostgresProductRepository(database)
		productUsecase := usecase.NewProductUsecase(productRepo)
		productHandler = handler.NewProductHandler(productUsecase)
	})

	ginkgo.AfterEach(func() {
		TruncateTables(database) // Clean the database after each test
		sqlDB.Close()
	})

	ginkgo.Context("POST /products", func() {
		ginkgo.It("should create a product successfully", func() {
			// Set up the request body for a valid product creation
			reqBody := map[string]string{"name": "New Product"}
			body, _ := json.Marshal(reqBody)

			// Create a Gin context to simulate the request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the CreateProduct handler
			productHandler.CreateProduct(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusCreated))

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["name"]).To(gomega.Equal("New Product"))
			gomega.Expect(response["id"]).ShouldNot(gomega.BeZero())
		})

		ginkgo.It("should fail with validation error if the name is empty", func() {
			// Set up the request body with an empty name
			reqBody := map[string]string{"name": ""}
			body, _ := json.Marshal(reqBody)

			// Create a Gin context to simulate the request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the CreateProduct handler
			productHandler.CreateProduct(c)

			// Verify the response (422 Unprocessable Entity)
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusUnprocessableEntity))

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.HaveKey("Name"))
			gomega.Expect(response["errors"].(map[string]interface{})["Name"]).To(gomega.Equal("Product name is required."))
		})

		ginkgo.It("should return 500 if the product creation fails at the use case level", func() {
			// Mock the use case to return an error during product creation
			mockUsecase := new(MockProductUsecase)
			productHandler := handler.NewProductHandler(mockUsecase)
			mockUsecase.On("CreateProduct", "Failing Product").Return(nil, errors.New("usecase error"))

			// Set up the request body
			reqBody := map[string]string{"name": "Failing Product"}
			body, _ := json.Marshal(reqBody)

			// Create a Gin context to simulate the request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the CreateProduct handler
			productHandler.CreateProduct(c)

			// Verify the response (500 Internal Server Error)
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusInternalServerError))

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.Equal("failed to create product"))
		})

		ginkgo.It("should return 400 Bad Request for invalid JSON body", func() {
			// Create an invalid JSON request body
			invalidJSON := []byte(`{"name":`)

			// Create a Gin context to simulate the request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(invalidJSON))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the CreateProduct handler
			productHandler.CreateProduct(c)

			// Verify the response (400 Bad Request)
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusBadRequest))

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.Equal("invalid JSON format"))

		})

	})
})
