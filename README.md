# Outfiro API

Outfiro is a textile e-commerce platform built using Go, Gin, PostgreSQL, JWT authentication, and Go Maroto for generating invoices. The platform allows users to browse products, manage their profiles, place orders, and utilize a referral-based wallet system for discounts and purchases. The API follows the MVC architecture, ensuring a modular and scalable approach.

## Key Features

### User Management
- **Signup & Login:** Users can register and log in using email/password or Google authentication.
- **Wallet System:** Users earn 100 points when they use a referral code, which is converted into wallet balance for purchases.
- **Cart & Wishlist:** Users can add, update, or remove products from their cart and wishlist.

### Product & Category Management
- **Product Listing:** Users can browse products with ratings and reviews.
- **Categories:** Products are categorized for better search and filtering.

### Order Processing
- **Order Placement:** Users can place orders and apply discounts using their wallet balance.
- **Order Invoice:** Invoices are generated using Go Maroto.
- **Order Tracking:** Users can track their order status from processing to delivery.

### Referral System & Wallet
- **Referral Code:** New users who enter a referral code during signup receive 100 points.
- **Wallet Conversion:** Referral points are converted into wallet balance for future purchases.

### Admin Panel
- **Admin Signup & Login:** Secure authentication for admin access.
- **Product & Category Management:** Admins can add, update, and delete products and categories.
- **Sales Reports & Graphs:** Admins can view sales analytics, including reports in PDF format using Go Maroto.

## Technology Stack

- **Backend:** Go (Gin framework)
- **Database:** PostgreSQL
- **Authentication:** JWT-based authentication
- **Payment Processing:** Integrated wallet system for using referral points
- **Invoice Generation:** Go Maroto package for PDF invoices

## Installation

To set up the project locally, follow these steps:

### Clone the Repository
```sh
git clone https://github.com/your-username/Outfiro-API.git
cd Outfiro-API
```

### Set Up Environment Variables
Create a `.env` file in the root directory and add the following:

```env
SERVER_IP=localhost:8080
DB_USER=your_database_username
DB_PASSWORD=your_database_password
DB_NAME=your_database_name
JWT_SECRET=your_jwt_secret_key
```

### Install Dependencies
```sh
go mod tidy
```

### Run the Application
```sh
go run .
```

## API Documentation
Detailed API documentation is available [here](#).
