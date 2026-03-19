# Coffee App - Postman Collection

This directory contains a complete Postman collection for testing the Coffee Shop Management API.

## Files

- **Coffee-API.postman_collection.json** - The complete Postman collection with all API endpoints

## Installation & Setup

### 1. Import the Collection into Postman

- Open Postman
- Click the **Import** button (top left)
- Select **File** tab and choose `Coffee-API.postman_collection.json`
- Click **Import**

### 2. Configure the Environment

The collection uses a variable `{{host}}` which defaults to `localhost:8080`.

To change the host:
1. Click on the collection name "Coffee App API"
2. Go to the **Variables** tab
3. Change the `host` variable value to your server address

**Default value:** `localhost:8080`

## API Endpoints Overview

### Menu Management
- `GET /menu` - Get all menu items
- `GET /menu/{id}` - Get a specific menu item
- `POST /menu` - Create a new menu item
- `PUT /menu/{id}` - Update a menu item
- `DELETE /menu/{id}` - Delete a menu item

**Example Menu Item:**
```json
{
  "product_id": "espresso-1",
  "name": "Espresso",
  "description": "Strong black coffee",
  "price": 2.50,
  "ingredients": [
    {
      "ingredient_id": "coffee-beans",
      "quantity": 20
    }
  ]
}
```

### Order Management
- `GET /orders` - Get all orders
- `GET /orders/{id}` - Get a specific order
- `POST /orders` - Create a new order
- `PUT /orders/{id}` - Update an order
- `POST /orders/{id}/close` - Mark an order as closed/completed
- `DELETE /orders/{id}` - Delete an order

**Example Order:**
```json
{
  "order_id": "order-1",
  "customer_name": "John Doe",
  "items": [
    {
      "product_id": "espresso-1",
      "quantity": 2
    }
  ],
  "status": "open",
  "total_price": 5.50,
  "created_at": "2024-01-01T10:00:00Z",
  "updated_at": "2024-01-01T10:05:00Z"
}
```

### Inventory Management
- `GET /inventory` - Get all inventory items
- `GET /inventory/{id}` - Get a specific inventory item
- `POST /inventory` - Create a new inventory item
- `PUT /inventory/{id}` - Update an inventory item
- `DELETE /inventory/{id}` - Delete an inventory item

**Example Inventory Item:**
```json
{
  "ingredient_id": "coffee-beans",
  "name": "Coffee Beans",
  "quantity": 500,
  "unit": "grams",
  "threshold": 100
}
```

### Reports & Analytics
- `GET /reports/total-sales` - Get total sales from completed orders
- `GET /reports/popular-items` - Get most popular menu items

**Example Total Sales Response:**
```json
{
  "total_sale": 250.75
}
```

**Example Popular Items Response:**
```json
[
  {
    "product_id": "espresso-1",
    "quantity": 25
  },
  {
    "product_id": "latte-1",
    "quantity": 18
  }
]
```

## Testing Workflow

### Quick Test Sequence

1. **Add Menu Items** (Setup)
   - Create Espresso (POST /menu)
   - Create Latte (POST /menu)
   - Create another drink (POST /menu)

2. **Add Inventory** (Setup)
   - Create Coffee Beans (POST /inventory)
   - Create Milk (POST /inventory)

3. **Create and Manage Orders** (Main Test)
   - Create an order (POST /orders)
   - View the order (GET /orders/{id})
   - Update the order (PUT /orders/{id})
   - Close the order (POST /orders/{id}/close)

4. **View Analytics** (Reports)
   - Check total sales (GET /reports/total-sales)
   - Check popular items (GET /reports/popular-items)

5. **Cleanup** (Optional)
   - Delete orders (DELETE /orders/{id})
   - Delete menu items (DELETE /menu/{id})
   - Delete inventory items (DELETE /inventory/{id})

## Notes

- All timestamps should be in ISO 8601 format (e.g., `2024-01-01T10:00:00Z`)
- Price fields are in decimal format (e.g., 2.50)
- Quantity fields are numbers (integers or decimals)
- The server validates required fields in requests
- HTTP Status Codes:
  - 200 OK - Successful request
  - 404 Not Found - Resource not found
  - 500 Internal Server Error - Server error

## Starting the API Server

Make sure the Coffee App API server is running:

```bash
cd /path/to/coffee
go run ./cmd/main.go
```

Or with Docker:
```bash
docker-compose up
```

The server will start on `http://localhost:8080` by default.

## Troubleshooting

- **Connection refused**: Make sure the API server is running
- **404 errors**: Check that the resource ID is correct
- **400/500 errors**: Review the request body format and required fields

For more information, refer to the main project README.
