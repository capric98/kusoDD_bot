package getsticker

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/capric98/kusoDD_bot/core"
)

var (
	tmpDir = "tmp_getsticker"
)

// pip(3) install tgs cairosvg numpy fonttools pillow scipy opencv-python
// Windows: https://github.com/tschoonj/GTK-for-Windows-Runtime-Environment-Installer

// /usr/bin/env python3 tgsconvert.py --input-format lottie --output-format video --sanitize sticker.tgs sticker.avi
// This almost works fine, but sometimes would get an incomplete video clip, which is beyond my ability to solve it.

func decodeTGS(pic []byte, filename string, msg core.Message) {
	defer func() {
		if e := recover(); e != nil {
			msg.Bot.Printf("%6s - getsticker failed: \"%v\".\n", "info", e)
		}
	}()
	ext := "." + msg.Message.CommandArguments()
	switch ext {
	case ".webp":
	case ".gif":
	default:
		ext = ".gif"
	}
	if _, e := os.Stat(tmpDir + "/" + filename + ext); os.IsNotExist(e) {
		if e := ioutil.WriteFile(tmpDir+"/"+filename+".tgs", pic, 0777); e != nil {
			msg.Bot.Printf("%6s - getsticker failed to write tgs file: \"%v\".\n", "info", e)
			return
		}
		defer func() { _ = os.Remove(tmpDir + "/" + filename + ".tgs") }()

		tgsCMD := exec.Command(python, script, "--input-format", "lottie", "--output-format", "video",
			"--sanitize", tmpDir+"/"+filename+".tgs", tmpDir+"/"+filename+".avi")
		tgsCMD.Stderr = os.Stderr
		_ = tgsCMD.Run()

		ffmpegCMD := exec.Command("ffmpeg", "-i", tmpDir+"/"+filename+".avi",
			"-vf", "format=rgb24,geq=r='if(gt(r(X,Y)+g(X,Y)+b(X,Y),32),r(X,Y),255)':g='if(gt(r(X,Y)+g(X,Y)+b(X,Y),32),g(X,Y),255)':b='if(gt(r(X,Y)+g(X,Y)+b(X,Y),32),b(X,Y),255)'",
			"-loop", "65535", tmpDir+"/"+filename+ext)
		_ = ffmpegCMD.Run()
		// ffmpegCMD.Stderr = os.Stderr
		_ = os.Remove(tmpDir + "/" + filename + ".avi")
	}

	fr, e := os.Open(tmpDir + "/" + filename + ext)
	if e != nil {
		msg.Bot.Printf("%6s - getsticker failed to open cached gif file: \"%v\".\n", "info", e)
		return
	}
	info, _ := fr.Stat()

	resp := core.NewDocumentUpload(
		msg.Message.Chat.ID,
		core.NewFileBytes(filename+ext, fr, info.Size()),
	)
	if _, e := msg.Bot.Send(resp); e != nil {
		msg.Bot.Printf("%6s - getsticker failed to send response: \"%v\".\n", "info", e)
	}
}

func checkTmp() {
	if _, e := os.Stat(tmpDir); !os.IsNotExist(e) {
		_ = os.RemoveAll(tmpDir)
	}
	_ = os.Mkdir(tmpDir, 0660)
}
