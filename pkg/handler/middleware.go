package handler

import (
	"github.com/gin-gonic/gin"

	"restaurant-assistant/pkg/base"
	"restaurant-assistant/pkg/common"
	"restaurant-assistant/pkg/entity"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessDetails, err := common.VerifyToken(ctx.Request)
		if err != nil {
			base.SendError(ctx, err)
			ctx.Abort()
			return
		}

		ctx.Set(entity.ContextTokenIDKey, accessDetails.TokenID)
		ctx.Set(entity.ContextTokenTypeKey, accessDetails.Type)
		ctx.Set(entity.ContextPairedTokenIdKey, accessDetails.PairedTokenID)

		ctx.Next()
	}
}
