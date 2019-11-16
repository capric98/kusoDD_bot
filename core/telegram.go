package core

import (
	"net/url"
	"strconv"
)

type Tgbot interface {
	SetWebHook() error
	CancelWebHook() error
	SendChatAction(interface{}, interface{}) error
	SendText(interface{}, string, bool) error
}

func (b *tgbot) SetWebHook() error {
	resp, err := b.client.Get(b.apiUrl + "setWebhook?url=" + b.hookSuffix + b.hookPath)
	b.Log("Set webhook.", 0)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (b *tgbot) CancelWebHook() error {
	resp, err := b.client.Get(b.apiUrl + "deleteWebhook")
	b.Log("Delete webhook.", 0)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (b *tgbot) SendChatAction(msg *Message, action string) error {
	resp, e := b.client.Get(b.apiUrl + "sendChatAction?chat_id=" + toStr(msg.Message.Chat.ID) + "&action=" + action)
	if e != nil {
		return e
	}
	resp.Body.Close()
	return nil
}

func (b *tgbot) SendText(m interface{}, text string, reply bool) error {
	msg := m.(*Message)
	furl := b.apiUrl + "sendmessage?chat_id=" + toStr(msg.Message.Chat.ID) + "&text=" + url.QueryEscape(text)
	if reply {
		furl += "&reply_to_message_id=" + toStr(msg.Message.MessageID)
	}

	resp, e := b.client.Get(furl)
	//b.Log(furl, 0)
	if e != nil {
		return e
	}
	resp.Body.Close()
	return nil
}

func toStr(n int64) string {
	return strconv.FormatInt(n, 10)
}
