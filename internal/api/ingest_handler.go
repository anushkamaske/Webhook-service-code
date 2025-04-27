package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "webhook-service/internal/queue"
)

func RegisterIngestRoutes(r *gin.Engine, pub *queue.Publisher) {
    i := r.Group("/ingest")
    i.POST("/:id", func(c *gin.Context) {
        id := c.Param("id")
        raw, err := c.GetRawData()
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        pub.Publish(id, map[string]interface{}{"subscription_id": id, "body": raw, "attempt": 1})
        c.Status(http.StatusAccepted)
    })
}
