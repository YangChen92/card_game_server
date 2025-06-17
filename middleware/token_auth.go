package middleware

import (
	"game_server/global"
	"game_server/proto/pb"
	"game_server/utils"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"

	// "go.starlark.net/lib/proto"
	"github.com/golang/protobuf/proto"
)

type TokenAuthMiddleware struct {
	znet.BaseRouter
}

func (m *TokenAuthMiddleware) PreHandle(req ziface.IRequest) {
	// 排除登录注册接口
	if req.GetMsgID() == global.MSG_REGISTER || req.GetMsgID() == global.MSG_LOGIN {
		return
	}

	// 从请求中解析Token
	var tokenAuth pb.TokenAuth
	if err := proto.Unmarshal(req.GetData(), &tokenAuth); err != nil {
		// 返回错误响应
		return
	}

	// 验证Token
	if userID, valid := utils.VerifyToken(tokenAuth.Token); !valid {
		// 返回未授权错误
		req.Abort()
	} else {
		// 将用户ID存入请求上下文
		req.Set("userID", userID)
	}
}

func (m *TokenAuthMiddleware) RouterHandler(req ziface.IRequest) {

}
func (m *TokenAuthMiddleware) Handle(req ziface.IRequest) {}

func (m *TokenAuthMiddleware) PostHandle(req ziface.IRequest) {

}
