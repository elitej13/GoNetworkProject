new Vue({
    el: '#chat',

    data: {
        ws: null,
        newMsg: '',        
        chatContent: '', 
        username: null,     
        color: '',
        joined: false, 
        id: ''
    },
    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/chat-ws');
        this.color = 'hsla(' + (Math.random() * 300 + 30) + ', 100%, 70%, 1)';
        this.id = '' + Math.random() * 360;
        this.ws.addEventListener('message', 
            function(e) {
                var msg = JSON.parse(e.data);
                self.chatContent += '<div class="messages"'
                + 'style="background:' + msg.color + ';"'
                + '>&nbsp;&nbsp;'
                + '[' + msg.username + '] '
                + msg.message
                + '</div>';
                var es = document.getElementsByClassName()
                var element = document.getElementById('chat-window');
                element.scrollTop = element.scrollHeight;
            });
    },
    methods: {
        send: function () {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        username: this.username,
                        color: this.color,
                        message: $('<p>').html(this.newMsg).text() // Strip out html
                    }
                ));
                this.newMsg = '';
            }
        },
        join: function () {
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        }
    }
});