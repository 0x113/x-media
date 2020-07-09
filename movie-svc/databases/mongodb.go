package databases

import (
	"context"
	"fmt"
	"time"

	"github.com/0x113/x-media/movie-svc/common"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB manages MongoDB connection
type MongoDB struct {
	Session mongo.Session
	DbName  string
}

// Init initializes mongo database
func (db *MongoDB) Init() error {
	log.Infoln("Connecting to the mongo database ...")
	db.DbName = common.Config.DbName
	// mongodb connection options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", common.Config.DbAddr))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Errorf("Couldn't connect to the mongo database, err: %v", err)
		return err
	}

	// check the connection
	// FIXME: when using this with docker it doesn't wait 10 sec
	/*
		if err := client.Ping(ctx, nil); err != nil {
			log.Errorf("Unable to connect to the MongoDB; err: %v", err)
			return err
		}
	*/

	// create session
	session, err := client.StartSession()
	if err != nil {
		log.Errorf("Couldn't create mongo session, err: %v", err)
		return err
	}
	log.Infoln("Successfuly connected to the Mongo database")

	db.Session = session
	return nil
}
