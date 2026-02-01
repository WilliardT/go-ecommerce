package models

import (
	"time"

	"github.com/google/uuid"
)

// postgreSQL database

type User struct {
	ID              uuid.UUID    `json:"id" db:"id"`
	First_Name      *string      `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name       *string      `json:"last_name" validate:"required,min=2,max=30"`
	Password        *string      `json:"password" validate:"required,min=6"`
	Email           *string      `json:"email" validate:"email,required"`
	Phone           *string      `json:"phone" validate:"required"`
	Token           *string      `json:"token"`
	Refresh_Token   *string      `json:"refresh_token"`
	Created_At      time.Time    `json:"created_at"`
	Updated_At      time.Time    `json:"updated_at"`
	User_ID         string       `json:"user_id"`
	UserCart        []PoductUser `json:"usercart"`
	Address_Details []Address    `json:"address"`
	Order_Status    []Order      `json:"order_Status"`
}
type Product struct {
	Product_ID   uuid.UUID `json:"product_id" db:"product_id"`
	Product_Name *string   `json:"product_name" db:"product_name"`
	Price        *uint64   `json:"price" db:"price"`
	Rating       *uint8    `json:"rating" db:"rating"`
	Image        *string   `json:"image" db:"image"`
}

type PoductUser struct {
	Product_ID   uuid.UUID `json:"product_id" db:"product_id"`
	Product_Name *string   `json:"product_name" db:"product_name"`
	Price        int       `json:"price" db:"price"`
	Rating       *uint     `json:"rating" db:"rating"`
	Image        *string   `json:"image" db:"image"`
}

// товар в корзине с деталями
type CartItem struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Price       uint64    `json:"price"`
	Rating      *uint8    `json:"rating"`
	Image       *string   `json:"image"`
	Quantity    int       `json:"quantity"`
}

type Address struct {
	Addres_ID uuid.UUID `json:"address_id" db:"address_id"`
	House     *string   `json:"house_name" db:"house_name"`
	Street    *string   `json:"street_name" db:"street_name"`
	City      *string   `json:"city_name" db:"city_name"`
	Pincode   *string   `json:"pincode_name" db:"pincode_name"`
	State     *string   `json:"state_name" db:"state_name"`
}

type Order struct {
	Order_ID       uuid.UUID    `json:"order_id" db:"order_id"`
	Order_Cart     []PoductUser `json:"order_cart" db:"order_cart"`
	Ordered_At     time.Time    `json:"ordered_at" db:"ordered_at"`
	Price          int          `json:"price" db:"price"`
	Discount       *int         `json:"discount" db:"discount"`
	Payment_Method Payment      `json:"payment_method" db:"payment_method"`
}

type Payment struct {
	Digital bool
	COD     bool
}
