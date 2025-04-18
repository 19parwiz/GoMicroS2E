package http

import (
	"github.com/19parwiz/inventory/internal/adapter/http/handler"
)

type ProductHandler interface {
	handler.ProductUseCase
}
