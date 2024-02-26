package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"saga-kafka-notclean/money-transfer-service/config"
	model "saga-kafka-notclean/money-transfer-service/shared"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var VALIDATE_ACCOUNT_CHANNEL string = "validate_account_channel"
var LIMIT_CUT_CHANNEL string = "limit_cut_channel"
var MONEY_CUT_CHANNEL string = "money_cut_channel"
var UPDATE_BALANCE_CHANNEL string = "update_balance_channel"

func NewPostgresDatabase() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.GetConfig().Database.Postgres.Host,
		config.GetConfig().Database.Postgres.User,
		config.GetConfig().Database.Postgres.DBName,
		config.GetConfig().Database.Postgres.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)
	}

	return db
}

func Produce(topic string, input string) error {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		return err
	}

	// Produce messages to topic (asynchronously)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(input),
	}, nil)
	if err != nil {
		return err
	}

	// Wait for message deliveries
	p.Flush(15 * 1000)
	defer p.Close()
	return nil
}

func ValidateAccountSend(ctx context.Context, input model.SaferRequest) error {
	log.Println("💡Validate Account activity starts")

	s, err := json.Marshal(input)
	if err != nil {
		log.Println(err)
		return err
	}

	err = Produce(VALIDATE_ACCOUNT_CHANNEL, string(s))
	if err != nil {
		log.Println("🔥Cant send to %s", VALIDATE_ACCOUNT_CHANNEL)
		log.Println(err)
	}

	log.Println("💡Validate Account Send successfully")
	return err
}

func UpdateStateCreated(ctx context.Context, input model.TransactionEntity) error {
	log.Println("💡Persist transaction to database starts")

	db := NewPostgresDatabase()
	res := db.Create(input)
	if res.Error != nil {
		log.Println("🔥Cannot connect to database")
		log.Println(res.Error.Error())
		return res.Error
	}
	log.Println("💡Persist transaction to database successfully")

	return nil
}

func CompensateTransaction(ctx context.Context, input model.TransactionEntity) error {
	log.Println("💡Compensate Transaction starts")

	res := NewPostgresDatabase().Save(input)
	if res.Error != nil {
		log.Println("🔥Cannot create compensate transaction")
		log.Println(res.Error.Error())
		return res.Error
	}

	log.Println("💡Compensate Transaction successfully")
	return nil
}

func LimitCutSend(ctx context.Context, input model.SaferRequest) error {
	log.Println("💡Limit cut Account activity starts")

	s, err := json.Marshal(input)
	if err != nil {
		log.Println("Cant send to %s", LIMIT_CUT_CHANNEL)
		log.Println(err.Error())
		return err
	}
	Produce(LIMIT_CUT_CHANNEL, string(s))

	log.Println("💡Limit cut Account activity successfully")
	return nil
}

func LimitCutReceiver(ctx context.Context, input model.SaferResponse) error {
	return nil
}
