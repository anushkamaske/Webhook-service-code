package mongo

import (
    "context"
    "time"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "webhook-service/internal/model"
)

type LogRepo struct {
    col *mongo.Collection
}

func NewLogRepo(client *mongo.Client, dbName string) *LogRepo {
    col := client.Database(dbName).Collection("delivery_logs")
    col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
        Keys:    bson.D{{Key: "timestamp", Value: 1}},
        Options: options.Index().SetExpireAfterSeconds(259200),
    })
    return &LogRepo{col: col}
}

func (r *LogRepo) Insert(entry *model.DeliveryLog) error {
    _, err := r.col.InsertOne(context.Background(), entry)
    return err
}

func (r *LogRepo) GetBySubscription(subID string, limit int64) ([]model.DeliveryLog, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    cur, err := r.col.Find(ctx, bson.M{"subscription_id": subID}, options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(limit))
    if err != nil {
        return nil, err
    }
    var logs []model.DeliveryLog
    if err := cur.All(ctx, &logs); err != nil {
        return nil, err
    }
    return logs, nil
}
