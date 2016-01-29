var React = require('react');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')
var InputValidationMixin = require('./InputValidationMixin.jsx')
var SaveBtn = require('./SaveBtn.jsx');

module.exports = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        return {
            cid: this.getNextIdAndIncrement() + '',
            name: this.props.name || '',
            description: this.props.description || '',
            address: this.props.address,
            id: '',
            modelNumber: this.props.modelNumber || '',
            token: this.props.token,
            showToken: false,
            errors: null
        }
    },

    toJson: function() {
        var s = this.state;
        return {
            clientId: s.cid,
            name: s.name,
            description: s.description,
            address: s.address,
            modelNumber: s.modelNumber,
            token: s.token
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
    },

    testConnection: function() {
        //TODO: How to know what to call
    },

    save: function() {
        var saveBtn = this.refs.saveBtn;
        saveBtn.saving();
        this.setState({ errors: null });

        var self = this;
        $.ajax({
            url: '/api/v1/systems/1/devices',
            type: 'POST',
            dataType: 'json',
            contentType: 'application/json; charset=utf-8',
            data: JSON.stringify(this.toJson()),
            success: function(data) {
                saveBtn.success();
            },
            error: function(xhr, status, err) {
                self.setState({ errors: JSON.parse(xhr.responseText || '{}').errors});
                saveBtn.failure();
            }
        });            
    },
    
    render: function() {
        //TODO:need unique name for id and htmlFor
        var device = this.state.device;

        var token
        if (this.props.showToken) {
            token = (
                <div className={this.addErr("form-group", "token")}>
                  <label className="control-label" htmlFor={this.uid("token")}>Security Token</label>
                  <input value={this.state.token} data-statepath="token" onChange={this.changed} className="token form-control" type="text" id={this.uid("token")}/>
                  {this.errMsg('token')}
                </div>
            );
        }
        
        return (
            <div className="cmp-DeviceInfo well">
              <div className={this.addErr("form-group", "name")}>
                <label className="control-label" htmlFor={this.uid("name")}>Name</label>
                <input value={this.state.name} data-statepath="name" onChange={this.changed} className="name form-control" type="text" id="name"/>
                {this.errMsg("name")}
              </div>
              <div className={this.addErr("form-group", "description")}>
                <label className="control-label" htmlFor={this.uid("description")}>Description</label>
                <input value={this.state.description} data-statepath="description" onChange={this.changed} className="description form-control" type="text" id={this.uid("description")}/>
                {this.errMsg("description")}
              </div>
              <div className={this.addErr("form-group", "modelNumber")}>
                <label className="control-label" htmlFor={this.uid("modelNumber")}>Model Number</label>
                <input value={this.state.modelNumber} readOnly={this.isReadOnly("modelNumber")} data-statepath="modelNumber" onChange={this.changed} className="modelNumber form-control" type="text" id={this.uid("modelNumber")}/>
                {this.errMsg("modelNumber")}
              </div>
              <div className={this.addErr("form-group", "address")}>
                <label className="control-label" htmlFor={this.uid("address")}>Address</label>
                <input value={this.state.address} data-statepath="address" onChange={this.changed} className="address form-control" type="text" id={this.uid("address")}/>
                {this.errMsg("address")}
              </div>
              {token}
              <button className="btn btn-primary" onClick={this.testConnection}>Test Connection</button>
              <div className="pull-right">
                <SaveBtn ref="saveBtn" clicked={this.save} text="Import" />
              </div>
            </div>
        );
    }
});
