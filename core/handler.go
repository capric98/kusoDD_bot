package core

import "net/http"

import "sync"

var msgPool = sync.Pool{
	New: func() interface{} {
		return new(Message)
	},
}

func (b *tgbot) Handle(r *http.Request) {
	msg := msgPool.Get()
	defer msgPool.Put(msg)
	if e := json.NewDecoder(r.Body).Decode(msg); e != nil {
		b.Log(e, 1)
	}

	pmsg := msg.(*Message)
	b.Log(pmsg.GetMsgLog(), 0)
	for i := 0; i < b.plugnum; i++ {
		done, err := b.plugins[i].Handle(pmsg)
		if err != nil {
			b.Log(err, 1)
		} else {
			if done {
				b.Log("Plugin: "+b.plugins[i].Name()+" handled the message "+pmsg.GetStrMsgID(), 0)
				//return
			}
		}
	}
}
