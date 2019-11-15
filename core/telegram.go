package core

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
