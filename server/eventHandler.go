package main

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"github.com/vishal1132/cafebucks/eventbus"
)

// eventHandler is the handler for events
func (s *server) eventHandler(msg kafka.Message, publish bool) {
	switch string(msg.Key) {
	case string(eventbus.OrderReceived):
		var event eventbus.EventC
		err := json.Unmarshal(msg.Value, &event)
		if err != nil {
			s.logger.Error().
				Err(err).
				Msg("error unmarshaling event")
			return
		}
		event.Event = eventbus.OrderAccept
		event.Order.Status = eventbus.OrderAccept
		b, err := json.Marshal(event)
		if err != nil {
			s.logger.Error().Err(err).Msg("error marshaling event to be pushed into kafka again")
		}
		if validateBeans(event.Order.Cof.Name) {
			if publish {
				err = s.EventBus.Publish(context.Background(), eventbus.OrderAccept, b)
			}
			if err != nil {
				s.logger.Error().Err(err).Msg("error pushing event to kafka")
			}
		}
	default:
		break
	}
}

func handle(eventType string, event []byte) {
	switch eventType {
	case string(eventbus.OrderReceived):

	}
}
