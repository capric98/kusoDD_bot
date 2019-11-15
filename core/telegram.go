package core

import "fmt"

func (b *tgbot) SetWebHook() error {
	resp, err := b.client.Get(b.apiUrl + "setWebhook?url=" + b.hookSuffix + b.hookPath)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (b *tgbot) CancelWebHook() error {
	resp, err := b.client.Get(b.apiUrl + "deleteWebhook")
	fmt.Println(b.apiUrl + "deleteWebhook")
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
