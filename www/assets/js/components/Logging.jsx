var React = require('react');
var ReactDOM = require('react-dom');
var LogLine = require('./LogLine.jsx');

var Logging = React.createClass({
    getInitialState: function() {
        return {
            items: [],
            connectionStatus: 'connecting'
        };
    },

    componentDidMount: function() {
        this.reconnect();
    },

    componentDidUpdate: function() {
        var lastLi = this.refs.lastLi;
        if (!lastLi) {
            return;
        }

        //TODO: Shouldn't set the body element like this, use events
        //TODO: If the user has scrolled away from the bottom, don't do this
        //until they scroll back to the bottom again, annoying to jump away
        $('body')[0].scrollTop = ReactDOM.findDOMNode(lastLi).offsetTop;
    },

    componentWillUnmount: function() {
        var conn = this.state.conn;
        if (!conn) {
            return;
        }
        conn.Close();
    },

    reconnect: function() {
        var oldConn = this.state.conn;
        if (oldConn) {
            oldConn.close();
        }

        var conn = new WebSocket("ws://" + window.location.host + "/api/v1/events/ws");
        var self = this;
        conn.onopen = function(evt) {
            self.setState({
                connectionStatus: 'connected'
            });
        };
        conn.onclose = function(evt) {
            conn = null;
            self.setState({
                conn: null,
                items: [],
                connectionStatus: 'disconnected'
            });
        };
        conn.onmessage = function(evt) {
            var item = JSON.parse(evt.data);
            item.datetime = new Date(item.datetime);
            self.setState({ items: self.state.items.concat(item)});
        };
        this.setState({
            conn: conn,
            connectionStatus: 'connecting'
        });

        //TODO: Fetch X previous log items from server?
    },

    clearClicked: function() {
        this.setState({ items: [] });
    },

    render: function() {
        var body;

        switch(this.state.connectionStatus) {
        case 'connected':
            var itemCount = this.state.items.length;
            body = this.state.items.map(function(item, i) {
                return <LogLine item={item} key={item.id} ref={itemCount === i+1 ? 'lastLi' : undefined}/>;
            });
            break;

        case 'connecting':
            body = <li className="spinner"><i className="fa fa-spinner fa-spin"></i></li>
            break;

        case 'disconnected':
            body = <li className="reconnect"><button className="btn btn-primary" onClick={this.reconnect}>Reconnect</button></li>
            break;
        }

        var hasEvents = this.state.items.length > 0;
        var waiting = !hasEvents && this.state.connectionStatus === 'connected';
        return (
            <div className="cmp-Logging">
              <h3 className={!waiting ? 'hidden' : ''}>Waiting for events...</h3>
              <ol className="list-unstyled">
                {body}
              </ol>
              <div className="footer text-center">
                <button className={(hasEvents ? '' : 'hidden') + ' btn btn-default'} onClick={this.clearClicked}>Clear</button>
              </div>
            </div>
        );
    }
});
module.exports = Logging;