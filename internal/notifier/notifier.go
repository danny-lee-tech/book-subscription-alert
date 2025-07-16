package notifier

import "github.com/xconstruct/go-pushbullet"

type PushBulleter struct {
	APIKey string
	Tag    string
	Title  string
}

type Notifier struct {
	PushBullet *PushBulleter
}

func (n *Notifier) Notify(messageString string) error {
	return n.PushBullet.postNoteToChannel(messageString)
}

func (n *Notifier) NotifyWithLink(messageString string, url string) error {
	return n.PushBullet.postLinkToChannel(messageString, url)
}

func (p *PushBulleter) postNoteToChannel(messageString string) error {
	pb := pushbullet.New(p.APIKey)

	err := pb.PushNoteToChannel(p.Tag, p.Title, messageString)
	if err != nil {
		return err
	}

	return nil
}

func (p *PushBulleter) postLinkToChannel(messageString string, url string) error {
	pb := pushbullet.New(p.APIKey)

	err := pb.PushLinkToChannel(p.Tag, p.Title, url, messageString)
	if err != nil {
		return err
	}

	return nil
}
