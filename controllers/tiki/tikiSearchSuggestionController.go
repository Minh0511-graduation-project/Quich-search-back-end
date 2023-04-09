package tiki

import (
	"Quick-search-back-end/configs"
	"Quick-search-back-end/middlewares"
	"Quick-search-back-end/models"
	"Quick-search-back-end/responses"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var tikiSuggestionCollection = configs.GetCollection(configs.DB, "tiki search suggestions")

func GetSuggestionsByKeyword() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		keyword := r.URL.Query().Get("keyword")
		var tikiSuggestions []models.SearchSuggestion
		defer cancel()

		results, err := tikiSuggestionCollection.Find(ctx, bson.M{"keyword": keyword})

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
			var suggestion models.SearchSuggestion
			if err = results.Decode(&suggestion); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				err := json.NewEncoder(rw).Encode(response)
				if err != nil {
					return
				}
			}

			tikiSuggestions = append(tikiSuggestions, suggestion)
		}

		sort.Slice(tikiSuggestions, func(i, j int) bool {
			return tikiSuggestions[i].UpdatedAt > tikiSuggestions[j].UpdatedAt
		})

		middlewares.HandleCors(rw)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": tikiSuggestions}}
		log.Println(response.Status)
		err = json.NewEncoder(rw).Encode(response)
		if err != nil {
			return
		}
	}
}

func GetAllSuggestions() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var tikiSuggestions []models.SearchSuggestion
		defer cancel()

		results, err := tikiSuggestionCollection.Find(ctx, bson.M{})

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
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
			var suggestion models.SearchSuggestion
			if err = results.Decode(&suggestion); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				log.Println(response)
				err := json.NewEncoder(rw).Encode(response)
				if err != nil {
					return
				}
			}

			tikiSuggestions = append(tikiSuggestions, suggestion)
		}

		sort.Slice(tikiSuggestions, func(i, j int) bool {
			return tikiSuggestions[i].UpdatedAt > tikiSuggestions[j].UpdatedAt
		})

		middlewares.HandleCors(rw)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": tikiSuggestions}}
		log.Println(response.Status)
		err = json.NewEncoder(rw).Encode(response)
		if err != nil {
			return
		}
	}
}

type RequestBody struct {
	ProductID        []int  `json:"product_id"`
	ExcludedBusiness int    `json:"excluded_business"`
	PaymentModel     string `json:"payment_model"`
}

func GetTikiTopSearchByCategory() http.HandlerFunc {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		tikiTopSearchUrl := os.Getenv("TIKI_TOP_SEARCH_URL")
		data := RequestBody{
			ProductID:        []int{249220939},
			ExcludedBusiness: 157998,
			PaymentModel:     "CPC",
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}
		client := &http.Client{}
		req, err := http.NewRequest("POST", tikiTopSearchUrl, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			fmt.Println(err)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		middlewares.HandleCors(rw)
		_, err = rw.Write(body)
		if err != nil {
			return
		}
	}
}

func GetTikiTopSearchSuggestion() http.HandlerFunc {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		tikiTopSearchUrl := os.Getenv("TIKI_TOP_SEARCH_SUGGESTION")
		client := &http.Client{}
		req, err := http.NewRequest("GET", tikiTopSearchUrl, nil)
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			fmt.Println(err)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		middlewares.HandleCors(rw)
		_, err = rw.Write(body)
		if err != nil {
			return
		}
	}
}
