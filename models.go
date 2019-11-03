package main

import (
    "time"
    "log"
    "context"
    "os"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Notification struct {
	MessageId  string  `json:"message_id" bson:"_id"`
	LotId      string  `json:"lot_id" bson:"lot_id"`
    PublicId   string  `json:"public_id" bson:"public_id"` // notification recipient
    PurchaseId string  `json:"purchase_id" bson:"purchase_id"`
	AuctionId  string  `json:"auction_id" bson:"auction_id"`
    ItemId     string  `json:"item_id" bson:"item_id"`
    Priority   int     `json:"priority" bson:"priority"`
    MessStatus string  `json:"message_status" bson:"message_status"`
	Deleted    bool    `json:"deleted" bson:"deleted"`
    DeleteDate string  `json:"delete_date" bson:"delete_date"`
    Read       bool    `json:"read" bson:"read"`
    ReadDate   string  `json:"read_date" bson:"read_date"`
	Message    string  `json:"message" bson:"message"`
	Created    string  `json:"created" bson:"created"`
    Type       int     `json:"type" bson:"type"`
    Label      string  `json:"label" bson:"label"`
}

type Mail struct {
    MessageId  string  `json:"message_id" bson:"_id"`
    PrevMessId string  `json:"previous_message_id" bson:"previous_message_id"` // enables message chains
    FromId     string  `json:"from_id" bson:"from_id"` // public_id of sender
    PublicId   string  `json:"public_id" bson:"public_id"` // notification recipient
    LotId      string  `json:"lot_id" bson:"lot_id"`
    PurchaseId string  `json:"purchase_id" bson:"purchase_id"`
    AuctionId  string  `json:"auction_id" bson:"auction_id"`
    ItemId     string  `json:"item_id" bson:"item_id"`
    Priority   int     `json:"priority" bson:"priority"`
    MessStatus string  `json:"message_status" bson:"message_status"`
    Deleted    bool    `json:"deleted" bson:"deleted"`
    DeleteDate string  `json:"delete_date" bson:"delete_date"`
    Read       bool    `json:"read" bson:"read"`
    ReadDate   string  `json:"read_date" bson:"read_date"`
    Message    string  `json:"message" bson:"message"`
    Created    string  `json:"created" bson:"created"`
    Type       int     `json:"type" bson:"type"`
    Subject    string  `json:"subject" bson:"subject"`
}

// ----------------------------------------------------------------------------

func getNotification(client *mongo.Client, messageId, publicId string) (*Notification, error) {

    filter := bson.D{{"_id", messageId},
                     {"public_id", publicId}}
    var result *Notification
    var err error

    collection := client.Database("poptape_messages").Collection("notifications")
    collection.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil {
        return nil, err
    }

    return result, nil
}

// ----------------------------------------------------------------------------

func setDeleteNotification(client *mongo.Client, messageId, publicId string) (*mongo.UpdateResult, error) {

    filter := bson.D{{"_id", messageId},
                     {"public_id", publicId},
                     {"deleted", false}}
    delete_datetime := time.Now()

    update := bson.D{
        {"$set", bson.D{{"deleted", true}, {"delete_date", delete_datetime.String()},
        }},
    }

    collection := client.Database("poptape_messages").Collection("notifications")
    updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        return nil, err
    }
    return updateResult, nil
}

// ----------------------------------------------------------------------------

func getAllNotifications(client *mongo.Client, publicId string, limit int64) ([]*Notification, error) {

    var results []*Notification

    findOptions := options.Find()
    findOptions.SetLimit(limit)

    filter := bson.D{{"deleted", false},
                     {"public_id", publicId}}

    collection := client.Database("poptape_messages").Collection("notifications")

    cur, err := collection.Find(context.TODO(), filter, findOptions)
    if err != nil {
        return results, err
    }

    for cur.Next(context.TODO()) {

        // create a value into which the single document can be decoded
        var n Notification
        err := cur.Decode(&n)
        if err != nil {
            return results, err
        }
        results = append(results, &n)
    }

    if err := cur.Err(); err != nil {
        return results, err
    }

    // close the cursor once finished
    cur.Close(context.TODO())

    return results, nil

}

// ----------------------------------------------------------------------------

func (n *Notification) createNotifications(client *mongo.Client) error {
    os.Stderr.WriteString("[[ models.go - createNotifications() ]]")

    collection := client.Database("poptape_messages").Collection("notifications")
    bsonInsert, err := bson.Marshal(n)
    if err != nil {
        return err
    }
    //log.Print(string(bsonInsert))
    insertResult, err := collection.InsertOne(context.TODO(), bsonInsert)
    log.Print(insertResult)
    if err != nil {
        return err
    }

    return nil
}

// ----------------------------------------------------------------------------
