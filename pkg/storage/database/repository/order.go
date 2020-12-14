package repository

import (
	"context"
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"
	. "github.com/volatiletech/sqlboiler/types"

	"restaurant-assistant/pkg/entity"
	"restaurant-assistant/pkg/storage/database/postgres/model"
)

type OrderRepository struct {
	db  *sql.DB
	ctx context.Context
}

func NewOrderRepository(db *sql.DB, ctx context.Context) *OrderRepository {
	return &OrderRepository{
		db:  db,
		ctx: ctx,
	}
}

func (r *OrderRepository) Context(ctx context.Context) *OrderRepository {
	return &OrderRepository{
		db:  r.db,
		ctx: ctx,
	}
}

func (r *OrderRepository) Create(o entity.Order) (*entity.Order, error) {
	order, err := convertOrderEntityToBaseModel(o)
	if err != nil {
		return nil, err
	}

	if err := order.Insert(r.ctx, r.db, boil.Infer()); err != nil {
		return nil, err
	}

	return convertOrderModelToEntity(*order)
}

func (r *OrderRepository) FindByID(id int) (*entity.Order, error) {
	order, err := model.CustomerOrders(Where("id = ?", id)).One(r.ctx, r.db)
	if err != nil {
		return nil, err
	}

	return convertOrderModelToEntity(*order)
}

func convertOrderEntityToBaseModel(o entity.Order) (*model.CustomerOrder, error) {
	var orderItems JSON
	err := orderItems.Marshal(o.Items)
	if err != nil {
		return nil, err
	}

	return &model.CustomerOrder{
		OrderData: orderItems,
		UserID:    o.UserID.String(),
	}, nil
}

func convertOrderEntityToModel(o entity.Order) (*model.CustomerOrder, error) {
	var orderItems JSON
	err := orderItems.Marshal(o.Items)
	if err != nil {
		return nil, err
	}

	return &model.CustomerOrder{
		ID:         o.ID,
		OrderData:  orderItems,
		UserID:     o.UserID.String(),
		CreateTime: o.CreateTime,
		UpdateTime: o.UpdateTime,
	}, nil
}

func convertOrderModelToEntity(o model.CustomerOrder) (*entity.Order, error) {
	userID, err := uuid.FromString(o.UserID)
	if err != nil {
		return nil, err
	}

	var orderItems []entity.OrderItem
	if err = o.OrderData.Unmarshal(&orderItems); err != nil {
		return nil, err
	}

	return &entity.Order{
		ID:         o.ID,
		Items:      orderItems,
		UserID:     userID,
		CreateTime: o.CreateTime,
		UpdateTime: o.UpdateTime,
	}, nil
}

func setupOrderTimes(
	ctx context.Context,
	exec boil.ContextExecutor,
	order *model.CustomerOrder,
) error {
	order.UpdateTime = time.Now().UTC().Round(time.Millisecond)

	if order.CreateTime.IsZero() {
		order.CreateTime = time.Now().UTC().Round(time.Millisecond)
	}

	return nil
}

//nolint:gochecknoinits
func init() {
	model.AddCustomerOrderHook(boil.BeforeInsertHook, setupOrderTimes)
	model.AddCustomerOrderHook(boil.BeforeUpsertHook, setupOrderTimes)
	model.AddCustomerOrderHook(boil.BeforeUpdateHook, setupOrderTimes)
}
