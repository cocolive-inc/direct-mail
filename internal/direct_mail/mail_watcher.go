package direct_mail

// メールのディレクトリを監視して、新しいメールが入ったら
// WEBAPIを使って、メールを返信する

import (
	"github.com/fsnotify/fsnotify"
	"fmt"
	"gitlab.com/nghinv/direct-mail/kit/cfg"
	"gitlab.com/nghinv/direct-mail/kit/log"
)

type MailWatcher struct {
	parser *MailParser

	watcher *fsnotify.Watcher
}

func NewMailWatcher() (*MailWatcher,error) {

	parser := NewMailParser()
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return nil, fmt.Errorf("[(MailWatcher)NewMailWatcher] Create file watcher error: %v", err)
	}

	return &MailWatcher{
		parser:parser,
		watcher:watcher,
	},nil
}

func (w *MailWatcher) Watch() error {

	defer w.watcher.Close()

	if err := w.watcher.Add(cfg.GetString("mail_store_path")); err != nil {
		return fmt.Errorf("[(MailWatcher)Watch] Add watch directory error: %v", err)
	}

	done := make(chan bool)

	// 一瞬に最大1000メール件を届くのが想定される
	MailReplyQueue = make(chan MailReply, 1000)

	// 同時に5スレッドを使って、メールを返信する
	replyDispatcher := NewDispatcher(5)
	replyDispatcher.Run()

	go func() {
		for {
			select {
			// watch for events
			case event := <- w.watcher.Events:

				if event.Op.String() == "CREATE" {

					mail, err := w.parser.Parser(event.Name)
					if err != nil {
						log.GetLogger().Error(fmt.Sprintf("[(MailWatcher) Watch] Error: %v", err))
					}

					reply := NewMailReply(mail)
					MailReplyQueue <- reply
				}
			}
		}
	}()

	<-done
	return nil
}