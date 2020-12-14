package repository

import (
	"context"
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"

	"restaurant-assistant/pkg/entity"
	"restaurant-assistant/pkg/storage/database/postgres/model"
)

type UserRepository struct {
	db  *sql.DB
	ctx context.Context
}

func NewUserRepository(db *sql.DB, ctx context.Context) *UserRepository {
	return &UserRepository{
		db:  db,
		ctx: ctx,
	}
}

func (r *UserRepository) Context(ctx context.Context) *UserRepository {
	return &UserRepository{
		db:  r.db,
		ctx: ctx,
	}
}

func (r *UserRepository) Create(u entity.User) (*entity.User, error) {
	user := convertUserEntityToBaseModel(u)

	err := user.Insert(r.ctx, r.db, boil.Infer())
	if err != nil {
		return nil, err
	}

	return convertUserModelToEntity(*user)
}

func (r *UserRepository) FindByID(id string) (*entity.User, error) {
	user, err := model.UserAccounts(Where("id = ?", id)).One(r.ctx, r.db)
	if err != nil {
		return nil, err
	}

	return convertUserModelToEntity(*user)
}

func (r *UserRepository) FindByName(u string) (*entity.User, error) {
	user, err := model.UserAccounts(Where("name = ?", u)).One(r.ctx, r.db)
	if err != nil {
		return nil, err
	}

	return convertUserModelToEntity(*user)
}

func (r *UserRepository) FindByNameAndPass(u string, p string) (*entity.User, error) {
	user, err := model.UserAccounts(Where("name = ? and password = ?", u, p)).One(r.ctx, r.db)
	if err != nil {
		return nil, err
	}

	return convertUserModelToEntity(*user)
}

func convertUserEntityToBaseModel(u entity.User) *model.UserAccount {
	return &model.UserAccount{
		Name:        u.Name,
		Password:    u.Password,
		Email:       u.Email,
		UserType:    u.Type,
		Description: u.Description,
	}
}

func convertUserEntityToModel(u entity.User) *model.UserAccount {
	return &model.UserAccount{
		ID:          u.ID.String(),
		Name:        u.Name,
		Password:    u.Password,
		Email:       u.Email,
		Description: u.Description,
		UserType:    u.Type,
		CreateTime:  u.CreateTime,
		UpdateTime:  u.UpdateTime,
	}
}

func convertUserModelToEntity(u model.UserAccount) (*entity.User, error) {
	userID, err := uuid.FromString(u.ID)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:          userID,
		Name:        u.Name,
		Password:    u.Password,
		Email:       u.Email,
		Description: u.Description,
		Type:        u.UserType,
		CreateTime:  u.CreateTime,
		UpdateTime:  u.UpdateTime,
	}, nil
}

func setupUserTimes(
	ctx context.Context,
	exec boil.ContextExecutor,
	user *model.UserAccount,
) error {
	user.UpdateTime = time.Now().UTC().Round(time.Millisecond)

	if user.CreateTime.IsZero() {
		user.CreateTime = time.Now().UTC().Round(time.Millisecond)
	}

	return nil
}

//nolint:gochecknoinits
func init() {
	model.AddUserAccountHook(boil.BeforeInsertHook, setupUserTimes)
	model.AddUserAccountHook(boil.BeforeUpsertHook, setupUserTimes)
	model.AddUserAccountHook(boil.BeforeUpdateHook, setupUserTimes)
}
