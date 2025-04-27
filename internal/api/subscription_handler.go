package api

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "webhook-service/internal/model"
    "webhook-service/internal/store/postgres"
)

func RegisterSubscriptionRoutes(r *gin.Engine, repo *postgres.SubscriptionRepo) {
    s := r.Group("/subscriptions")
    s.POST("", func(c *gin.Context) {
        var sub model.Subscription
        if err := c.ShouldBindJSON(&sub); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        repo.Create(&sub)
        c.JSON(http.StatusCreated, sub)
    })

    s.GET("/:id", func(c *gin.Context) {
        id := c.Param("id")
        sub, err := repo.GetByID(id)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
            return
        }
        c.JSON(http.StatusOK, sub)
    })

    s.GET("", func(c *gin.Context) {
        subs, _ := repo.List()
        c.JSON(http.StatusOK, subs)
    })

    s.PUT("/:id", func(c *gin.Context) {
        id := c.Param("id")
        var sub model.Subscription
        if err := c.ShouldBindJSON(&sub); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        sub.ID = id
        repo.Update(&sub)
        c.JSON(http.StatusOK, sub)
    })

    s.DELETE("/:id", func(c *gin.Context) {
        id := c.Param("id")
        repo.Delete(id)
        c.Status(http.StatusNoContent)
    })
}
