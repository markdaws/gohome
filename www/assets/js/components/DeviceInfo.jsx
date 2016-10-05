var React = require('react');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')
var InputValidationMixin = require('./InputValidationMixin.jsx')
var SaveBtn = require('./SaveBtn.jsx');
var Api = require('../utils/API.js');

var DeviceInfo = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        //TODO: need state?
        return {
            name: this.props.name || '',
            description: this.props.description || '',
            address: this.props.address,
            id: this.props.id,
            clientId: this.props.clientId,
            modelNumber: this.props.modelNumber || '',
            token: this.props.token,
            showToken: false,
            errors: null,
            saveButtonStatus: ''
        }
    },

    toJson: function() {
        var s = this.state;
        return {
            clientId: this.props.clientId,
            name: s.name,
            description: s.description,
            address: s.address,
            modelNumber: s.modelNumber,
            token: s.token,
            id: s.id,
        };
    },

    componentWillReceiveProps: function(nextProps) {
        var device = this.state.device;
        if (nextProps.name != "") {
            this.setState({ name: nextProps.name });
        }
        if (nextProps.description != "") {
            this.setState({ description: nextProps.description });
        }
        if (nextProps.address != "") {
            this.setState({ address: nextProps.address });
        }
        if (nextProps.token != "") {
            this.setState({ token: nextProps.token });
        }
        if (nextProps.id != "" ) {
            this.setState({ id: nextProps.id });
        }
        if (nextProps.clientId != "") {
            this.setState({ clientId: nextProps.clientId });
        }
    },

    testConnection: function() {
        //TODO: How to know what to call
    },

    save: function() {
        this.setState({ errors: null });

        Api.deviceCreate(this.toJson(), function(err, data) {
            if (err) {
                this.setState({
                    saveButtonStatus: 'error',
                    errors: err.validationErrors
                });
                return;
            }
            
            //TODO: Update list of devices with response from server, via redux
            this.setState({ saveButtonStatus: 'success' });

            //TODO: Update list of devices with saved device information
            //this.props.savedDevice(data);
        }.bind(this));
    },

    deleteDevice: function() {
        this.props.deviceDelete(this.state.id, this.state.clientId);
    },

    _changed: function(evt) {
        this.setState({ saveButtonStatus: '' });
        this.changed(evt);
    },

    render: function() {
        var device = this.state.device;

        var token
        if (this.props.showToken) {
            token = (
                <div className={this.addErr("form-group", "token")}>
                    <label className="control-label" htmlFor={this.uid("token")}>Security Token</label>
                    <input
                        value={this.state.token}
                        data-statepath="token"
                        onChange={this._changed}
                        className="token form-control"
                        type="text"
                        id={this.uid("token")} />
                    {this.errMsg('token')}
                </div>
            );
        }
        
        return (
            <div className="cmp-DeviceInfo well">
                <button className="btn btn-link btnDelete pull-right" onClick={this.deleteDevice}>
                    <i className="glyphicon glyphicon-trash"></i>
                </button>
                <div className={this.addErr("form-group", "name")}>
                    <label className="control-label" htmlFor={this.uid("name")}>Name*</label>
                    <input
                        value={this.state.name}
                        data-statepath="name"
                        onChange={this._changed}
                        className="name form-control"
                        type="text"
                        id={this.uid("name")} />
                    {this.errMsg("name")}
                </div>
                <div className={this.addErr("form-group", "id")}>
                    <label className="control-label" htmlFor={this.uid("id")}>ID</label>
                    <input
                        value={this.state.id}
                        readOnly={this.isReadOnly("id")}
                        data-statepath="id"
                        onChange={this._changed}
                        className="id form-control"
                        type="text"
                        id={this.uid("id")} />
                    {this.errMsg("id")}
                </div>
                <div className={this.addErr("form-group", "description")}>
                    <label className="control-label" htmlFor={this.uid("description")}>Description</label>
                    <input
                        value={this.state.description}
                        data-statepath="description"
                        onChange={this._changed}
                        className="description form-control"
                        type="text"
                        id={this.uid("description")} />
                    {this.errMsg("description")}
                </div>
                <div className={this.addErr("form-group", "modelNumber")}>
                    <label className="control-label" htmlFor={this.uid("modelNumber")}>Model Number</label>
                    <input
                        value={this.state.modelNumber}
                        readOnly={this.isReadOnly("modelNumber")}
                        data-statepath="modelNumber"
                        onChange={this._changed}
                        className="modelNumber form-control"
                        type="text"
                        id={this.uid("modelNumber")} />
                    {this.errMsg("modelNumber")}
                </div>
                <div className={this.addErr("form-group", "address")}>
                    <label className="control-label" htmlFor={this.uid("address")}>Address</label>
                    <input
                        value={this.state.address}
                        data-statepath="address"
                        onChange={this._changed}
                        className="address form-control"
                        type="text"
                        id={this.uid("address")} />
                    {this.errMsg("address")}
                </div>
                {token}
                <button className="btn btn-primary" onClick={this.testConnection}>Test Connection</button>
                <div className="pull-right">
                    <SaveBtn
                        clicked={this.save}
                        text="Save"
                        status={this.state.saveButtonStatus}/>
                </div>
            </div>
        );
    }
});
module.exports = DeviceInfo;
