package saucenao

// var json = jsoniter.ConfigCompatibleWithStandardLibrary

// type tgbot interface {
// 	GetFile(map[string]string) string
// 	SendMessage(map[string]string) error
// 	SendDocument(map[string]string, string, []byte) string

// 	GetConfig(name string) map[string]interface{}
// 	Log(interface{}, int)
// }
// type message interface {
// 	GetChatIDStr() string
// 	GetPhotoFileID() string
// 	GetMsgIDStr() string
// 	GetReplyMsgIDStr() string
// 	GetReplyToPhotoFileID() string
// }

// type saucenaoresp struct {
// 	Header struct {
// 		UserID     string `json:"user_id"`
// 		ShortLimit string `json:"short_limit"`
// 		LongLimit  string `json:"long_limit"`
// 	} `json:"header"`
// 	Results []struct {
// 		Header struct {
// 			Similarity string `json:"similarity"`
// 			Thumbnail  string `json:"thumbnail"`
// 		} `json:"header"`
// 		Data struct {
// 			ExtUrls []string `json:"ext_urls"`
// 			Title   string   `json:"title"`
// 			Author  string   `json:"author_name"`
// 			PixivID int      `json:"pixiv_id"`
// 			MemName string   `json:"member_name"`
// 		} `json:"data"`
// 	} `json:"results"`
// } //simple

// var (
// 	client = &http.Client{
// 		Timeout: 10 * time.Second,
// 	}
// 	prefix     = "https://saucenao.com/search.php?db=999&output_type=2&numres=1"
// 	ErrNoPhoto = errors.New("saucenao: No photo in this message.")
// )

// func New(b interface{}) func(interface{}, interface{}) error {
// 	conf := b.(tgbot).GetConfig("sauceNAO")
// 	prefix += "&api_key=" + conf["key"].(string) + "&url="
// 	return Handle
// }

// func Handle(m interface{}, b interface{}) error {
// 	return handle(m.(message), b.(tgbot))
// }

// func handle(msg message, bot tgbot) error {
// 	ID := msg.GetPhotoFileID()
// 	paras := map[string]string{
// 		// 1. A send Photo, B reply it with /whatpic
// 		//    -> bot will reply to A message
// 		// 2. B send photo with caption /whatpic
// 		//    -> bot will send message directly
// 		//    (reply_to_message_id == "")
// 		"reply_to_message_id": msg.GetReplyMsgIDStr(),
// 		"chat_id":             msg.GetChatIDStr(),
// 		"parse_mode":          "Markdown",
// 	}

// 	if ID == "" {
// 		ID = msg.GetReplyToPhotoFileID()
// 		if ID == "" {
// 			paras["text"] = "未找到图片，请对图片内容回复该命令。"
// 			bot.SendMessage(paras)
// 			return ErrNoPhoto
// 		}
// 	}
// 	u := bot.GetFile(map[string]string{"file_id": ID})
// 	bot.Log("sauceNAO: pic url -> "+u, 0)
// 	resp, err := client.Get(prefix + url.PathEscape(u))
// 	if err != nil {
// 		return err
// 	}

// 	var sresp saucenaoresp
// 	err = json.NewDecoder(resp.Body).Decode(&sresp)
// 	resp.Body.Close()
// 	if err != nil {
// 		return err
// 	}

// 	if sresp.Results[0].Data.PixivID != 0 {
// 		paras["text"] = "\n*Similarity :* " + sresp.Results[0].Header.Similarity
// 		paras["text"] += "%\n*Illustrator:* " + sresp.Results[0].Data.MemName +
// 			"\n*Pixiv ID     :* [" + strconv.Itoa(sresp.Results[0].Data.PixivID) + "](" + sresp.Results[0].Data.ExtUrls[0] + ")"
// 	} else {
// 		paras["text"] = sresp.Results[0].Data.ExtUrls[0]
// 		paras["text"] += "\n*Similarity:* " + sresp.Results[0].Header.Similarity
// 	}

// 	return bot.SendMessage(paras)
// }
