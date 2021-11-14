package main

import (
	"net/http"

	"github.com/Soontao/similar"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	sim := similar.NewSimilar(10000)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "similar",
			"status":  http.StatusOK,
		})
	})

	r.POST("/remember", func(c *gin.Context) {
		payload := &RememberPayload{}
		if err := c.BindJSON(payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "payload can not be accepted",
			})
			return
		}
		for _, set := range payload.Sentences {
			sim.Remember(set)
		}
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Remembered",
		})
	})

	r.POST("/find", func(c *gin.Context) {
		payload := &FindPayload{}
		if err := c.BindJSON(payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "payload can not be accepted",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "executed",
			"results": ToFindResults(sim.FindSimilar(payload.Sentence, payload.Threshold)),
		})
	})

	r.POST("/findone", func(c *gin.Context) {
		payload := &FindOnePayload{}
		if err := c.BindJSON(payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "payload can not be accepted",
			})
			return
		}
		r := ToFindResult(sim.FindMostSimilar(payload.Sentence))
		c.JSON(http.StatusOK, gin.H{
			"message": "executed",
			"results": r,
		})
	})

	r.Run()

}
