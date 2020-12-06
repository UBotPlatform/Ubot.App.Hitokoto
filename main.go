package main

import (
	"encoding/json"
	"net/http"
	"strings"

	ubot "github.com/UBotPlatform/UBot.Common.Go"
)

var api *ubot.AppApi
var typeMappings = [...]struct {
	Name         string
	InternalType string
}{
	{"动画", "a"},
	{"动漫", "a"},
	{"漫画", "b"},
	{"游戏", "c"},
	{"文学", "d"},
	{"原创", "e"},
	{"网络", "f"},
	{"其他", "g"},
	{"其她", "g"},
	{"其它", "g"},
	{"影视", "h"},
	{"诗词", "i"},
	{"网易", "j"},
	{"云村", "j"},
	{"哲学", "k"},
	{"机灵", "l"},
}

type HitokotoResponse struct {
	ID         int    `json:"id"`
	UUID       string `json:"uuid"`
	Hitokoto   string `json:"hitokoto"`
	Type       string `json:"type"`
	From       string `json:"from"`
	FromWho    string `json:"from_who"`
	Creator    string `json:"creator"`
	CreatorUID int    `json:"creator_uid"`
	Reviewer   int    `json:"reviewer"`
	CommitFrom string `json:"commit_from"`
	CreatedAt  string `json:"created_at"`
	Length     int    `json:"length"`
}

func onReceiveChatMessage(bot string, msgType ubot.MsgType, source string, sender string, message string, info ubot.MsgInfo) (ubot.EventResultType, error) {
	if strings.Contains(message, "一言") {
		var url string
		url = "https://v1.hitokoto.cn/?encode=json"
		for _, typeMapping := range typeMappings {
			if strings.Contains(message, typeMapping.Name) {
				url += "&c=" + typeMapping.InternalType
			}
		}
		for {
			var builder ubot.MsgBuilder
			resp, err := http.Get(url)
			if err != nil {
				break
			}
			defer resp.Body.Close()
			var result HitokotoResponse
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				break
			}
			builder.WriteString("『")
			builder.WriteString(result.Hitokoto)
			builder.WriteString("』\n——")
			if result.FromWho != "" && result.FromWho != result.From {
				builder.WriteString(result.FromWho)
			}
			builder.WriteString("『")
			builder.WriteString(result.From)
			builder.WriteString("』")
			_ = api.SendChatMessage(bot, msgType, source, sender, builder.String())
			break //nolint
		}
		return ubot.CompleteEvent, nil
	}
	return ubot.IgnoreEvent, nil
}

func main() {
	err := ubot.HostApp("Hitokoto", func(e *ubot.AppApi) *ubot.App {
		api = e
		return &ubot.App{
			OnReceiveChatMessage: onReceiveChatMessage,
		}
	})
	ubot.AssertNoError(err)
}
