package database

import (
	"errors"
)

var (
	ErrRecordNotFound     = errors.New("can`t find the produt")
	ErrCantDecodeProducts = errors.New("can't find (decode) products")
	ErrUserIdIsNotValid   = errors.New("this user is not valid")
	ErrCantUpdateUser     = errors.New("can't add this product to the cart")
	ErrCantRemoveItemCart = errors.New("can't remove this item from the cart")
	ErrCantGetItem        = errors.New("can't get the item from the cart")
	ErrCantBuyCartItem    = errors.New("cannot update the purchcase")
)

func AddProductToCart() {

}

func RemoveCartItem() {

}

func BuyItemFromCart() {

}

func InstantBuyer() {

}
