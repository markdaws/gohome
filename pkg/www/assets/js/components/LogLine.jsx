var React = require('react');

var LogLine = React.createClass({
    render: function() {
        return (
            <li className="cmp-LogLine">
              <span className="datetime">{this.props.item.datetime.toLocaleString()}</span>
              <span className="deviceName"> [{this.props.item.deviceName}]</span>
              <span> : {this.props.item.friendlyMessage}</span>
              <span className="rawMessage"> [Raw: {this.props.item.rawMessage}]</span>
            </li>
        );
    }
});
module.exports = LogLine;
