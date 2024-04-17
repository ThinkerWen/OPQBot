package main

import (
	"context"
	"encoding/base64"
	"os"
	"time"

	"github.com/charmbracelet/log"
	_ "github.com/mcoo/OPQBot/session/provider"
	"github.com/opq-osc/OPQBot/v2"
	"github.com/opq-osc/OPQBot/v2/apiBuilder"
	"github.com/opq-osc/OPQBot/v2/events"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	core, err := OPQBot.NewCore(apiUrl, OPQBot.WithMaxRetryCount(5), OPQBot.WithAutoSignToken(123123, 123123))
	if err != nil {
		log.Fatal(err)
	}
	// 群组相关功能
	go HandleGroupMsg(core)
	core.On(events.EventNameGroupMsg, func(ctx context.Context, event events.IEvent) {
		if event.GetMsgType() == events.MsgTypeGroupMsg {
			groupMsg := event.ParseGroupMsg()
			if groupMsg.AtBot() {
				text := groupMsg.ExcludeAtInfo().ParseTextMsg().GetTextContent()
				if text == "test" {
					apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg("test").DoWithCallBack(ctx, func(iApiBuilder *apiBuilder.Response, err error) {
						response, err := iApiBuilder.GetGroupMessageResponse()
						if err != nil {
							return
						}
						time.Sleep(time.Second * 1)
						apiBuilder.New(apiUrl, event.GetCurrentQQ()).GroupManager().RevokeMsg().ToGUin(groupMsg.GetGroupUin()).MsgSeq(response.MsgSeq).MsgRandom(response.MsgTime).Do(ctx)
					})
				}
				if text == "reply" {
					apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).TextMsg("回复").ReplyMsg(groupMsg.GetMsgSeq(), groupMsg.GetMsgUid()).Do(ctx)
				}
				if text == "pic" {
					p, err := os.ReadFile("example/test.jpg")
					if err != nil {
						panic(err)
					}

					pic, err := apiBuilder.New(apiUrl, event.GetCurrentQQ()).Upload().GroupPic().SetBase64Buf(base64.StdEncoding.EncodeToString(p)).DoUpload(ctx)
					if err != nil {
						panic(err)
					}
					log.Debug(pic)
					apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).PicMsg(pic).Do(ctx)
				}
				if text == "voice" {
					v, err := os.ReadFile("example/test.amr")
					if err != nil {
						panic(err)
					}

					voice, err := apiBuilder.New(apiUrl, event.GetCurrentQQ()).Upload().GroupVoice().SetBase64Buf(base64.StdEncoding.EncodeToString(v)).DoUpload(ctx)
					if err != nil {
						panic(err)
					}
					log.Debug(voice)
					apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().GroupMsg().ToUin(groupMsg.GetGroupUin()).VoiceMsg(voice).Do(ctx)
				}
			}
		}

	})
	core.On(events.EventNameFriendMsg, func(ctx context.Context, event events.IEvent) {
		if event.GetMsgType() == events.MsgTypeFriendsMsg {
			friendMsg := event.ParseFriendMsg()
			if friendMsg.ParseTextMsg().GetTextContent() == "test" && friendMsg.GetSenderUin() != event.GetCurrentQQ() {
				apiBuilder.New(apiUrl, event.GetCurrentQQ()).SendMsg().FriendMsg().ToUin(friendMsg.GetFriendUin()).TextMsg("test").Do(ctx)
			}
		}
	})
	err = core.ListenAndWait(context.Background())
	if err != nil {
		panic(err)
	}
}
