package main

import (
    "encoding/json"
    "net/http"
    "log"
    "context"
    "io"
    "io/ioutil"
    "bytes"
    "github.com/gorilla/mux"
    "strings"
    "github.com/google/uuid"
    "github.com/xeipuuv/gojsonschema"
	"fmt"
)

// ----------------------------------------------------------------------------

func (a *App) getStatus(w http.ResponseWriter, r *http.Request) {

    err := a.Mongo.Ping(context.TODO(), nil)

    if err != nil {
        log.Fatal(err)
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
        return
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	mess := `{"message": "System running and ping successful..."}`
    io.WriteString(w, mess)
}

// ----------------------------------------------------------------------------

func (a *App) getMessageInfo(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    b, st, mess := bouncerSaysOk(r)
    if !b {
        w.WriteHeader(st)
        io.WriteString(w, mess)
        return
    }
    // if successful mess contains only the public_id
    publicId := mess

    metadata, err := getMetadata(a.Mongo, publicId)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
        return
    }

    jsonData, _ := json.Marshal(metadata)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)

}


// ----------------------------------------------------------------------------

func (a *App) createNotifications(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    b, st, mess := bouncerSaysOk(r)
    if !b {
        w.WriteHeader(st)
        io.WriteString(w, mess)
        return
    }
    // if successful mess contains only the public_id
    publicId := mess

    // read the request body so we can validate using json schema
    // we also want to limit how big the read should be to stop
    // bad guys trying to use massive request bodys
    limitedRead := io.LimitReader(r.Body, 10000)
    body, err := ioutil.ReadAll(limitedRead)
    if err != nil {
        log.Fatal(err)
    }

    schemaLoader := gojsonschema.NewReferenceLoader(a.NotificationURI)
    documentLoader := gojsonschema.NewStringLoader(string(body))
    result, err := gojsonschema.Validate(schemaLoader, documentLoader)
    if err != nil {
        log.Print(err)
        w.WriteHeader(http.StatusBadRequest)
        io.WriteString(w, "")
        return
    }
    if !result.Valid() {
        // if we have any errors then build a
        // list of them to return to the client
        schErrs := make(map[string]string)
        for _, err := range result.Errors() {
            log.Print(err)
            if err.Field() == "(root)" {
                s := strings.Split(err.Description(), " ")
                if s[0] != "Additional" {
                    schErrs[s[0]] = "missing field"
                } else {
                    schErrs[s[2]] = "this additional property not allowed"
                }
            } else {
                schErrs[err.Field()] = err.Description()
            }
        }
        errsJson, err := json.Marshal(schErrs)
        if err != nil {
            log.Fatal(err)
        }
        w.WriteHeader(http.StatusBadRequest)
        w.Write(errsJson)
        return
    }

    // rebuild the request body as it gets consumed by reading
    r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

    // in theory since we get here after json schema validation
    // there should not be any errors...
    var n Notification
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&n); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        mess := fmt.Sprintf("{ \"error\": \"%s\" }",err)
        io.WriteString(w, mess)
        return
    }
    defer r.Body.Close()

    messageId, err := uuid.NewRandom()
    if err !=nil {
        log.Fatal(err)
    }

    n.MessageId = messageId.String()
    n.FromId = publicId

    if err := n.createNotifications(a.Mongo); err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
        return
    }

    w.WriteHeader(http.StatusCreated)
    io.WriteString(w, `{ "message": "you no worry, everything's tickety-boo" }`)

}

// ----------------------------------------------------------------------------

func (a *App) getAllMyNotifications(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    b, st, mess := bouncerSaysOk(r)
    if !b {
        w.WriteHeader(st)
		io.WriteString(w, mess)
		return
    }
	// successfully authenticated which means mess is the public_id
	publicId := mess

	notifs, err := getAllNotifications(a.Mongo, publicId, a.DocumentLimit)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
		return
	}

    jsonData, _ := json.Marshal(notifs)
    if len(notifs) == 0 {
        w.WriteHeader(http.StatusNotFound)
        io.WriteString(w, `{}`)
        return
    } else {
        w.WriteHeader(http.StatusOK)
    }
    w.Write(jsonData)

}

// ----------------------------------------------------------------------------

func (a *App) getAllMyMails(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    b, st, mess := bouncerSaysOk(r)
    if !b {
        w.WriteHeader(st)
        io.WriteString(w, mess)
        return
    }
    // successfully authenticated which means mess is the public_id
    publicId := mess

    mails, err := getAllMails(a.Mongo, publicId, a.DocumentLimit)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
        return
    }

    jsonData, _ := json.Marshal(mails)
    if len(mails) == 0 {
        w.WriteHeader(http.StatusNotFound)
        io.WriteString(w, `{}`)
        return
    } else {
        w.WriteHeader(http.StatusOK)
    }
    w.Write(jsonData)

}

// ----------------------------------------------------------------------------

