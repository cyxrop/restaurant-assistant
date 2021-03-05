package handler

import (
	"github.com/gin-gonic/gin"

	"restaurant-assistant/pkg/base"
	"restaurant-assistant/pkg/entity"
	"restaurant-assistant/pkg/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessDetails, err := utils.VerifyToken(ctx.Request)
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
