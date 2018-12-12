new Vue({
    el: '#game',

    data: {
        ws: null
    },
    
    created: function() {
        this.ws = new WebSocket('ws://' + window.location.host + '/game-ws');
        this.ws.addEventListener('message', 
            function(e) {
                var msg = JSON.parse(e.data);
                var pos = JSON.parse(msg);
                var player = document.getElementById('player');
                player.style.left = pos.x + "px";
                player.style.top = pos.y + "px"; 
            });
    },

    methods: { }
});