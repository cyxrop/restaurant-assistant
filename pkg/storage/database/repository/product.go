package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"

	"restaurant-assistant/pkg/entity"
	"restaurant-assistant/pkg/storage/database/postgres/model"
)

type ProductRepository struct {
	db  *sql.DB
	ctx context.Context
}

func NewProductRepository(db *sql.DB, ctx context.Context) *ProductRepository {
	return &ProductRepository{
		db:  db,
		ctx: ctx,
	}
}

func (r *ProductRepository) Context(ctx context.Context) *ProductRepository {
	return &ProductRepository{
		db:  r.db,
		ctx: ctx,
	}
}

func (r *ProductRepository) Create(p entity.Product) (*entity.Product, error) {
	product := convertProductEntityToBaseModel(p)

	err := product.Insert(r.ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return convertProductModelToEntity(*product)
}

func (r *ProductRepository) FindByID(id int) (*entity.Product, error) {
	product, err := model.Products(Where("id = ?", id)).One(r.ctx, r.db)
	if err != nil {
		return nil, err
	}

	return convertProductModelToEntity(*product)
}

func (r *ProductRepository) FindByIDs(ids []int) ([]*entity.Product, error) {
	products, err := model.Products(model.ProductWhere.ID.IN(ids)).All(r.ctx, r.db)
	if err != nil {
		return nil, err
	}

	convertedProducts := make([]*entity.Product, len(products))
	for i, product := range products {
		convertedProduct, err := convertProductModelToEntity(*product)
		if err != nil {
			return nil, err
		}

		convertedProducts[i] = convertedProduct
	}

	return convertedProducts, nil
}

func (r *ProductRepository) FindByName(name string) (*entity.Product, error) {
	product, err := model.Products(Where("name = ?", name)).One(r.ctx, r.db)
	if err != nil {
		return nil, err
	}

	return convertProductModelToEntity(*product)
}

func convertProductEntityToBaseModel(p entity.Product) *model.Product {
	return &model.Product{
		Name:        p.Name,
		Description: p.Description,
	}
}

func convertProductEntityToModel(p entity.Product) *model.Product {
	return &model.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreateTime:  p.CreateTime,
		UpdateTime:  p.UpdateTime,
	}
}

func convertProductModelToEntity(p model.Product) (*entity.Product, error) {
	return &entity.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		CreateTime:  p.CreateTime,
		UpdateTime:  p.UpdateTime,
	}, nil
}

func setupProductTimes(
	ctx context.Context,
	exec boil.ContextExecutor,
	product *model.Product,
) error {
	product.UpdateTime = time.Now().UTC().Round(time.Millisecond)

	if product.CreateTime.IsZero() {
		product.CreateTime = time.Now().UTC().Round(time.Millisecond)
	}

	return nil
}

//nolint:gochecknoinits
func init() {
	model.AddProductHook(boil.BeforeInsertHook, setupProductTimes)
	model.AddProductHook(boil.BeforeUpsertHook, setupProductTimes)
	model.AddProductHook(boil.BeforeUpdateHook, setupProductTimes)
}
