package main

import (
)

func (a *App) initializeRoutes() {

    // endpoints
    a.Router.HandleFunc("/messages/status", a.getStatus).Methods("GET")
    //a.Router.HandleFunc("/messages/info", a.getMessageInfo).Methods("GET")
    //a.Router.HandleFunc("/messages/all", a.getAllMyMessages).Methods("GET")
    //a.Router.HandleFunc("/messages/mails", a.getAllMyMails).Methods("GET")
	//a.Router.HandleFunc("/messages/mails", a.createMessages).Methods("POST")
    //a.Router.HandleFunc("/messages/mails/{messageId}", a.getMailMessage).Methods("GET")
    a.Router.HandleFunc("/messages/notifications", a.createNotifications).Methods("POST")
    a.Router.HandleFunc("/messages/notifications", a.getAllMyNotifications).Methods("GET")
    a.Router.HandleFunc("/messages/notifications/{messageId}", a.getNotification).Methods("GET")
	a.Router.HandleFunc("/messages/notifications/{messageId}", a.deleteNotification).Methods("DELETE")

}
