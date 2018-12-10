new Vue({
    el: '#game',

    data: {
        ws: null,
        X: 0,
        Y: 0
    },
    created: function() {
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('position', function(e) {
            var msg = JSON.parse(e.data);
            var canvas = document.getElementById("canvas");
            var context = canvas.getContext("2d");
            context.fillRect(msg.X, msg.Y, 150, 100);
        });
    }
    //,
    // methods: {
    //     send: function () {
    //         if (this.newMsg != '') {
    //             this.ws.send(
    //                 JSON.stringify({
    //                     y: 0
    //                     x: 0
    //                 }
    //             ));
    //         }
    //     }
    // }
});