package main

import (
	"fmt"
	"github.com/mohamedattahri/mail"
)

func main() {
	msg := NewMessage()
	msg.SetFrom(&Address{"Al Bumin", "a.bumin@example.name"})
	msg.To().Add(&Address{"Polly Ester", "p.ester@example.com"})
	msg.SetSubject("Message with HTML, alternative text, and an attachment")
	mixed := NewMultipart("multipart/mixed", msg)

	// filename is the name that will be suggested to a user who would like to
	// download the attachment, but also the ID with which you can refer to the
	// attachment in a cid URI scheme.
	filename := "gopher.jpg"

	// The src of the image in this HTML is set to use the attachment with the
	// Content-ID filename.
	html := fmt.Sprintf("<html><body><img src=\"cid:%s\"/></body></html>", filename)
	mixed.AddText("text/html", bytes.NewReader([]byte(html)))

	// Load the photo and add the attachment with filename.
	attachment, _ := ioutil.ReadFile("path/of/image.jpg")
	mixed.AddAttachment(Attachment, filename, "image/jpeg", bytes.NewReader(attachment))

	// Closing mixed, the parent part.
	mixed.Close()

	fmt.Println(msg)
}
