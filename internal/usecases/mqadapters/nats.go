package mqadapters

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/AleksandrVishniakov/tgbots-tgapi/models"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	streamName = "UPDATES"
	updatesSubject = "updates.>"
)

type NATSAdapter struct {
	log    *slog.Logger
	conn   *nats.Conn
	js     jetstream.JetStream
	stream jetstream.Stream
}

func NewNATSAdapter(
	ctx context.Context,
	log *slog.Logger,
	conn *nats.Conn,
) (*NATSAdapter, error){
	n :=  &NATSAdapter{
		log:  log,
		conn: conn,
	}

	err := n.init(ctx)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (n *NATSAdapter) PublishUpdate(ctx context.Context, subjectID int64, update *models.Update) error {
	const src = "NATSAdapter"
	log := n.log.With(slog.String("src", src))

	data, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("%s marshal data: %w", src, err)
	}

	_, err = n.js.Publish(ctx, fmt.Sprintf("updates.%d", subjectID), data)
	if err != nil {
		return fmt.Errorf("%s publish update: %w", src, err)
	}

	log.DebugContext(ctx, "update published",
		slog.Int64("subjectID", subjectID),
		slog.Int("updateID", update.UpdateID),
	)

	return nil
}

func (n *NATSAdapter) init(ctx context.Context) error {
	const src = "NATSAdapter.init"
	log := n.log.With(slog.String("src", src))

	js, err := jetstream.New(n.conn)
	if err != nil {
		return fmt.Errorf("%s create JetStream: %w", src, err)
	}
	n.js = js

	cfg := jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{updatesSubject},
		Storage:  jetstream.FileStorage,
	}

	stream, err := js.CreateStream(ctx, cfg)
	if err != nil {
		return fmt.Errorf("%s create stream: %w", src, err)
	}
	n.stream = stream

	log.DebugContext(ctx, "create stream", slog.String("name", streamName))
	return nil
}
