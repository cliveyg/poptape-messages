FROM ubuntu:18.04
# as base                                                                                                

RUN mkdir -p /messages
RUN mkdir -p /messages/schemas
COPY messages /messages
COPY .env /messages
COPY schemas/ /messages/schemas
WORKDIR /messages

# Make port 8090 available to the world outside this container
EXPOSE 8090

# Run app.py when the container launches
CMD ["./messages"]
