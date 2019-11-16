package core

import "net/http"

// var msgPool = sync.Pool{
// 	New: func() interface{} {
// 		return new(Message)
// 	},
// }

func (b *tgbot) Handle(r *http.Request) {
	defer func() { recover() }()
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
	for i := 0; i < b.plugnum; i++ {
		go b.plugins[i].Handle(pmsg, b)
	}
}
