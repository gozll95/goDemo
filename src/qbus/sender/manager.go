package send

import (
	"github.com/dolab/gogo"
	"github.com/qbusapi/app/models"

	"time"
)

type MessageManager struct {
	Queue    *MessageQueue
	Sender   Sender
	SendChan chan *models.MessageModel

	Config *Config
	Logger gogo.Logger
}

func NewMessageManager(sender Sender, config *Config, logger gogo.Logger) *MessageManager {
	manager := &MessageManager{
		Queue: &MessageQueue{
			First: nil,
			Tail:  nil,
		},
		Sender:   sender,
		SendChan: make(chan *models.MessageModel, config.Ratelimit+1),
		Config:   config,
		Logger:   logger,
	}

	go manager.dispatch()

	for i := 0; i < config.Concurrent; i++ {
		go manager.sending()
	}

	return manager
}

func (manager *MessageManager) SendMessage(message *models.MessageModel) (err error) {
	err = message.Save()
	if err != nil {
		manager.Logger.Errorf("%#v message.Save() %v", err)
		return
	}

	manager.Queue.AddMessage(message)

	return
}

func (manager *MessageManager) dispatch() {
	ticker := time.Tick(time.Second * 60)
	for now := range ticker {
		manager.Logger.Infof("dispatching... at %v", now)
		for {
			message := manager.Queue.FetchMessage()
			if message == nil {
				break
			}

			manager.SendChan <- message

			if len(manager.SendChan) > manager.Config.Ratelimit {
				break
			}
		}

		manager.Logger.Info("dispatched... at")
	}
}

// TODO retry when send mail failed
func (manager *MessageManager) sending() {
	for {
		message := <-manager.SendChan

		err := manager.Sender.Send(message)
		manager.Logger.Infof("%#v, message sended", message)
		if err != nil {
			manager.Logger.Errorf("manager.Sender %v Send(%#v) %v", manager.Config.Send.Provider, message, err)

			message.Status = models.MessageStatusFailed
			err = message.Save()
			if err != nil {
				manager.Logger.Errorf("%#v message.Save() %v", message, err)
			}

			continue
		}

		message.Status = models.MessageStatusSuccess
		err = message.Save()
		if err != nil {
			manager.Logger.Errorf("%#v message.Save() %v", message, err)
		}

	}
}
