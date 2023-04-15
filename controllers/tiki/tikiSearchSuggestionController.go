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

func GetTikiTopSearchByCategory() http.HandlerFunc {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			middlewares.HandleCors(rw)
			return
		}
		tikiTopSearchUrl := os.Getenv("TIKI_TOP_SEARCH_URL")
		topDisplay := r.URL.Query().Get("topDisplay")
		tikiTopSearchUrl = tikiTopSearchUrl + topDisplay
		var data models.GetTikiTopSearchByCategoryRequestBody
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		client := &http.Client{}
		jsonData, err := json.Marshal(data)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		req, err := http.NewRequest("POST", tikiTopSearchUrl, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
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
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
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
