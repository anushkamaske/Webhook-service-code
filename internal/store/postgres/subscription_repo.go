package postgres

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "webhook-service/internal/model"
)

type SubscriptionRepo struct {
    db *gorm.DB
}

func NewPostgres(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`)
    db.AutoMigrate(&model.Subscription{})
    return db, nil
}

func NewSubscriptionRepo(db *gorm.DB) *SubscriptionRepo {
    return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Create(sub *model.Subscription) error {
    return r.db.Create(sub).Error
}

func (r *SubscriptionRepo) GetByID(id string) (*model.Subscription, error) {
    var sub model.Subscription
    if err := r.db.First(&sub, "id = ?", id).Error; err != nil {
        return nil, err
    }
    return &sub, nil
}

func (r *SubscriptionRepo) List() ([]model.Subscription, error) {
    var subs []model.Subscription
    err := r.db.Find(&subs).Error
    return subs, err
}

func (r *SubscriptionRepo) Update(sub *model.Subscription) error {
    return r.db.Save(sub).Error
}

func (r *SubscriptionRepo) Delete(id string) error {
    return r.db.Delete(&model.Subscription{}, "id = ?", id).Error
}
