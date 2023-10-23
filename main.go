package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	rueidis_store "github.com/eko/gocache/store/rueidis/v4"
	"github.com/gin-gonic/gin"
	"github.com/redis/rueidis"
)

var cacheManager *cache.Cache[string]

func getHandler(c *gin.Context) {
	key := c.Param("key")
	value, err := cacheManager.Get(context.Background(), key)
	if err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	c.JSON(http.StatusOK, value)
}

func setHandler(c *gin.Context) {
	key := c.Param("key")
	val := c.Param("val")
	if err := cacheManager.Set(context.Background(), key, val); err != nil {
		c.JSON(http.StatusOK, err)
		return
	}
	c.JSON(http.StatusOK, key)
}

func main() {

	redisHost := os.Getenv("rhost")
	client, err := rueidis.NewClient(rueidis.ClientOption{InitAddress: []string{redisHost + ":6379"}})
	if err != nil {
		panic(err)
	}

	cacheManager = cache.New[string](rueidis_store.NewRueidis(
		client,
		store.WithExpiration(15*time.Minute),
		store.WithClientSideCaching(15*time.Minute)),
	)

	router := gin.New()

	router.GET("/get/:key", getHandler)
	router.GET("/set/:key/:val", setHandler)

	router.Run(":8080")
}
