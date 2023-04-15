package searchCount

import (
	"Quick-search-back-end/configs"
	"Quick-search-back-end/middlewares"
	"Quick-search-back-end/models"
	"Quick-search-back-end/responses"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var tikiSuggestionCollection = configs.GetCollection(configs.DB, "search keyword count")

func GetCountByKeyword() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		keyword := r.URL.Query().Get("keyword")
		decodeKeyword, err := url.QueryUnescape(keyword)
		if err != nil {
			return
		}
		site := r.URL.Query().Get("site")
		var wordCountStat []models.KeywordCount
		defer cancel()

		results, err := tikiSuggestionCollection.Find(ctx,
			bson.M{
				"keyword": decodeKeyword,
				"site":    site,
			})

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
			var stat models.KeywordCount
			if err = results.Decode(&stat); err != nil {
				rw.WriteHeader(http.StatusInternalServerError)
				response := responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}}
				err := json.NewEncoder(rw).Encode(response)
				if err != nil {
					return
				}
			}

			wordCountStat = append(wordCountStat, stat)
		}

		middlewares.HandleCors(rw)

		rw.WriteHeader(http.StatusOK)
		response := responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": wordCountStat}}
		log.Println(response.Status)
		err = json.NewEncoder(rw).Encode(response)
		if err != nil {
			return
		}
	}
}

func GetShopeeTopSearch() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		findOptions := options.Find()
		topDisplay, err := strconv.ParseInt(r.URL.Query().Get("topDisplay"), 10, 64)
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
		findOptions.SetSort(bson.M{"count": -1})
		findOptions.SetLimit(topDisplay)
		filter := bson.M{"site": "shopee"}
		cursor, err := tikiSuggestionCollection.Find(ctx, filter, findOptions)
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
		defer func(cursor *mongo.Cursor, ctx context.Context) {
			err := cursor.Close(ctx)
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
		}(cursor, ctx)

		var keywordCounts []models.KeywordCount
		for cursor.Next(ctx) {
			var keywordCount models.KeywordCount
			err := cursor.Decode(&keywordCount)
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
			keywordCounts = append(keywordCounts, keywordCount)
		}
		err = cursor.Err()
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
		response := map[string]interface{}{
			"data": keywordCounts,
		}
		userResponse := map[string]interface{}{
			"status":  http.StatusOK,
			"message": "success",
			"data":    response,
		}

		middlewares.HandleCors(rw)

		rw.WriteHeader(http.StatusOK)
		err = json.NewEncoder(rw).Encode(userResponse)
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
	}
}
