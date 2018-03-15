package cousume

import (
	"time"
	"strings"
	"os"
	"github.com/Shopify/sarama"
	"log-shiper/tool"
)


func WriteToKakc(c chan sarama.ProducerMessage, brokers string){
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
    config.Producer.Return.Successes = true
    config.Producer.Timeout = 5 * time.Second
    p, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
    if err != nil {
        tool.Logger.Error(err.Error())
        os.Exit(-1)
	}
	defer func() {
		if err := p.Close(); err != nil {
            tool.Logger.Error(err.Error())
            os.Exit(-1)
		}
	}()
	for {
		if msg, ok := <- c; ok{
			partition, offset, err := p.SendMessage(&msg)
			if err != nil {
				tool.Logger.Error(err.Error())
				os.Exit(-1)
			}
			tool.Logger.Info("Message is stored in /partition(%d)/offset(%d)\n", partition, offset)
		}else{
			tool.Logger.Error("read data form channel error.")
			os.Exit(-1)
		}
	}
}