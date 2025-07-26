package handler

import (
	"game_server/global"
	"game_server/model"
	"game_server/proto/pb"
	"game_server/utils"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/znet"
	"github.com/golang/protobuf/proto"
)

type UserHandler struct {
	znet.BaseRouter
}

// 注册处理
func (h *UserHandler) Register(req ziface.IRequest) {
	// 1. 频率限制
	ip := req.GetConnection().RemoteAddr().String()
	if !utils.RegisterLimit(ip) {
		sendResponse(req, 429, "注册频率过高")
		return
	}

	// 2. 解析请求
	var user pb.User
	if err := proto.Unmarshal(req.GetData(), &user); err != nil {
		sendResponse(req, 400, "请求解析错误")
		return
	}

	// 3. 密码加密
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		sendResponse(req, 500, "系统错误")
		return
	}

	// 4. 创建用户模型
	u := model.User{
		Username:      user.Username,
		Password:      hashedPwd,
		Email:         user.Email,
		Source:        user.Source,
		HeadImg:       user.HeadImg,
		Nickname:      user.Nickname,
		UserID:        user.UserId,
		Exp:           user.Exp,
		Money:         user.Money,
		DeviceID:      user.DeviceId,
		RegTime:       user.RegTime,
		LastLoginTime: user.LastLoginTime,
		LastIP:        ip,
	}

	// 5. 存储数据库
	if err := model.CreateUser(&u); err != nil {
		sendResponse(req, 400, "用户已存在")
		return
	}

	// 6. 记录日志
	global.Log.Infof("新用户注册: %s [%s]", user.Username, user.Source)

	sendResponse(req, 200, "注册成功")
}

// 登录处理
func (h *UserHandler) Login(req ziface.IRequest) {
	var user pb.User
	if err := proto.Unmarshal(req.GetData(), &user); err != nil {
		sendResponse(req, 400, "请求解析错误")
		return
	}

	// 1. 查询用户
	u, err := model.GetUserByName(user.Username)
	if err != nil {
		sendResponse(req, 404, "用户不存在")
		return
	}

	// 2. 验证密码
	if !utils.CheckPassword(user.Password, u.Password) {
		sendResponse(req, 401, "密码错误")
		return
	}

	// 3. 生成Token
	token, err := utils.GenerateToken(u.ID)
	if err != nil {
		sendResponse(req, 500, "系统错误")
		return
	}

	// 4. 记录日志
	global.Log.Infof("用户登录: %s [%s]", user.Username, user.Source)

	// 5. 返回响应
	resp := &pb.Response{
		Code:  200,
		Msg:   "登录成功",
		Token: token,
	}
	data, _ := proto.Marshal(resp)
	req.GetConnection().SendMsg(global.MSG_LOGIN_RES, data)
}

// 获取玩家信息
func (h *UserHandler) GetUserInfo(req ziface.IRequest) {
	var reqMsg pb.GetUserInfoRequest
	if err := proto.Unmarshal(req.GetData(), &reqMsg); err != nil {
		sendResponse(req, 400, "请求解析错误")
		return
	}
	// 1. 获取token
	token := reqMsg.Token
	if token == "" { // token为空
		sendResponse(req, 401, "请先登录")
		return
	}

	// 2. 查询用户
	userID, flag := utils.VerifyToken(token)
	if !flag {
		sendResponse(req, 401, "请先登录")
		return
	}
	u, err := model.GetUserData(userID)
	if err != nil {
		sendResponse(req, 404, "用户不存在")
		return
	}

	// 3. 返回响应
	resp := &pb.UserResponse{
		Response: &pb.Response{
			Code: 200,
			Msg:  "获取成功",
		},
		Data: &pb.User{
			Username:      u.Username,
			Email:         u.Email,
			Source:        u.Source,
			HeadImg:       u.HeadImg,
			Nickname:      u.Nickname,
			UserId:        u.UserID,
			Exp:           u.Exp,
			Money:         u.Money,
			DeviceId:      u.DeviceID,
			RegTime:       u.RegTime,
			LastLoginTime: u.LastLoginTime,
			LastIp:        u.LastIP,
		},
	}
	data, _ := proto.Marshal(resp)
	req.GetConnection().SendMsg(global.MSG_COMMON_RES, data)
}

func (h *UserHandler) UpdateUserInfo(req ziface.IRequest) {
	var reqMsg pb.UpdateUserInfoRequest
	if err := proto.Unmarshal(req.GetData(), &reqMsg); err != nil {
		sendResponse(req, 400, "请求解析错误")
		return
	}
	// 1. 获取token
	token := reqMsg.Token
	if token == "" { // token为空
		sendResponse(req, 401, "请先登录")
		return
	}

	// 2. 查询用户
	userID, flag := utils.VerifyToken(token)
	if !flag {
		sendResponse(req, 401, "请先登录")
		return
	}
	u, err := model.GetUserData(userID)
	if err != nil {
		sendResponse(req, 404, "用户不存在")
		return
	}

	// 3. 更新用户信息
	u.Nickname = reqMsg.Nickname
	u.HeadImg = reqMsg.HeadImg
	if err := model.UpdateUser(u); err != nil {
		sendResponse(req, 400, "更新失败")
		return
	}

	// 4. 返回响应
	resp := &pb.Response{
		Code: 200,
		Msg:  "更新成功",
	}
	data, _ := proto.Marshal(resp)
	req.GetConnection().SendMsg(global.MSG_COMMON_RES, data)
}

// 辅助函数：发送响应
func sendResponse(req ziface.IRequest, code int32, msg string) {
	resp := &pb.Response{Code: code, Msg: msg}
	data, _ := proto.Marshal(resp)
	req.GetConnection().SendMsg(global.MSG_COMMON_RES, data)
}
