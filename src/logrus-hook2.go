import (
    "github.com/vladoatanasov/logrus_amqp"
    "gopkg.in/olivere/elastic.v5"
    "gopkg.in/sohlich/elogrus.v2"
    log "github.com/sirupsen/logrus"
    "github.com/pkg/errors"
)

// config logrus log to amqp
func ConfigAmqpLogger(server, username, password, exchange, exchangeType, virtualHost, routingKey string) {
    hook := logrus_amqp.NewAMQPHookWithType(server, username, password, exchange, exchangeType, virtualHost, routingKey)
    log.AddHook(hook)
}

// config logrus log to es
func ConfigESLogger(esUrl string, esHOst string, index string) {
    client, err := elastic.NewClient(elastic.SetURL(esUrl))
    if err != nil {
        log.Errorf("config es logger error. %+v", errors.WithStack(err))
    }
    esHook, err := elogrus.NewElasticHook(client, esHOst, log.DebugLevel, index)
    if err != nil {
        log.Errorf("config es logger error. %+v", errors.WithStack(err))
    }
    log.AddHook(esHook)
}