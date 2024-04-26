package natsclient

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log/slog"
	"time"
)

const (
	AddBalance = "add_balance"
	Pays       = "pays"
	Notify     = "notify"
)

type Client struct {
	conn *nats.Conn
	Js   jetstream.JetStream
}

func NewClient(uri string) *Client {
	nc, _ := nats.Connect(uri)
	c := &Client{conn: nc}
	c.migrate()
	return c
}

func (c *Client) migrate() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()
	stream, err := jetstream.New(c.conn)
	if err != nil {
		slog.Error("jetStream new", "err", err)
		panic(err)
	}
	c.Js = stream
	checkSteam(ctx, stream, AddBalance, Pays, Notify)
	fmt.Println("---")

}

func checkSteam(ctx context.Context, js jetstream.JetStream, streamNames ...string) {
	for _, streamName := range streamNames {
		_, err := js.Stream(ctx, streamName)

		if err != nil {
			if errors.As(err, &jetstream.ErrStreamNotFound) {
				streamConfig := jetstream.StreamConfig{
					Name:                 streamName,
					Description:          "",
					Subjects:             []string{streamName},
					Retention:            jetstream.WorkQueuePolicy,
					MaxConsumers:         -1,
					MaxMsgs:              -1,
					MaxBytes:             -1,
					Discard:              jetstream.DiscardNew,
					DiscardNewPerSubject: false,
					MaxAge:               0,
					MaxMsgsPerSubject:    0,
					MaxMsgSize:           0,
					Storage:              jetstream.MemoryStorage,
					Replicas:             1,
					NoAck:                false,
					Duplicates:           0,
					Placement:            nil,
					Mirror:               nil,
					Sources:              nil,
					Sealed:               false,
					DenyDelete:           false,
					DenyPurge:            false,
					AllowRollup:          false,
					Compression:          0,
					FirstSeq:             0,
					SubjectTransform:     nil,
					RePublish:            nil,
					AllowDirect:          false,
					MirrorDirect:         false,
					ConsumerLimits:       jetstream.StreamConsumerLimits{},
					Metadata:             nil,
					Template:             "",
				}
				createStream, err := js.CreateStream(ctx, streamConfig)
				if err != nil {
					slog.Error("create stream", "err", err)
					panic(err)
				}
				slog.Info("NATS check, stream created", "stream", streamName)

				_, err = createStream.CreateConsumer(ctx, jetstream.ConsumerConfig{
					Name:               streamName,
					Durable:            streamName,
					Description:        "",
					DeliverPolicy:      0,
					OptStartSeq:        0,
					OptStartTime:       nil,
					AckPolicy:          0,
					AckWait:            0,
					MaxDeliver:         0,
					BackOff:            nil,
					FilterSubject:      "",
					ReplayPolicy:       0,
					RateLimit:          0,
					SampleFrequency:    "",
					MaxWaiting:         0,
					MaxAckPending:      0,
					HeadersOnly:        false,
					MaxRequestBatch:    0,
					MaxRequestExpires:  0,
					MaxRequestMaxBytes: 0,
					InactiveThreshold:  0,
					Replicas:           0,
					MemoryStorage:      false,
					FilterSubjects:     nil,
					Metadata:           nil,
				})
				if err != nil {
					slog.Error("create consumer", "err", err)
					panic(err)
				}
				slog.Info("NATS check, consumer created", "consume", streamName)

			}
			//slog.Error("check streams", "err", err, "stream", streamName)
			//return
		}

	}
}
