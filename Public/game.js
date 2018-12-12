new Vue({
    el: '#game',

    data: {
        ws: null,
        down: false
    },
    
    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/game-ws');
        this.ws.addEventListener('message', 
            function(e) {
                var msg = JSON.parse(e.data);
                var pos = JSON.parse(msg);
                var player = document.getElementById('player');
                player.style.left = pos.x + "px";
                player.style.top = pos.y + "px"; 
            });
        document.addEventListener('mousedown',
            function(e) {
                self.down = true;
            });
        document.addEventListener('mouseup',
            function(e) {
                self.down = false;
            });   
        document.addEventListener('mousemove',
            function(e) {
                if(self.down) {
                    self.ws.send(
                        JSON.stringify(
                            {
                                x: event.clientX,
                                y: event.clientY 
                            }
                        )
                    );
                }
            }, false);
    }
});