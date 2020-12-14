package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"restaurant-assistant/pkg/app/service"
	"restaurant-assistant/pkg/base"
	"restaurant-assistant/pkg/common"
	"restaurant-assistant/pkg/entity"
)

type RestaurantAssistantServer struct {
	log *zap.Logger // todo: create logger interface
	as  *service.AuthenticationService
	us  *service.UserService
	ps  *service.ProductService
	os  *service.OrderService
}

func NewRestaurantAssistantServer(
	log *zap.Logger,
	as *service.AuthenticationService,
	us *service.UserService,
	ps *service.ProductService,
	os *service.OrderService,
) (*RestaurantAssistantServer, error) {
	if log == nil {
		return nil, errors.New("logger is not defined")
	}
	if as == nil {
		return nil, errors.New("authentication service is not defined")
	}
	if us == nil {
		return nil, errors.New("user service is not defined")
	}
	if ps == nil {
		return nil, errors.New("product service is not defined")
	}
	if os == nil {
		return nil, errors.New("order service is not defined")
	}

	return &RestaurantAssistantServer{
		log: log,
		as:  as,
		us:  us,
		ps:  ps,
		os:  os,
	}, nil
}

func (ras *RestaurantAssistantServer) Login(ctx *gin.Context) {
	var u entity.User
	if err := ctx.ShouldBindJSON(&u); err != nil {
		base.SendNewError(ctx, base.ErrInvalidParameters, "invalid json provided")
		return
	}

	tokenPair, err := ras.as.Login(ctx, u.Name, u.Password)
	if err != nil {
		base.SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, common.FormatResponse(tokenPair))
}

func (ras *RestaurantAssistantServer) RefreshToken(ctx *gin.Context) {
	tokenPair, err := ras.as.RefreshToken(ctx)
	if err != nil {
		base.SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, common.FormatResponse(tokenPair))
}

func (ras *RestaurantAssistantServer) Logout(ctx *gin.Context) {
	err := ras.as.Logout(ctx)
	if err != nil {
		base.SendError(ctx, err)
		return
	}

	response := entity.NewBaseResponse("successfully logged out")
	ctx.JSON(http.StatusOK, common.FormatResponse(response))
}

func (ras *RestaurantAssistantServer) CreateUser(ctx *gin.Context) {
	var u entity.User
	if err := ctx.ShouldBindJSON(&u); err != nil {
		base.SendNewError(ctx, base.ErrInvalidParameters, "invalid json provided")
		return
	}

	user, err := ras.us.CreateUser(ctx, u)
	if err != nil {
		base.SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, common.FormatResponse(user))
}

func (ras *RestaurantAssistantServer) UpdateUser(ctx *gin.Context) {
	var u entity.User
	if err := ctx.ShouldBindJSON(&u); err != nil {
		base.SendNewError(ctx, base.ErrInvalidParameters, "invalid json provided")
		return
	}

	if base.IsEmptyUUID(u.ID) {
		base.SendNewError(ctx, base.ErrInvalidParameters, "user id required")
		return
	}

	clientUser, err := ras.as.GetUserFromContext(ctx)
	if err != nil {
		base.SendError(ctx, err)
		return
	}

	if (clientUser.ID.String() != u.ID.String()) && !clientUser.IsAdmin() {
		base.SendNewError(ctx, base.ErrNoPermissions, "only the user can change himself or the admin")
		return
	}

	fmt.Printf("user: %+v", &u)

	ctx.JSON(http.StatusOK, common.FormatResponse(&u))
}

func (ras *RestaurantAssistantServer) CreateProduct(ctx *gin.Context) {
	user, err := ras.as.GetUserFromContext(ctx)
	if err != nil {
		base.SendError(ctx, err)
		return
	}

	if !user.IsAdmin() {
		base.SendNewError(ctx, base.ErrNoPermissions, "no permission to this action")
		return
	}

	var product entity.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		base.SendNewError(ctx, base.ErrInvalidParameters, "invalid json provided")
		return
	}

	newProduct, err := ras.ps.CreateProduct(ctx, product)
	if err != nil {
		base.SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, common.FormatResponse(newProduct))
}

func (ras *RestaurantAssistantServer) CreateOrder(ctx *gin.Context) {
	user, err := ras.as.GetUserFromContext(ctx)
	if err != nil {
		base.SendError(ctx, err)
		return
	}

	var order entity.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		base.SendNewError(ctx, base.ErrInvalidParameters, "invalid json provided")
		return
	}

	order.UserID = user.ID

	newOrder, err := ras.os.CreateOrder(ctx, order)
	if err != nil {
		base.SendError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, common.FormatResponse(newOrder))

	ctx.Next()
}
