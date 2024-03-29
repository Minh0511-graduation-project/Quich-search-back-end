package shopee

import (
	"Quick-search-back-end/configs"
	"Quick-search-back-end/middlewares"
	"Quick-search-back-end/models"
	"Quick-search-back-end/responses"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"
)

var shopeeProductCollection = configs.GetCollection(configs.DB, "shopee products")

func GetProductsBySearchTerm() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		searchTerm := r.URL.Query().Get("searchTerm")
		decodedSearchTerm, err := url.QueryUnescape(searchTerm)
		if err != nil {
			return
		}
		var shopeeProduct []models.Product
		defer cancel()

		results, err := shopeeProductCollection.Find(ctx, bson.M{"searchTerm": decodedSearchTerm})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			log.Println(response)
			err := json.NewEncoder(rw).Encode(response)
			if err != nil {
				return
			}
			return
		}

		//reading from the db in an optimal way
		defer func(results *mongo.Cursor, ctx context.Context) {
			err := results.Close(ctx)
			if err != nil {
				return
			}
		}(results, ctx)
		for results.Next(ctx) {
			var singleProduct models.Product
			if err = results.Decode(&singleProduct); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				log.Println(response)
				err := json.NewEncoder(rw).Encode(response)
				if err != nil {
					return
				}
			}

			shopeeProduct = append(shopeeProduct, singleProduct)
		}

		sort.Slice(shopeeProduct, func(i, j int) bool {
			return shopeeProduct[i].UpdatedAt > shopeeProduct[j].UpdatedAt
		})

		middlewares.HandleCors(rw)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": shopeeProduct}}
		log.Println(response.Status)
		err = json.NewEncoder(rw).Encode(response)
		if err != nil {
			return
		}
	}
}

func GetAllProducts() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var shopeeProduct []models.Product
		defer cancel()

		results, err := shopeeProductCollection.Find(ctx, bson.M{})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
			log.Println(response)
			err := json.NewEncoder(rw).Encode(response)
			if err != nil {
				return
			}
			return
		}

		//reading from the db in an optimal way
		defer func(results *mongo.Cursor, ctx context.Context) {
			err := results.Close(ctx)
			if err != nil {
				return
			}
		}(results, ctx)
		for results.Next(ctx) {
			var singleProduct models.Product
			if err = results.Decode(&singleProduct); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				log.Println(response)
				err := json.NewEncoder(rw).Encode(response)
				if err != nil {
					return
				}
			}

			shopeeProduct = append(shopeeProduct, singleProduct)
		}

		sort.Slice(shopeeProduct, func(i, j int) bool {
			return shopeeProduct[i].UpdatedAt > shopeeProduct[j].UpdatedAt
		})

		middlewares.HandleCors(rw)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": shopeeProduct}}
		log.Println(response.Status)
		err = json.NewEncoder(rw).Encode(response)
		if err != nil {
			return
		}
	}
}
