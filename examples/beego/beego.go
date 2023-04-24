package main

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/Marek-maxs/gowechat"
	"github.com/Marek-maxs/gowechat/mp/message"
	"github.com/Marek-maxs/gowechat/mp/user"
	"github.com/Marek-maxs/gowechat/wxcontext"
)

var appURL = "http://localhost:8001"

//配置微信参数
var config = wxcontext.Config{
}

func hello(ctx *context.Context) {
	//微信平台mp
	var wechat = gowechat.NewWechat(config)
	mp, err := wechat.MpMgr()
	//fmt.Println(mp)
	if err != nil {
		return
	}
	user := mp.GetUser()
	isf , err := user.GetUserInfo("oiMmL6M-EbNAlNP8a9-rPPUM1wQg")
	if err != nil {
		fmt.Println(err.Error())
	}
	// 小程序的UnionID 是：oJXNh6mdBfpoP4jOwMfwGz1vfpHc
	fmt.Println("userInfo:",isf.UnionID)
//	userInfo, err := user.GetUserInfo("o3x9R6GQL-deGXtoluzGp__IuTyw")
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	fmt.Println(userInfo)
//fmt.Println(ctx.Request)
	// 传入request和responseWriter
	msgHandler := mp.GetMsgHandler(ctx.Request, ctx.ResponseWriter)
	fmt.Println(msgHandler)

	//设置接收消息的处理方法
	msgHandler.SetHandleMessageFunc(func(msg message.MixMessage) *message.Reply {
		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		return &message.Reply{message.MsgTypeText, text}
	})

	//处理消息接收以及回复
	err = msgHandler.Handle()
	if err != nil {
		fmt.Println(err)
	}
}

//wxOAuth 微信公众平台，网页授权
func wxOAuth(ctx *context.Context) {
	var wechat = gowechat.NewWechat(config)
	mp, err := wechat.MpMgr()
	if err != nil {
		return
	}

	oauthHandler := mp.GetPageOAuthHandler(ctx.Request, ctx.ResponseWriter, appURL+"/oauth")

	oauthHandler.SetFuncCheckOpenIDExisting(func(openID string) (existing bool, stopNow bool) {
		//看自己的系统中是否已经存在此openID的用户
		//如果已经存在， 调用自己的Login 方法，设置cookie等，return true
		//如果还不存在，return false, handler会自动去取用户信息
		return false, true
	})

	oauthHandler.SetFuncAfterGetUserInfo(func(user user.Info) bool {
		//已获得用户信息，这里用信息做注册使用
		//调用自己的Login方法，设置cookie等
		return false
	})

	oauthHandler.Handle()
}

func main() {
	beego.Any("/follow", hello)
	beego.Any("/oauth", wxOAuth) //需要网页授权的页面url  /oauth?target=url
	beego.Run(":8879")
}
