package myerrors

import "errors"

var (
	ErrNotFound      = errors.New("Not found")                     // 404
	ErrFailOpenJson  = errors.New("Error with oppening json file") // 5##
	ErrFailUnmarshal = errors.New("Failed to unmarshal")
	ErrFailMarshal   = errors.New("Failed to marshal") // 5##
	ErrInvalidJson   = errors.New("Invalid JSON")
	ErrFailWrite     = errors.New("Failed write")

	ErrOrderClosed = errors.New("Order is already closed") //
	ErrEmptyOrder  = errors.New("After validating of your order - it became empty")
	ErrIDExist     = errors.New("ID already exists")
	ErrAbsentItem  = errors.New("No such items in the menu")
	ErrNoItems     = errors.New("No items")

	ErrIdRequired           = errors.New("ID field is required")
	ErrNameRequired         = errors.New("Name field is required")
	ErrDescriptionRequired  = errors.New("Description field is required")
	ErrPriceRequired        = errors.New("Price field is required")
	ErrIngredientsRequired  = errors.New("Ingredients field is required")
	ErrInvalidQuantity      = errors.New("Quantity field is invalid")
	ErrUnitRequired         = errors.New("Unit field is required")
	ErrItemsRequired        = errors.New("Items field are required")
	ErrNotEnoughIngridients = errors.New("Not enough ingridients")
)