func (a *App) getNotification(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    b, st, mess := bouncerSaysOk(r)
    if !b {
        w.WriteHeader(st)
        io.WriteString(w, mess)
        return
    }
    // if check is successful then mess contains only the public_id
    publicId := mess

    vars := mux.Vars(r)
    messageId := vars["messageId"]

	if !IsValidUUID(messageId) {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{ "message": "Not a valid ID" }`)
		return
	}

    n, err := getNotification(a.Mongo, messageId, publicId)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
        return
    }

    if n == nil || n.Deleted == true {
        w.WriteHeader(http.StatusNotFound)
        return
    }

	jsonData, _ := json.Marshal(n)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}

// ----------------------------------------------------------------------------

func (a *App) getMail(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    b, st, mess := bouncerSaysOk(r)
    if !b {
        w.WriteHeader(st)
        io.WriteString(w, mess)
        return
    }
    // if check is successful then mess contains only the public_id
    publicId := mess

    vars := mux.Vars(r)
    messageId := vars["messageId"]

    if !IsValidUUID(messageId) {
        w.WriteHeader(http.StatusBadRequest)
        io.WriteString(w, `{ "message": "Not a valid ID" }`)
        return
    }

    m, err := getMail(a.Mongo, messageId, publicId)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
        return
    }

    if m == nil || m.Deleted == true {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    jsonData, _ := json.Marshal(m)
    w.WriteHeader(http.StatusOK)
    w.Write(jsonData)

}

// ----------------------------------------------------------------------------

func (a *App) deleteNotification(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    b, st, mess := bouncerSaysOk(r)
    if !b {
        w.WriteHeader(st)
        io.WriteString(w, mess)
        return
    }
	publicId := mess

    vars := mux.Vars(r)
    messageId := vars["messageId"]

    if !IsValidUUID(messageId) {
        w.WriteHeader(http.StatusBadRequest)
        io.WriteString(w, `{ "message": "Not a valid ID" }`)
        return
    }

	res, err := setDeleteNotification(a.Mongo, messageId, publicId)
    if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
        return
    }

    if res.ModifiedCount == 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusGone)

}

// ----------------------------------------------------------------------------

func (a *App) deleteMail(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    b, st, mess := bouncerSaysOk(r)
    if !b {
        w.WriteHeader(st)
        io.WriteString(w, mess)
        return
    }
    publicId := mess

    vars := mux.Vars(r)
    messageId := vars["messageId"]

    if !IsValidUUID(messageId) {
        w.WriteHeader(http.StatusBadRequest)
        io.WriteString(w, `{ "message": "Not a valid ID" }`)
        return
    }

    res, err := setDeleteMail(a.Mongo, messageId, publicId)
    if err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
        return
    }

    if res.ModifiedCount == 0 {
        w.WriteHeader(http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusGone)

}

// ----------------------------------------------------------------------------

func (a *App) createMails(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")

    b, st, mess := bouncerSaysOk(r)
    if !b {
        w.WriteHeader(st)
        io.WriteString(w, mess)
        return
    }
    // if successful mess contains only the public_id
    publicId := mess

    // read the request body so we can validate using json schema
    // we also want to limit how big the read should be to stop
    // bad guys trying to use massive request bodys
    limitedRead := io.LimitReader(r.Body, 10000)
    body, err := ioutil.ReadAll(limitedRead)
    if err != nil {
        log.Fatal(err)
    }

    schemaLoader := gojsonschema.NewReferenceLoader(a.MailURI)
    documentLoader := gojsonschema.NewStringLoader(string(body))
    result, err := gojsonschema.Validate(schemaLoader, documentLoader)
    if err != nil {
        log.Print(err)
        w.WriteHeader(http.StatusBadRequest)
        io.WriteString(w, "")
        return
    }
    if !result.Valid() {
        // if we have any errors then build a
        // list of them to return to the client
        schErrs := make(map[string]string)
        for _, err := range result.Errors() {
            log.Print(err)
            if err.Field() == "(root)" {
                s := strings.Split(err.Description(), " ")
                if s[0] != "Additional" {
                    schErrs[s[0]] = "missing field"
                } else {
                    schErrs[s[2]] = "this additional property not allowed"
                }
            } else {
                schErrs[err.Field()] = err.Description()
            }
        }
        errsJson, err := json.Marshal(schErrs)
        if err != nil {
            log.Fatal(err)
        }
        w.WriteHeader(http.StatusBadRequest)
        w.Write(errsJson)
        return
    }

    // rebuild the request body as it gets consumed by reading
    r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

    // in theory since we get here after json schema validation
    // there should not be any errors...
    //TODO: change this so we can use a common struct and make both
    // notifications and mails use the same code
    var m Mail
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&m); err != nil {
        w.WriteHeader(http.StatusBadRequest)
        mess := fmt.Sprintf("{ \"error\": \"%s\" }",err)
        io.WriteString(w, mess)
        return
    }
    defer r.Body.Close()

    // create message uuid
    messageId, err := uuid.NewRandom()
    if err !=nil {
        log.Fatal(err)
    }

    m.MessageId = messageId.String()
    m.FromId = publicId

    if err := m.createMails(a.Mongo); err != nil {
        log.Print(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        io.WriteString(w, `{ "message": "Oopsy somthing went wrong" }`)
        return
    }

    w.WriteHeader(http.StatusCreated)
    io.WriteString(w, `{ "message": "you no worry, everything's tickety-boo" }`)

}
