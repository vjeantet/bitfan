package watch

import (
	"log"
	"os"

	"github.com/etrepat/postman/handler"
	"github.com/etrepat/postman/imap"
	"github.com/etrepat/postman/version"
)

const (
	DELIVERY_MODE_POSTBACK = "postback"
	DELIVERY_MODE_LOGGER   = "logger"
)

var (
	DefaultLogger  = log.New(os.Stdout, "[watch] ", log.LstdFlags)
	DELIVERY_MODES = map[string]bool{
		DELIVERY_MODE_POSTBACK: true,
		DELIVERY_MODE_LOGGER:   true}
)

type Flags struct {
	Host          string
	Port          uint
	Ssl           bool
	Username      string
	Password      string
	Mailbox       string
	Mode          string
	PostbackUrl   string
	PostEncoded   bool
	PostParamName string
}

type Watch struct {
	mailbox  string
	handlers []handler.MessageHandler
	client   *imap.ImapClient
	logger   *log.Logger
	chMsgs   chan []string
}

func (w *Watch) Mailbox() string {
	return w.mailbox
}

func (w *Watch) SetMailbox(value string) {
	w.mailbox = value
}

func (w *Watch) SetLogger(logger *log.Logger) {
	w.logger = logger
}

func (w *Watch) Logger() *log.Logger {
	return w.logger
}

func (w *Watch) AddHandler(handler handler.MessageHandler) {
	w.handlers = append(w.handlers, handler)
}

func (w *Watch) Handlers() []handler.MessageHandler {
	return w.handlers
}

func (w *Watch) Start() {
	w.logger.Println("Starting ", version.VersionShort())

	w.chMsgs = make(chan []string)

	go w.handleIncoming()

	w.logger.Printf("Handling incoming messages with:")
	for i := 0; i < len(w.handlers); i++ {
		w.logger.Printf("> %s", w.handlers[i].Describe())
	}

	err := w.monitorMailbox()
	if err != nil {
		w.logger.Fatalln(err)
	}
}

func (w *Watch) Stop() {
	// TODO: Unimplemented for now
}

func (w *Watch) handleIncoming() {
	var err error

	for {
		messages := <-w.chMsgs

		for _, msg := range messages {
			for _, handler := range w.handlers {
				err = handler.Deliver(msg)
				if err != nil {
					w.logger.Println(err)
				} else {
					w.logger.Println("Delivered successfully")
				}
			}
		}
	}
}

func (w *Watch) monitorMailbox() error {
	var messages []string
	var err error

	w.logger.Printf("Initiating connection to %s", w.client.Addr())
	err = w.client.Connect()
	if err != nil {
		return err
	}

	defer w.client.Disconnect()

	w.logger.Printf("Switching to %s", w.mailbox)
	err = w.client.Select(w.mailbox)
	if err != nil {
		return err
	}

	w.logger.Printf("Checking for new (unseen) messages")
	messages, err = w.client.Unseen()
	if err != nil {
		return err
	}

	if len(messages) != 0 {
		w.logger.Printf("Detected %d new (unseen) messages. Delivering...", len(messages))
		w.chMsgs <- messages
	}

	for {
		w.logger.Printf("Waiting for new messages")
		messages, err = w.client.Incoming()
		if err != nil {
			return err
		}

		if len(messages) != 0 {
			w.logger.Printf("Detected %d new (unseen) messages. Delivering...", len(messages))
			w.chMsgs <- messages
		}
	}

	return nil
}

func NewFlags() *Flags {
	return &Flags{}
}

func New(flags *Flags, handlers ...handler.MessageHandler) *Watch {
	watch := &Watch{
		mailbox: flags.Mailbox,
		client:  imap.NewClient(flags.Host, flags.Port, flags.Ssl, flags.Username, flags.Password),
		logger:  DefaultLogger}

	if len(handlers) != 0 {
		for _, hnd := range handlers {
			watch.AddHandler(hnd)
		}
	} else {
		switch flags.Mode {
		case DELIVERY_MODE_POSTBACK:
			watch.AddHandler(handler.New(handler.POSTBACK_HANDLER, flags.PostbackUrl, flags.PostEncoded, flags.PostParamName))

		case DELIVERY_MODE_LOGGER:
			watch.AddHandler(handler.New(handler.LOGGER_HANDLER, DefaultLogger))
		}
	}

	return watch
}

func DeliveryModeValid(mode string) bool {
	return DELIVERY_MODES[mode]
}

func ValidDeliveryModes() []string {
	modes := make([]string, len(DELIVERY_MODES))
	i := 0

	for k, _ := range DELIVERY_MODES {
		modes[i] = k
		i++
	}

	return modes
}
