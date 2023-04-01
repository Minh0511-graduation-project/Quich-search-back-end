package searchCount

import (
	"Quick-search-back-end/configs"
	"Quick-search-back-end/middlewares"
	"Quick-search-back-end/models"
	"Quick-search-back-end/responses"
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"net/url"
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
