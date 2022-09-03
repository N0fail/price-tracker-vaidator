package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gitlab.ozon.dev/N0fail/price-tracker-validator/internal/error_codes"
	"gitlab.ozon.dev/N0fail/price-tracker-validator/internal/kafka/config"
	"gitlab.ozon.dev/N0fail/price-tracker-validator/internal/kafka/counter"
	"log"
	"time"
)

type RequestProducerI interface {
	ProductCreate(request config.ProductCreateRequest) error
	ProductDelete(request config.ProductDeleteRequest) error
	PriceTimeStampAdd(request config.PriceTimeStampAddRequest) error
	ProductList(request config.ProductListRequest) error
	PriceHistory(request config.PriceHistoryRequest) error
}

type RequestProducer struct {
	sp                 sarama.SyncProducer
	outRequestsCounter *counter.Counter
	errorsCounter      *counter.Counter
}

func (r *RequestProducer) ProductCreate(request config.ProductCreateRequest) error {
	r.outRequestsCounter.Inc()
	str, err := json.Marshal(request)
	if err != nil {
		r.errorsCounter.Inc()
		logrus.Errorf("kafka ProductCreate error in json.Marshal: %v", err.Error())
		return error_codes.ErrExternalProblem
	}

	_, _, err = r.sp.SendMessage(&sarama.ProducerMessage{
		Topic: config.ProductCreateTopic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%v", r.outRequestsCounter.Get())),
		Value: sarama.ByteEncoder(str),
	})

	if err != nil {
		r.errorsCounter.Inc()
		logrus.Errorf("kafka ProductCreate error in SendMessage: %v", err.Error())
		return error_codes.GetInternal(err)
	}

	logrus.Infof("success ProductCreate code: %v, name: %v", request.Code, request.Name)

	return nil
}

func (r *RequestProducer) ProductDelete(request config.ProductDeleteRequest) error {
	r.outRequestsCounter.Inc()
	str, err := json.Marshal(request)
	if err != nil {
		r.errorsCounter.Inc()
		logrus.Errorf("kafka ProductDelete error in json.Marshal: %v", err.Error())
		return error_codes.ErrExternalProblem
	}

	_, _, err = r.sp.SendMessage(&sarama.ProducerMessage{
		Topic: config.ProductDeleteTopic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%v", r.outRequestsCounter.Get())),
		Value: sarama.ByteEncoder(str),
	})

	if err != nil {
		r.errorsCounter.Inc()
		logrus.Errorf("kafka ProductDelete error in SendMessage: %v", err.Error())
		return error_codes.GetInternal(err)
	}

	logrus.Infof("success ProductDelete code: %v", request.Code)

	return nil
}

func (r *RequestProducer) PriceTimeStampAdd(request config.PriceTimeStampAddRequest) error {
	r.outRequestsCounter.Inc()
	str, err := json.Marshal(request)
	if err != nil {
		r.errorsCounter.Inc()
		logrus.Errorf("kafka PriceTimeStampAdd error in json.Marshal: %v", err.Error())
		return error_codes.ErrExternalProblem
	}

	_, _, err = r.sp.SendMessage(&sarama.ProducerMessage{
		Topic: config.PriceTimeStampAddTopic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%v", r.outRequestsCounter.Get())),
		Value: sarama.ByteEncoder(str),
	})

	if err != nil {
		r.errorsCounter.Inc()
		logrus.Errorf("kafka PriceTimeStampAdd error in SendMessage: %v", err.Error())
		return error_codes.GetInternal(err)
	}

	logrus.Infof("success PriceTimeStampAdd code: %v, price: %v, date: %v", request.Code, request.Price, time.Unix(request.Ts, 0).Format("2 Jan 2006 15:04"))

	return nil
}

func (r *RequestProducer) ProductList(request config.ProductListRequest) error {
	r.outRequestsCounter.Inc()
	str, err := json.Marshal(request)
	if err != nil {
		r.errorsCounter.Inc()
		logrus.Errorf("kafka ProductList error in json.Marshal: %v", err.Error())
		return error_codes.ErrExternalProblem
	}

	_, _, err = r.sp.SendMessage(&sarama.ProducerMessage{
		Topic: config.ProductListTopic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%v", r.outRequestsCounter.Get())),
		Value: sarama.ByteEncoder(str),
	})

	if err != nil {
		r.errorsCounter.Inc()
		logrus.Errorf("kafka ProductList error in SendMessage: %v", err.Error())
		return error_codes.GetInternal(err)
	}

	logrus.Infof("success ProductList")

	return nil
}

func (r *RequestProducer) PriceHistory(request config.PriceHistoryRequest) error {
	r.outRequestsCounter.Inc()
	str, err := json.Marshal(request)
	if err != nil {
		r.errorsCounter.Inc()
		logrus.Errorf("kafka PriceHistory error in json.Marshal: %v", err.Error())
		return error_codes.ErrExternalProblem
	}

	_, _, err = r.sp.SendMessage(&sarama.ProducerMessage{
		Topic: config.PriceHistoryTopic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%v", r.outRequestsCounter.Get())),
		Value: sarama.ByteEncoder(str),
	})

	if err != nil {
		r.errorsCounter.Inc()
		logrus.Errorf("kafka PriceHistory error in SendMessage: %v", err.Error())
		return error_codes.GetInternal(err)
	}

	logrus.Infof("success PriceHistory")

	return nil
}

func New() RequestProducerI {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(config.Brokers, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	requestProducer := &RequestProducer{
		sp:                 producer,
		outRequestsCounter: counter.New("outRequestsCounter"),
		errorsCounter:      counter.New("errorsCounter"),
	}
	return requestProducer
}
