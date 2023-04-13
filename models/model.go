package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id         primitive.ObjectID `json:"_id,omitempty" validate:"required"`
	Name       string             `json:"name,omitempty" validate:"required"`
	ImageUrl   string             `json:"imageUrl,omitempty" validate:"required"`
	Price      string             `json:"price,omitempty" validate:"required"`
	SearchTerm string             `json:"searchTerm,omitempty" validate:"required"`
	Site       string             `json:"site,omitempty" validate:"required"`
	UpdatedAt  float64            `json:"updatedAt,omitempty" validate:"required"`
	ProductUrl string             `json:"productUrl,omitempty" validate:"required"`
}

type SearchSuggestion struct {
	Id          primitive.ObjectID `json:"id,omitempty" validate:"required"`
	Keyword     string             `json:"keyword,omitempty" validate:"required"`
	Site        string             `json:"site,omitempty" validate:"required"`
	Suggestions []string           `json:"suggestions,omitempty" validate:"required"`
	UpdatedAt   float64            `json:"updatedAt,omitempty" validate:"required"`
}

type KeywordCount struct {
	Id        primitive.ObjectID `json:"_id,omitempty" validate:"required"`
	Keyword   string             `json:"keyword,omitempty" validate:"required"`
	Site      string             `json:"site,omitempty" validate:"required"`
	Count     float64            `json:"count,omitempty" validate:"required"`
	UpdatedAt float64            `json:"updatedAt,omitempty" validate:"required"`
}

type GetTikiTopSearchByCategoryRequestBody struct {
	ProductID        []int  `json:"product_id"`
	ExcludedBusiness int    `json:"excluded_business"`
	PaymentModel     string `json:"payment_model"`
}
