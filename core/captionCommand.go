package core

const (
	NAME = "kusoDD_bot"
)

func (m Message) CaptionCommand() string {
	// return the first command in message caption.
	for k := range m.Message.CaptionEntities {
		ce := m.Message.CaptionEntities[k]
		if ce.Type == "bot_command" {
			p := ce.Offset + 1
			for ; p < ce.Offset+ce.Length && m.Message.Caption[p] != '@'; p++ {
			}
			if p == ce.Offset+ce.Length {
				return m.Message.Caption[ce.Offset+1 : p]
			} else {
				if m.Message.Caption[p+1:ce.Offset+ce.Length] == NAME {
					return m.Message.Caption[ce.Offset+1 : p]
				}
			}
		}
	}
	return ""
}
