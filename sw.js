self.addEventListener('push', function(event) {
    var data = event.data.json();
    var title = data.title;
    var options = {
        body: data.message,
        icon: data.icon
    };
    event.waitUntil(self.registration.showNotification(title, options));
});
