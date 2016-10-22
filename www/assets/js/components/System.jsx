var React = require('react');
var Import = require('./Import.jsx');
var SystemDeviceList = require('./SystemDeviceList.jsx');

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
        var body, importBtn
        if (this.state.importing) {
            body = <Import/>
            importBtn = <button className="btn btn-danger pull-right btnExitImport" onClick={this.cancelImport}>Exit Import</button>
        } else {
            if (this.props.devices.length === 0) {
                body = <h5 className="emptyMessage">You don't have any devices. Click on the "Import" button to start, or you can manually update the .json file if you know what you are doing ;)</h5>
            } else {
                body = <SystemDeviceList devices={this.props.devices}/>                
            }

            importBtn = <button className="btn btn-primary" onClick={this.importProduct}>Import</button>
        }
        return (
            <div className="cmp-System">
              {importBtn}
              {body}
            </div>
        );
    }
});
module.exports = System;
