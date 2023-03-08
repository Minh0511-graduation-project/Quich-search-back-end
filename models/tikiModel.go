package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TikiProduct struct {
	Id         primitive.ObjectID `json:"_id,omitempty" validate:"required"`
	Name       string             `json:"name,omitempty" validate:"required"`
	ImageUrl   string             `json:"imageUrl,omitempty" validate:"required"`
	Price      string             `json:"price,omitempty" validate:"required"`
	SearchTerm string             `json:"searchTerm,omitempty" validate:"required"`
	Site       string             `json:"site,omitempty" validate:"required"`
}

type TikiSearchSuggestion struct {
	Id          primitive.ObjectID `json:"id,omitempty" validate:"required"`
	Keyword     string             `json:"keyword,omitempty" validate:"required"`
	Site        string             `json:"site,omitempty" validate:"required"`
	Suggestions []string           `json:"suggestions,omitempty" validate:"required"`
}
