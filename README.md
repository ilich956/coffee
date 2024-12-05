# Coffee Management System


## Overview

Coffee Management System is a backend application for a coffee shop management system. This project is designed to simulate real-world systems used by businesses to streamline operations such as order management, inventory tracking, and menu updates.

## Build and Run
- Clone the repository.
- Navigate to the project root.
- Build the application:
```sh
go build -o coffee ./cmd/main.go
```
Run the server:
```sh
./coffee
```

## Features
- Order Management: Create, update, delete, and close orders.
- Inventory Tracking: Monitor and update ingredient stock levels.
- Menu Management: Add, edit, and remove menu items with associated ingredients and prices.
- Aggregations: Calculate metrics such as total sales and popular menu items.
- Error Handling: Ensure graceful handling of errors with meaningful responses.
- Logging: Use Go’s log/slog package for recording application events.ices as needed, and keep the offerings up to date.



### API Endpoints

Implement the following RESTful API endpoints:

- **Orders:**

  - `POST /orders`: Create a new order.
  - `GET /orders`: Retrieve all orders.
  - `GET /orders/{id}`: Retrieve a specific order by ID.
  - `PUT /orders/{id}`: Update an existing order.
  - `DELETE /orders/{id}`: Delete an order.
  - `POST /orders/{id}/close`: Close an order.

- **Menu Items:**

  - `POST /menu`: Add a new menu item.
  - `GET /menu`: Retrieve all menu items.
  - `GET /menu/{id}`: Retrieve a specific menu item.
  - `PUT /menu/{id}`: Update a menu item.
  - `DELETE /menu/{id}`: Delete a menu item.

- **Inventory:**

  - `POST /inventory`: Add a new inventory item.
  - `GET /inventory`: Retrieve all inventory items.
  - `GET /inventory/{id}`: Retrieve a specific inventory item.
  - `PUT /inventory/{id}`: Update an inventory item.
  - `DELETE /inventory/{id}`: Delete an inventory item.

- **Aggregations:**

  - `GET /reports/total-sales`: Get the total sales amount.
  - `GET /reports/popular-items`: Get a list of popular menu items.

**Examples:**

- **Create Order Request:**
```http 
POST /orders
Content-Type: application/json

{
  "customer_name": "John Doe",
  "items": [
    {
      "product_id": "espresso",
      "quantity": 2
    },
    {
      "product_id": "croissant",
      "quantity": 1
    }
  ]
}
```

- **Total Sales Aggregation Response:**
```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "total_sales": 1500.50
}
```


#### Examples of JSON Files:

1. `orders.json`:

```json
[
  {
    "order_id": "order123",
    "customer_name": "Alice Smith",
    "items": [
      {
        "product_id": "latte",
        "quantity": 2
      },
      {
        "product_id": "muffin",
        "quantity": 1
      }
    ],
    "status": "open",
    "created_at": "2023-10-01T09:00:00Z"
  },
  {
    "order_id": "order124",
    "customer_name": "Bob Johnson",
    "items": [
      {
        "product_id": "espresso",
        "quantity": 1
      }
    ],
    "status": "closed",
    "created_at": "2023-10-01T09:30:00Z"
  }
]
```

2. `menu_items.json`:
```json
[
  {
    "product_id": "latte",
    "name": "Caffe Latte",
    "description": "Espresso with steamed milk",
    "price": 3.50,
    "ingredients": [
      {
        "ingredient_id": "espresso_shot",
        "quantity": 1
      },
      {
        "ingredient_id": "milk",
        "quantity": 200
      }
    ]
  },
  {
    "product_id": "muffin",
    "name": "Blueberry Muffin",
    "description": "Freshly baked muffin with blueberries",
    "price": 2.00,
    "ingredients": [
      {
        "ingredient_id": "flour",
        "quantity": 100
      },
      {
        "ingredient_id": "blueberries",
        "quantity": 20
      },
      {
        "ingredient_id": "sugar",
        "quantity": 30
      }
    ]
  },
  {
    "product_id": "espresso",
    "name": "Espresso",
    "description": "Strong and bold coffee",
    "price": 2.50,
    "ingredients": [
      {
        "ingredient_id": "espresso_shot",
        "quantity": 1
      }
    ]
  }
]
```

**Note:** The ingredients field in each menu item lists the ingredients required to prepare that item. The quantity is specified in units appropriate for the ingredient (e.g., grams, milliliters).

3. `inventory.json`:

```json
[
  {
    "ingredient_id": "espresso_shot",
    "name": "Espresso Shot",
    "quantity": 500, // Number of shots
    "unit": "shots"
  },
  {
    "ingredient_id": "milk",
    "name": "Milk",
    "quantity": 5000, // In milliliters
    "unit": "ml"
  },
  {
    "ingredient_id": "flour",
    "name": "Flour",
    "quantity": 10000, // In grams
    "unit": "g"
  },
  {
    "ingredient_id": "blueberries",
    "name": "Blueberries",
    "quantity": 2000,  // In grams
    "unit": "g"
  },
  {
    "ingredient_id": "sugar",
    "name": "Sugar",
    "quantity": 5000, // In grams
    "unit": "g"
  }
]
```
