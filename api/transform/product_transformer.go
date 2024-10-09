// /api/handler/transformer/product_transformer.go
package transformer

import (
	"inventory_management/api/handler/dto"
	"inventory_management/internal/entity"
)

// TransformProductEntityToResponse transforms an entity.Product to a dto.ProductResponse
func TransformProductEntityToResponse(p *entity.Product) *dto.ProductResponse {
	return &dto.ProductResponse{
		ID:   p.ID(),
		Name: p.Name(),
		SKU:  p.SKU(),
	}
}
