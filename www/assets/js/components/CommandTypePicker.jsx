var React = require('react');

var CommandTypePicker = React.createClass({
    getInitialState: function() {
        return {
            value: ''
        };
    },

    selected: function(evt) {
        this.setState({ value: '' });
        this.props.changed && this.props.changed(evt.target.value);
    },
    
    render: function() {
        return (
            <div className="cmp-CommandTypePicker">
              <select className="form-control" onChange={this.selected} value={this.state.value}>
                <option value="">Select...</option>
                <option value="buttonPress">Button Press</option>
                <option value="buttonRelease">Button Release</option>
                <option value="sceneSet">Scene Set</option>
                <option value="zoneSetLevel">Zone Set Level</option>
              </select>
            </div>
        );
    }
});
module.exports = CommandTypePicker;

//TODO: turnOn/turnOff