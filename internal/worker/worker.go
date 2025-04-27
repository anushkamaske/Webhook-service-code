package worker

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/streadway/amqp"
    "webhook-service/internal/cache"
    "webhook-service/internal/config"
    "webhook-service/internal/model"
    "webhook-service/internal/queue"
    "webhook-service/internal/store/postgres"
    "webhook-service/internal/store/mongo"
)

type Worker struct {
    subRepo *postgres.SubscriptionRepo
    logRepo *mongo.LogRepo
    pub     *queue.Publisher
    cache   *cache.RedisClient
    cfg     *config.Config
}

func NewWorker(sr *postgres.SubscriptionRepo, lr *mongo.LogRepo, pub *queue.Publisher, c *cache.RedisClient, cfg *config.Config) *Worker {
    return &Worker{subRepo: sr, logRepo: lr, pub: pub, cache: c, cfg: cfg}
}

func (w *Worker) Start() {
    msgs, err := w.pub.Channel.Consume(w.pub.QueueName, "", false, false, false, false, nil)
    if err != nil {
        log.Fatalf("consume: %v", err)
    }
    for msg := range msgs {
        var job struct {
            SubscriptionID string          `json:"subscription_id"`
            Payload        json.RawMessage `json:"body"`
            Attempt        int             `json:"attempt"`
        }
        if err := json.Unmarshal(msg.Body, &job); err != nil {
            log.Println("unmarshal job:", err)
            msg.Nack(false, false)
            continue
        }
        w.handleJob(context.Background(), &job)
        msg.Ack(false)
    }
}

func (w *Worker) handleJob(ctx context.Context, job *struct {
    SubscriptionID string
    Payload        []byte
    Attempt        int
}) {
    sub, err := w.subRepo.GetByID(job.SubscriptionID)
    if err != nil {
        log.Println("sub fetch:", err)
        return
    }

    reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    statusCode := 0
    errMsg := ""

    req, err := http.NewRequestWithContext(reqCtx, "POST", sub.TargetURL, bytes.NewReader(job.Payload))
    if err != nil {
        errMsg = err.Error()
    } else {
        req.Header.Set("Content-Type", "application/json")
        resp, err := http.DefaultClient.Do(req)
        if err != nil {
            errMsg = err.Error()
        } else {
            statusCode = resp.StatusCode
            resp.Body.Close()
            if statusCode < 200 || statusCode >= 300 {
                errMsg = fmt.Sprintf("received status code %d", statusCode)
            }
        }
    }

    entry := &model.DeliveryLog{
        DeliveryID:     fmt.Sprintf("%s-%s", job.SubscriptionID, time.Now().Format(time.RFC3339Nano)),
        SubscriptionID: job.SubscriptionID,
        Attempt:        job.Attempt,
        Timestamp:      time.Now(),
        Status:         "Failed",
        HTTPStatus:     statusCode,
        ErrorMessage:   errMsg,
    }
    if statusCode >= 200 && statusCode < 300 {
        entry.Status = "Success"
        entry.ErrorMessage = ""
    }
    if err := w.logRepo.Insert(entry); err != nil {
        log.Println("log insert:", err)
    }

    if statusCode < 200 || statusCode >= 300 {
        if job.Attempt < 5 {
            next := job.Attempt + 1
            go func(attempt int) {
                time.Sleep(time.Duration(attempt) * 10 * time.Second)
                w.pub.Publish(job.SubscriptionID, map[string]interface{}{
                    "subscription_id": job.SubscriptionID,
                    "body":            job.Payload,
                    "attempt":         attempt,
                })
            }(next)
        }
    }
}
