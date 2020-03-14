package getsticker

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/capric98/kusoDD_bot/core"
)

// pip(3) install tgs cairosvg numpy fonttools pillow scipy opencv-python
// Windows: https://github.com/tschoonj/GTK-for-Windows-Runtime-Environment-Installer

// /usr/bin/env python3 tgsconvert.py --input-format lottie --output-format video --sanitize sticker.tgs sticker.avi

func decodeTGS(pic []byte, filename string, msg core.Message) {
	defer func() {
		if e := recover(); e != nil {
			msg.Bot.Printf("%6s - getsticker failed: \"%v\".\n", "info", e)
		}
	}()
	if _, e := os.Stat("tmp/" + filename + ".gif"); os.IsNotExist(e) {
		if e := ioutil.WriteFile("tmp/"+filename+".tgs", pic, 0777); e != nil {
			msg.Bot.Printf("%6s - getsticker failed to write tgs file: \"%v\".\n", "info", e)
			return
		}
		defer func() { _ = os.Remove("tmp/" + filename + ".tgs") }()

		tgsCMD := exec.Command(python, script, "--input-format", "lottie", "--output-format", "video",
			"--sanitize", "tmp/"+filename+".tgs", "tmp/"+filename+".avi")
		tgsCMD.Stderr = os.Stderr
		_ = tgsCMD.Run()

		ffmpegCMD := exec.Command("ffmpeg", "-i", "tmp/"+filename+".avi", "tmp/"+filename+".gif")
		_ = ffmpegCMD.Run()
		_ = os.Remove("tmp/" + filename + ".avi")
	}

	fr, e := os.Open("tmp/" + filename + ".gif")
	if e != nil {
		msg.Bot.Printf("%6s - getsticker failed to open cached gif file: \"%v\".\n", "info", e)
		return
	}
	info, _ := fr.Stat()

	resp := core.NewDocumentUpload(
		msg.Message.Chat.ID,
		core.NewFileBytes(filename+".gif", fr, info.Size()),
	)
	if _, e := msg.Bot.Send(resp); e != nil {
		msg.Bot.Printf("%6s - getsticker failed to send response: \"%v\".\n", "info", e)
	}
}

func checkTmp() {
	if _, e := os.Stat("tmp"); os.IsNotExist(e) {
		_ = os.Mkdir("tmp", 0660)
	}
}
