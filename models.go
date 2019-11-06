package main

import (
    "time"
    //"log"
    "context"
    "os"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Notification struct {
	MessageId  string  `json:"message_id" bson:"_id"`
	LotId      string  `json:"lot_id" bson:"lot_id"`
    FromId     string  `json:"from_id" bson:"from_id"` // public_id of sender
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
    Starred    bool    `json:"starred" bson:"starred"`
}

type MessageMeta struct {
    TotalNotifications  int64 `json:"total_notifications" bson:"total_notifications"`
    UnreadNotifications int64 `json:"unread_notifications" bson:"unread_notifications"`
    TotalMails          int64 `json:"total_mails" bson:"total_mails"`
    UnreadMails         int64 `json:"unread_mails" bson:"unread_mails"`
}

// ----------------------------------------------------------------------------

func getMetadata(client *mongo.Client, publicId string) (MessageMeta, error) {

    var m MessageMeta

    nc := client.Database("poptape_messages").Collection("notifications")
    mc := client.Database("poptape_messages").Collection("mails")

    filter := bson.D{{"public_id", publicId},
                     {"read", false},
                     {"deleted", false}}
    nUnreadCount, enuc := nc.CountDocuments(context.TODO(), filter)
    if enuc != nil {
        return m, enuc
    }
    mUnreadCount, emuc := mc.CountDocuments(context.TODO(), filter)
    if emuc != nil {
        return m, emuc
    }

    filter = bson.D{{"public_id", publicId},
                    {"read", true},
                    {"deleted", false}}
    nReadCount, enrc := nc.CountDocuments(context.TODO(), filter)
    if enrc != nil {
        return m, enrc
    }
    mReadCount, emrc := mc.CountDocuments(context.TODO(), filter)
    if emrc != nil {
        return m, emrc
    }

    m.TotalNotifications = nUnreadCount + nReadCount
    m.UnreadNotifications = nUnreadCount
    m.TotalMails = mUnreadCount + mReadCount
    m.UnreadMails = mUnreadCount

    return m, nil
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

func getMail(client *mongo.Client, messageId, publicId string) (*Mail, error) {

    filter := bson.D{{"_id", messageId},
                     {"public_id", publicId}}
    var result *Mail
    var err error

    collection := client.Database("poptape_messages").Collection("mails")
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
        {"$set", bson.D{{"deleted", true}, {"delete_date", delete_datetime.Format(time.RFC3339)},
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

func setDeleteMail(client *mongo.Client, messageId, publicId string) (*mongo.UpdateResult, error) {

    filter := bson.D{{"_id", messageId},
                     {"public_id", publicId},
                     {"deleted", false}}
    delete_datetime := time.Now()

    update := bson.D{
        {"$set", bson.D{{"deleted", true}, {"delete_date", delete_datetime.String()},
        }},
    }

    collection := client.Database("poptape_messages").Collection("mails")
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

func getAllMails(client *mongo.Client, publicId string, limit int64) ([]*Mail, error) {

    var results []*Mail

    findOptions := options.Find()
    findOptions.SetLimit(limit)

    filter := bson.D{{"deleted", false},
                     {"public_id", publicId}}

    collection := client.Database("poptape_messages").Collection("mails")

    cur, err := collection.Find(context.TODO(), filter, findOptions)
    if err != nil {
        return results, err
    }

    for cur.Next(context.TODO()) {

        // create a value into which the single document can be decoded
        var m Mail
        err := cur.Decode(&m)
        if err != nil {
            return results, err
        }
        results = append(results, &m)
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
    create_datetime := time.Now()
    n.Created = create_datetime.Format(time.RFC3339)
    bsonInsert, err := bson.Marshal(n)
    if err != nil {
        return err
    }
    insertResult, err := collection.InsertOne(context.TODO(), bsonInsert)
    _ = insertResult
    if err != nil {
        return err
    }

    return nil
}

// ----------------------------------------------------------------------------

func (m *Mail) createMails(client *mongo.Client) error {
    os.Stderr.WriteString("[[ models.go - createMails() ]]")

    collection := client.Database("poptape_messages").Collection("mails")
    create_datetime := time.Now()
    m.Created = create_datetime.Format(time.RFC3339)
    bsonInsert, err := bson.Marshal(m)
    if err != nil {
        return err
    }
    insertResult, err := collection.InsertOne(context.TODO(), bsonInsert)
    _ = insertResult
    if err != nil {
        return err
    }

    return nil
}

// ----------------------------------------------------------------------------
