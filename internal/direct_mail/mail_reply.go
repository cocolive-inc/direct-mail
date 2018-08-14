package direct_mail

// WEBAPIを使って、メールを返信する

import (
	"gitlab.com/nghinv/direct-mail/kit/log"
)

type MailReply struct {
	toAddress string
}

func NewMailReply(toAddress string) MailReply {
	return MailReply{
		toAddress: toAddress,
	}
}

func (r MailReply) Reply() {
	// TODO APIができたら、対応する
	log.GetLogger().Error("Reply Mail To " + r.toAddress + " success!" )
}

// 一瞬に、多くメールを届いたら、短い間メールを返信するできるため、
// マルチRoutineでメールを返信する

var MailReplyQueue chan MailReply

type ReplyWorker struct {
	WorkerPool chan chan MailReply
	ReplyChannel chan MailReply
	done chan bool
}

func NewReplyWorker(workerPool chan chan MailReply) ReplyWorker {
	return ReplyWorker{
		WorkerPool:workerPool,
		ReplyChannel:make(chan MailReply),
		done:make(chan bool),
	}
}

func (w ReplyWorker) Start() {
	go func() {
		for {

			w.WorkerPool <- w.ReplyChannel

			select {
			case job := <-w.ReplyChannel:
				job.Reply()

			case <-w.done:
				return
			}
		}
	}()
}

func (w ReplyWorker) Stop() {
	go func() {
		w.done <- true
	}()
}

type Dispatcher struct {
	WorkerPool chan chan MailReply
	maxWorker int
}

func NewDispatcher(maxWorker int) *Dispatcher {
	pool := make(chan chan MailReply, maxWorker)
	return &Dispatcher{WorkerPool: pool, maxWorker: maxWorker}
}

func (d *Dispatcher) Run() {

	// starting n number of workers
	for i := 0; i < d.maxWorker; i++ {
		worker := NewReplyWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <- MailReplyQueue:

			go func(job MailReply) {

				jobChannel := <-d.WorkerPool

				jobChannel <- job
			}(job)
		}
	}
}