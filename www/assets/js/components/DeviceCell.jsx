var React = require('react');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'DeviceCell',
    prefix: 'b-'
});
require('../../css/components/DeviceCell.less')

var DeviceCell = React.createClass({
    getInitialState: function() {
        return {
            checkboxChecked: this.props.checkboxChecked,
        };
    },

    getDefaultProps: function() {
        return {
            showCheckbox: false,
            checkboxChecked: true,
            hasError: false,
            hasSuccess: false
        };
    },

    isChecked: function() {
        return this.state.checkboxChecked;
    },

    checkboxClicked: function(evt) {
        evt.stopPropagation();
    },

    checkboxChanged: function(evt) {
        this.setState({checkboxChecked: evt.target.checked});
        this.props.chkBxChanged && this.props.chkBxChanged(this.props.id, evt.target.checked);
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

        var state = '';
        if (this.props.hasError) {
            state = 'error';
        } else if (this.props.hasSuccess) {
            state = 'success';
        }
        return (
            <div {...classes('', state)}>
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
module.exports = DeviceCell;
