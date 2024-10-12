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

var _ = ginkgo.Describe("UpdateProductName E2E Tests", func() {
	var productHandler *handler.ProductHandler
	var database *gorm.DB
	var sqlDB *sql.DB

	var createdProductID uint

	// Helper function to create a product before testing updates
	CreateProduct := func() uint {
		// Create a new product via the API
		reqBody := map[string]string{"name": "Initial Product"}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		productHandler.CreateProduct(c)

		// Extract the created product's ID from the response
		var response map[string]interface{}
		_ = json.NewDecoder(w.Body).Decode(&response)
		return uint(response["id"].(float64)) // Return the created product ID
	}

	ginkgo.BeforeEach(func() {
		// Initialize test environment
		database, sqlDB = db.InitDB(true) // Assuming `true` loads the test environment
		TruncateTables(database)          // Clean up before each test

		// Initialize the handler with a real database connection (no mocks)
		productRepo := repository.NewPostgresProductRepository(database)
		productUsecase := usecase.NewProductUsecase(productRepo)
		productHandler = handler.NewProductHandler(productUsecase)

		// Create a product before testing updates
		createdProductID = CreateProduct()
	})

	ginkgo.AfterEach(func() {
		TruncateTables(database) // Clean up after each test
		sqlDB.Close()
	})

	ginkgo.Context("PUT /products/:id", func() {
		ginkgo.It("should update a product's name successfully", func() {
			// Create the request body
			reqBody := map[string]string{"name": "Updated Name"}
			body, _ := json.Marshal(reqBody)

			// Create a new request with the product ID
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(createdProductID))}}
			c.Request = httptest.NewRequest("PUT", "/api/v1/products/"+strconv.Itoa(int(createdProductID)), bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the handler to update the product
			productHandler.UpdateProductName(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusOK))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["name"]).To(gomega.Equal("Updated Name"))
		})

		ginkgo.It("should return 422 if the request body validation fails", func() {
			// Create an invalid request body (empty name)
			reqBody := map[string]string{"name": ""}
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(createdProductID))}}
			c.Request = httptest.NewRequest("PUT", "/api/v1/products/"+strconv.Itoa(int(createdProductID)), bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the handler to update the product
			productHandler.UpdateProductName(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusUnprocessableEntity))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.HaveKey("Name"))
		})

		ginkgo.It("should return 404 if the product is not found", func() {
			// Use a non-existent product ID
			nonExistentProductID := createdProductID + 999

			reqBody := map[string]string{"name": "Updated Name"}
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(nonExistentProductID))}}
			c.Request = httptest.NewRequest("PUT", "/api/v1/products/"+strconv.Itoa(int(nonExistentProductID)), bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the handler to update the product
			productHandler.UpdateProductName(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusNotFound))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.Equal("product not found"))
		})

		ginkgo.It("should return 400 if the product ID is invalid", func() {
			// Create a valid JSON request body, but the product ID will be invalid
			reqBody := map[string]string{"name": "Updated Name"}
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: "invalid-id"}} // Invalid product ID
			c.Request = httptest.NewRequest("PUT", "/api/v1/products/invalid-id", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the handler to update the product
			productHandler.UpdateProductName(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusBadRequest))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.Equal("invalid product ID"))
		})

		ginkgo.It("should return 500 if there's an internal server error", func() {
			// Mock the use case to simulate an internal server error
			mockUsecase := new(MockProductUsecase)
			productHandler := handler.NewProductHandler(mockUsecase)

			// Simulate the use case returning an error
			mockUsecase.On("UpdateProductName", createdProductID, "Error Case").Return(nil, errors.New("internal server error"))

			// Create the request body
			reqBody := map[string]string{"name": "Error Case"}
			body, _ := json.Marshal(reqBody)

			// Create a new request with the product ID
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(createdProductID))}}
			c.Request = httptest.NewRequest("PUT", "/api/v1/products/"+strconv.Itoa(int(createdProductID)), bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the handler
			productHandler.UpdateProductName(c)

			// Verify the response
			gomega.Expect(w.Code).To(gomega.Equal(http.StatusInternalServerError))
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(response["errors"]).To(gomega.Equal("failed to update product"))
		})

	})

	ginkgo.It("should return 400 if there is an error while reading the request body", func() {
		// Set up the request body with invalid JSON
		invalidJSON := []byte(`{"name":`)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(createdProductID))}}
		c.Request = httptest.NewRequest("PUT", "/api/v1/products/"+strconv.Itoa(int(createdProductID)), bytes.NewBuffer(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the handler to update the product name
		productHandler.UpdateProductName(c)

		// Verify the response
		gomega.Expect(w.Code).To(gomega.Equal(http.StatusBadRequest))
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		gomega.Expect(response["errors"]).To(gomega.Equal("invalid JSON format")) // Adjust the error message if necessary
	})

})
