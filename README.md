# Notifications system

The system consist of `frontend`, `backend`, `router`, and a simple `slack-service`

- Frontend - serves html form for sending messages
- Backend - accepts http request with message and writes to kafka notifications topic
- Router - service reads from notifications topic and distributes the message to appropriate channel topics
- Slack-service - reads from the slack topic and send the message to Slack via webhook url

Each of the services could be scaled horizontally

Delivery is assured by Kafka.
- If backend service is down it produces an error for frontend.
- If any of the other services is down messages are accumulated by kafka and will be consumed when the service becomes available.


# Usage
```
git clone git@github.com:bitmaskit/notifications.git
cd notifications
make
```