var React = require('react');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'SystemDeviceListGridCell',
    prefix: 'b-'
});
require('../../css/components/SystemDeviceListGridCell.less')

var SystemDeviceListGridCell = React.createClass({
    getInitialState: function() {
        return {
            checkboxChecked: this.props.checkboxChecked
        };
    },

    getDefaultProps: function() {
        return {
            //TODO: Set false
            showCheckbox: true,
            checkboxChecked: true
        };
    },

    checkboxClicked: function(evt) {
        evt.stopPropagation();
    },

    checkboxChanged: function(evt) {
        this.setState({checkboxChecked: evt.target.checked});
        this.props.chkBxChanged && this.props.chkBxChanged(evt.target.checked);
    },
    
    render: function() {
        var chkBx;
        if (this.props.showCheckbox) {
            chkBx = (
                <input
                    {...classes('checkbox')}
                    type="checkbox"
                    onChange={this.checkboxChanged}
                    onClick={this.checkboxClicked}
                    checked={this.state.checkboxChecked}
                ></input>
            );
        }

        return (
            <div {...classes()}>
                {chkBx}
                <div {...classes('icon')}>
                    <i className="icon ion-cube"></i>
                </div>
                <div {...classes('name')}>
                    {this.props.device.name}
                </div>
            </div>
        );
    }
});
module.exports = SystemDeviceListGridCell;
