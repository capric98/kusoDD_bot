package core

import "net/http"

// var msgPool = sync.Pool{
// 	New: func() interface{} {
// 		return new(Message)
// 	},
// }

func (b *tgbot) Handle(r *http.Request) {
	//msg := msgPool.Get()
	//pmsg := msg.(*Message)
	//We have to cleanup msg before putting it back to the pool,
	//which is not worthwhile...
	pmsg := &Message{}

	//defer func() {
	//	CleanPut(pmsg)
	//	msgPool.Put(msg)
	//}()

	if e := json.NewDecoder(r.Body).Decode(pmsg); e != nil {
		b.Log(e, 1)
	}

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
