var React = require('react');
var Import = require('./Import.jsx');
var DeviceList = require('./DeviceList.jsx');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'System',
    prefix: 'b-'
});
require('../../css/components/System.less')

var System = React.createClass({
    getInitialState: function() {
        return {
            importing: false,
        };
    },

    importProduct: function() {
        this.setState({ importing: true });
    },

    cancelImport: function() {
        this.setState({ importing: false });
    },

    render: function() {
        var body, header
        if (this.state.importing) {
            body = <Import/>;
            header = (
                <button {...classes('exit', '', 'btn btn-default pull-right')} onClick={this.cancelImport}>
                    <i className="fa fa-times"></i>
                </button>
            );
        } else {
            if (this.props.devices.length === 0) {
                body = (
                    <h5 {...classes('empty-message')}>
                        You don't have any hardware. Click on the "Import" button to get started.
                    </h5>
                );
            } else {
                body = <DeviceList devices={this.props.devices}/>
            }

            header = (
                <div {...classes('header', '', 'clearfix')}>
                    <button className="btn btn-default pull-right" onClick={this.importProduct}>
                        <i className="ion-arrow-down-c"></i>
                    </button>
                </div>
            );
        }
        return (
            <div {...classes()}>
              {header}
              {body}
            </div>
        );
    }
});
module.exports = System;
