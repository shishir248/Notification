import {PushNotificationClient} from "./public/notification_grpc_web_pb.js";
import {Subscription, Notification, Response} from "./public/notification_pb.js"

const client = new PushNotificationClient("http://localhost:50051");

export function subscribe(email) {
  const request = new Subscription();
  request.setEmail(email);
  client.subscribe(request, {}, (err, response) => {
    if (err) {
        console.error(err);
    } else {
        console.log(response.getMessage());
    }
  });
}

export function unsubscribe(email) {
  const request = new Subscription();
  request.setEmail(email);
  client.unsubscribe(request, {}, (err, response) => {
    if (err) {
        console.error(err);
    } else {
        console.log(response.getMessage());
    }
  });
}

export function receiveNotification() {
  const call = client.send();
  call.on("data", (notification) => {
    console.log(notification.getTitle(), notification.getMessage());
  });
  call.on("error", (err) => {
    console.error(err);
  });
}
