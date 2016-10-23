var React = require('react');
var ReactRedux = require('react-redux');
var Import = require('./Import.jsx');
var SystemDeviceList = require('./SystemDeviceList.jsx');

var System = React.createClass({
    getInitialState: function() {
        return {
            importing: false,
        };
    },

    newClicked: function() {
        this.props.deviceNew();
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
            body = <Import/>
            header = <button className="btn btn-danger pull-right btnExitImport" onClick={this.cancelImport}>Exit Import</button>
        } else {
            if (this.props.devices.length === 0) {
                body = <h5 className="emptyMessage">You don't have any devices. Click on the "Import" button to start, or you can manually update the .json file if you know what you are doing ;)</h5>
            } else {
                body = <SystemDeviceList devices={this.props.devices}/>                
            }

            header = (
                <div className="header clearfix">
                    {/* TODO: Add back when this is really needed"
                    <button className="btn btn-default pull-right" onClick={this.newClicked}>
                        <i className="ion-plus"></i>
                    </button>*/}
                    <button className="btn btn-default pull-right" onClick={this.importProduct}>
                        <i className="ion-arrow-down-c"></i>
                    </button>
                </div>
            );
        }
        return (
            <div className="cmp-System">
              {header}
              {body}
            </div>
        );
    }
});

/*
function mapDispatchToProps(dispatch) {
    return {
        deviceNew: function() {
            dispatch(SystemActions.deviceNew());
        }
    };
}
module.exports = ReactRedux.connect(null, mapDispatchToProps)(System);
*/
module.exports = System;
