var React = require('react');
var Trigger = require('./Trigger.jsx');

module.exports = React.createClass({
    handleClick: function(trigger) {
        this.props.selected(trigger);
    },

    render: function() {
        var self = this;
        var triggerNodes = this.props.triggers.map(function(trigger) {
            return (
                <Trigger data={trigger} selected={self.handleClick} key={trigger.name} />
            );
        });

        return (
            <div className="cmp-TriggerList clearfix">
              {triggerNodes}
            </div>
        );
    }
});
