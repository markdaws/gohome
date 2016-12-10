var React = require('react');
var AutomationCell = require('./AutomationCell.jsx');
var Automation = require('./Automation.jsx');
var Grid = require('./Grid.jsx');
var BEMHelper = require('react-bem-helper');
var Feature = require('../feature.js');

var classes = new BEMHelper({
    name: 'AutomationList',
    prefix: 'b-'
});
require('../../css/components/AutomationList.less')

var AutomationList = React.createClass({
    getDefaultProps: function() {
        return {
            automations: []
        };
    },

    render: function() {
        var gridCells = this.props.automations.map(function(automation) {
            return {
                // Note: automations don't have an ID
                key: automation.name,
                cell: <AutomationCell automation={automation} />,
                content: <Automation automation={automation} key={automation.name}/>
            };
        });

        return (
            <div {...classes()}>
                <h2 {...classes('header')}>Automation</h2>
                <Grid cells={gridCells} />
            </div>
        );
    }
});

module.exports = AutomationList;
