package e2e_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"inventory_management/api/handler"
	"inventory_management/internal/entity"
	"inventory_management/internal/repository"
	"inventory_management/internal/usecase"
	"inventory_management/pkg/db"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// Mocking the Product Repository
type MockProductRepository struct {
	mock.Mock
}

type MockProductUsecase struct {
	mock.Mock
}

// CreateProduct mock method
func (m *MockProductUsecase) CreateProduct(name string) (*entity.Product, error) {
	args := m.Called(name)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// GetProductByID mock method
func (m *MockProductUsecase) GetProductByID(id uint) (*entity.Product, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// UpdateProductName mock method
func (m *MockProductUsecase) UpdateProductName(id uint, name string) (*entity.Product, error) {
	args := m.Called(id, name)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockProductRepository) Save(p *entity.Product) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockProductRepository) FindByID(id uint) (*entity.Product, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.Product), args.Error(1)
	}
	return nil, args.Error(1)
}

// Helper function to truncate tables
func TruncateTables(database *gorm.DB) {
	database.Exec("TRUNCATE TABLE products RESTART IDENTITY CASCADE;")
}

var _ = Describe("Product Handler E2E Tests (Direct Handler Calls)", func() {
	var productHandler *handler.ProductHandler
	var database *gorm.DB
	var productID uint

	BeforeEach(func() {
		// Use static configuration for the test environment
		database, _ = db.InitDB(true) // Pass 'true' to use static test configuration
		TruncateTables(database)      // Ensure tables are clean before each test

		productRepo := repository.NewPostgresProductRepository(database)
		productUsecase := usecase.NewProductUsecase(productRepo)
		productHandler = handler.NewProductHandler(productUsecase)

		// Create a product using the handler
		reqBody := map[string]string{"name": "Test Product"}
		body, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")

		// Call the CreateProduct handler
		productHandler.CreateProduct(c)

		// Print the response to check if the product is created correctly
		fmt.Printf("Create Product Response: %v\n", w.Body.String())

		// Extract the product ID from the response for future GET requests
		var response map[string]interface{}
		json.NewDecoder(w.Body).Decode(&response)
		productID = uint(response["id"].(float64))

		// Debugging to ensure productID is not zero
		fmt.Printf("Product ID: %d\n", productID)
	})

	AfterEach(func() {
		TruncateTables(database) // Clean the database after each test
	})

	Context("POST /products (direct handler call)", func() {
		It("should return 500 if the usecase returns an error during product creation", func() {
			mockUsecase := new(MockProductUsecase)
			productHandler := handler.NewProductHandler(mockUsecase)
			mockUsecase.On("CreateProduct", "Test Product").Return(nil, errors.New("usecase error"))

			reqBody := map[string]string{"name": "Test Product"}
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			productHandler.CreateProduct(c)

			Expect(w.Code).To(Equal(http.StatusInternalServerError))
			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)
			Expect(response["error"]).To(Equal("Failed to create product"))
		})

		It("should create a product successfully", func() {
			// Test successful product creation
			reqBody := map[string]string{"name": "Another Test Product"}
			body, _ := json.Marshal(reqBody)

			// Create a Gin context to simulate the request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the handler
			productHandler.CreateProduct(c)

			// Verify the response
			Expect(w.Code).To(Equal(http.StatusCreated))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)
			Expect(response["name"]).To(Equal("Another Test Product"))
			Expect(response["id"]).ShouldNot(BeZero())
		})

		It("should fail with empty name and return 422 Unprocessable Entity", func() {
			// Set up the invalid request payload (empty name)
			reqBody := map[string]string{"name": ""}
			body, _ := json.Marshal(reqBody)

			// Create a Gin context to simulate the request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the handler
			productHandler.CreateProduct(c)

			// Verify the response code (422 Unprocessable Entity)
			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			// Verify the error message in the response
			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			// Print the actual error message for debugging
			fmt.Printf("Actual error message: %v\n", response["error"])

			// Check for the expected custom error message for the "name" field
			Expect(response["errors"]).To(HaveKey("Name"))
			Expect(response["errors"].(map[string]interface{})["Name"]).To(Equal("Product name is required."))
		})

		It("should fail with a name that is too short and return 422 Unprocessable Entity", func() {
			// Set up the request payload with a short name (less than 2 characters)
			reqBody := map[string]string{"name": "A"} // Name shorter than the minimum length
			body, _ := json.Marshal(reqBody)

			// Create a Gin context to simulate the request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the handler
			productHandler.CreateProduct(c)

			// Verify the response code (422 Unprocessable Entity)
			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			// Verify the error message in the response
			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			// Print the actual error message for debugging
			fmt.Printf("Actual error message: %v\n", response["errors"])

			// Check for the expected custom error message for the "name" field
			Expect(response["errors"]).To(HaveKey("Name"))
			Expect(response["errors"].(map[string]interface{})["Name"]).To(Equal("Product name must be at least 2 characters long."))
		})

		It("should fail with a name that is too long and return 422 Unprocessable Entity", func() {
			// Set up the request payload with a name longer than 255 characters
			longName := make([]byte, 256) // Generate a string longer than 255 characters
			for i := range longName {
				longName[i] = 'A'
			}

			reqBody := map[string]string{"name": string(longName)}
			body, _ := json.Marshal(reqBody)

			// Create a Gin context to simulate the request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the handler
			productHandler.CreateProduct(c)

			// Verify the response code (422 Unprocessable Entity)
			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			// Verify the error message in the response
			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			// Print the actual error message for debugging
			fmt.Printf("Actual error message: %v\n", response["errors"])

			// Check for the expected custom error message for the "name" field
			Expect(response["errors"]).To(HaveKey("Name"))
			Expect(response["errors"].(map[string]interface{})["Name"]).To(Equal("Product name must be less than 255 characters long."))
		})

		It("should return 500 when the product cannot be saved to the database", func() {
			// Simulate a database error
			mockRepo := new(MockProductRepository)
			mockRepo.On("Save", mock.Anything).Return(errors.New("database error"))
			productUsecase := usecase.NewProductUsecase(mockRepo)
			productHandler := handler.NewProductHandler(productUsecase)

			// Create a Gin context to simulate the request
			reqBody := map[string]string{"name": "Test Product"}
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			// Call the handler
			productHandler.CreateProduct(c)

			// Verify the response code (500 Internal Server Error)
			Expect(w.Code).To(Equal(http.StatusInternalServerError))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			// Check the error message
			Expect(response["error"]).To(Equal("Failed to create product"))
		})
	})

	Context("PUT /products/:id (direct handler call)", func() {
		It("should rename the product successfully", func() {
			// Rename the product using the ID obtained from the create test
			reqBody := map[string]string{"name": "Updated Product Name"}
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set the route parameter manually
			c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(productID))}}

			c.Request = httptest.NewRequest("PUT", "/api/v1/products/"+strconv.Itoa(int(productID)), bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			productHandler.UpdateProductName(c)

			// Verify the product was renamed
			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)
			Expect(response["name"]).To(Equal("Updated Product Name"))
		})

		It("should return 404 when renaming a non-existent product", func() {
			reqBody := map[string]string{"name": "New Name"}
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set the route parameter manually
			c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(99999))}}

			c.Request = httptest.NewRequest("PUT", "/api/v1/products/99999", bytes.NewBuffer(body)) // Non-existent ID
			c.Request.Header.Set("Content-Type", "application/json")

			productHandler.UpdateProductName(c)

			Expect(w.Code).To(Equal(http.StatusNotFound))
		})

		It("should return 422 when the new name is empty", func() {
			// Try renaming the product with an empty name
			reqBody := map[string]string{"name": ""}
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("PUT", "/api/v1/products/"+strconv.Itoa(int(productID)), bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			productHandler.UpdateProductName(c)

			// Verify it returns a 422 Unprocessable Entity
			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			// Check for the expected custom error message for the "name" field
			Expect(response["errors"]).To(HaveKey("Name"))
			Expect(response["errors"].(map[string]interface{})["Name"]).To(Equal("Product name is required."))
		})

		It("should return 422 when the new name is too short", func() {
			// Try renaming the product with a name that's too short
			reqBody := map[string]string{"name": "A"} // Name shorter than the minimum length
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("PUT", "/api/v1/products/"+strconv.Itoa(int(productID)), bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			productHandler.UpdateProductName(c)

			// Verify it returns a 422 Unprocessable Entity
			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			// Check for the expected custom error message for the "name" field
			Expect(response["errors"]).To(HaveKey("Name"))
			Expect(response["errors"].(map[string]interface{})["Name"]).To(Equal("Product name must be at least 2 characters long."))
		})

		It("should return 422 when the new name exceeds the maximum length", func() {
			// Try renaming the product with a name that's too long
			longName := make([]byte, 256) // Generate a string longer than 255 characters
			for i := range longName {
				longName[i] = 'A'
			}

			reqBody := map[string]string{"name": string(longName)}
			body, _ := json.Marshal(reqBody)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("PUT", "/api/v1/products/"+strconv.Itoa(int(productID)), bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			productHandler.UpdateProductName(c)

			// Verify it returns a 422 Unprocessable Entity
			Expect(w.Code).To(Equal(http.StatusUnprocessableEntity))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			// Check for the expected custom error message for the "name" field
			Expect(response["errors"]).To(HaveKey("Name"))
			Expect(response["errors"].(map[string]interface{})["Name"]).To(Equal("Product name must be less than 255 characters long."))
		})

	})

	Context("GET /products/:id (direct handler call)", func() {
		It("should return 400 for invalid product IDs", func() {
			invalidIDs := []string{"invalid", "0"}
			for _, id := range invalidIDs {
				w := httptest.NewRecorder()
				c, _ := gin.CreateTestContext(w)
				c.Params = gin.Params{{Key: "id", Value: id}}

				productHandler.GetProduct(c)

				Expect(w.Code).To(Equal(http.StatusBadRequest))

				var response map[string]interface{}
				json.NewDecoder(w.Body).Decode(&response)
				Expect(response["error"]).To(Equal("Invalid product ID"))
			}
		})

		It("should retrieve the created product", func() {
			// Create a Gin context to simulate the request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(int(productID))}}

			// Call the handler
			productHandler.GetProduct(c)

			// Verify the response
			Expect(w.Code).To(Equal(http.StatusOK))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			// Debugging to ensure the response is correct
			fmt.Printf("Get Product Response: %v\n", w.Body.String())

			// Verify the product name is "Test Product"
			Expect(response["name"]).To(Equal("Test Product"))
		})

		It("should return 404 for non-existent product", func() {
			// Mock the repository to return a "product not found" error
			mockRepo := new(MockProductRepository)
			mockRepo.On("FindByID", uint(99999)).Return(nil, errors.New("product not found"))
			productUsecase := usecase.NewProductUsecase(mockRepo)
			productHandler := handler.NewProductHandler(productUsecase)

			// Create a Gin context to simulate the request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: "99999"}} // Non-existent ID

			// Call the handler
			productHandler.GetProduct(c)

			// Verify the response
			Expect(w.Code).To(Equal(http.StatusNotFound))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)
			Expect(response["error"]).To(Equal("Product not found")) // Match error message
		})

		It("should return 500 if the database returns an error", func() {
			// Simulate a database error when retrieving a product
			mockRepo := new(MockProductRepository)
			mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("database error"))
			productUsecase := usecase.NewProductUsecase(mockRepo)
			productHandler := handler.NewProductHandler(productUsecase)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: "1"}}

			// Call the handler
			productHandler.GetProduct(c)

			// Verify the response
			Expect(w.Code).To(Equal(http.StatusInternalServerError))

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)
			Expect(response["error"]).To(Equal("Database error"))
		})

	})
})

func TestE2E(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "E2E Product Handler Suite")
}
