# poptape-messages
Messaging microservice written in Go with MongoDB

### API routes

```
/messages/notifications [GET] (Authenticated)
```

Returns a list of non-deleted notifications for the authenticated user.
Expected normal return codes: [200, 401]




### To Do:
* Refactor to use common code
* ~~Write notifications~~
* Write mails
* Write all tests
* ~~Dockerize~~
* Documentation
