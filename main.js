const check = () => {
  if (!("serviceWorker" in navigator)) {
    throw new Error("No Service Worker support!");
  }
  if (!("PushManager" in window)) {
    throw new Error("No Push API Support!");
  }
};


const registerServiceWorker = async () => {
  const swRegistration = await navigator.serviceWorker.register("sw.js"); //notice the file name
  return swRegistration;
};

const requestNotificationPermission = async () => {
  const permission = await window.Notification.requestPermission();
  if (permission !== "granted") {
    throw new Error("Permission not granted for Notification");
  }
};

const getAPIData = async () => {
  let lastSale;
  // await axios.get("https://qcapi.nasdaq.com/api/quote/aapl/info?assetclass=stocks", -1).then((response) => {
  //   lastSale = response.data.data.primaryData.lastSalePrice;
  // });
  console.log("hey");
  return lastSale;
};

const showLocalNotification = (title, body, swRegistration) => {
  const options = {
    "//": "Visual Options",
    body,
  };
  swRegistration.showNotification(title, options);
};

const sendNotification = (swRegistration) => {
  let lastSale;
  let message = "";
  setInterval(async () => {
    let newData = await getAPIData();
    if (lastSale && lastSale !== newData) {
      message = "Price changed for apple from " + lastSale + " to " + newData;
      showLocalNotification("AAPL Stock changed", message, swRegistration);
    }
    else if (lastSale){
        message = "Last sale price was " + lastSale;
        showLocalNotification("AAPL Stock", message, swRegistration);
    }
    lastSale = newData;
  }, 10000);
  clearInterval();
  return message;
};

// Connect to the WebSocket server
const socket = io('http://localhost:8080');

// Listen for notifications
socket.on('connect', function() {
    socket.on('notification', function(data) {
        // Show the notification
        self.registration.showNotification(data.title, {
            body: data.message,
            icon: data.icon,
        });
    });
});

const main = async () => {
  check();

  const swRegistration = await registerServiceWorker();
  const permission = await requestNotificationPermission();
  sendNotification(swRegistration);
};
main();
