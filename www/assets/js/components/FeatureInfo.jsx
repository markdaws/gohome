var React = require('react');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var FeatureTypePicker = require('./FeatureTypePicker.jsx');
var SaveBtn = require('./SaveBtn.jsx');
var Api = require('../utils/API.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'FeatureInfo',
    prefix: 'b-'
});
require('../../css/components/FeatureInfo.less')

var FeatureInfo = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],

    getDefaultProps: function() {
        return {
            showSaveBtn: false
        };
    },

    getInitialState: function() {
        var f = this.props.feature
        return {
            id: f.id,
            aid: f.aid,
            name: f.name,
            description: f.description,
            address: f.address,
            attrs: f.attrs,
            type: f.type,
            errors: this.props.errors,
            dirty: false,
            saveButtonStatus: '',
        }
    },

    toJson: function() {
        var s = this.state;
        var f = this.props.feature;
        return {
            id: f.id,
            aid: s.aid,
            name: s.name,
            description: s.description,
            address: s.address,
            type: s.type,
            deviceId: f.deviceId,
            type: s.type,
            attrs: s.attrs
        }
    },

    setErrors: function(errors) {
        this.setState({ errors: errors });
    },

    _changed: function(evt) {
        this.setState({ saveButtonStatus: '' }, function() {
            this.props.changed && this.props.changed(this);
        }.bind(this));
        this.changed(evt);
    },

    typeChanged: function(type) {
        this.changed({
            target: {
                getAttribute: function() {
                    return 'type'
                },
                value: type
            }
        }, function() {
            this.props.changed && this.props.changed(this);
        });
    },

    save: function() {
        this.setState({ errors: null });
        Api.featureUpdate(this.props.feature.deviceId, this.props.feature.id, this.toJson(), function(err, featureData) {
            if (err && !err.validation) {
                //TODO: Dispatch general error
                this.setState({ saveButtonStatus: 'error' })
                return;
            } else if (err && err.validation) {
                this.setState({
                    saveButtonStatus: 'error',
                    errors: err.validation.errors[this.state.id]
                });
                return;
            }

            this.setState({ saveButtonStatus: 'success' });
            this.props.updatedFeature(featureData);
        }.bind(this));
    },

    render: function() {
        var saveBtn;
        if (this.props.showSaveBtn && this.state.dirty) {
            saveBtn = (
                <SaveBtn
                    {...classes('save')}
                    clicked={this.save}
                    text="Save"
                    status={this.state.saveButtonStatus}/>
            );
        }

        return (
            <div {...classes('', '', 'well well-sm')}>
              <div className={this.addErr('form-group', 'name')}>
                <label {...classes('label', '', 'control-label')} htmlFor={this.uid('name')}>Name*</label>
                <input
                    value={this.state.name}
                    data-statepath="name"
                    onChange={this._changed}
                    className="name form-control"
                    type="text"
                    id={this.uid('name')} />
                {this.errMsg('name')}
              </div>
              <div className={this.addErr("form-group", "type")}>
                  <label {...classes('label', '', 'control-label')} htmlFor={this.uid("type")}>Type</label>
                  <FeatureTypePicker type={this.state.type} changed={this.typeChanged}/>
                  {this.errMsg('type')}
              </div>
              <div className={this.addErr("form-group", "id")}>
                  <label {...classes('label', '', 'control-label')} htmlFor={this.uid("id")}>ID</label>
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
              <div className={this.addErr("form-group", 'description')}>
                  <label {...classes('label', '', 'control-label')} htmlFor={this.uid("description")}>Description</label>
                  <input
                      value={this.state.description}
                      data-statepath="description"
                      onChange={this._changed}
                      className="description form-control"
                      type="text"
                      id={this.uid("description")} />
                  {this.errMsg('description')}
              </div>
              <div className={this.addErr("form-group", "address")}>
                  <label {...classes('label', '', 'control-label')} htmlFor={this.uid("address")}>Address</label>
                  <input
                      value={this.state.address}
                      data-statepath="address"
                      onChange={this._changed}
                      className="address form-control"
                      type="text"
                      id={this.uid("address")} />
                  {this.errMsg('address')}
              </div>
              <div className={this.addErr("form-group", "aid")}>
                  <label {...classes('label', '', 'control-label')} htmlFor={this.uid("aid")}>Automation ID</label>
                  <input
                      value={this.state.aid}
                      readOnly={this.isReadOnly("aid")}
                      data-statepath="aid"
                      onChange={this._changed}
                      className="id form-control"
                      type="text"
                      id={this.uid("aid")} />
                  {this.errMsg("aid")}
              </div>
              <div className="pull-right">
                  {saveBtn}
              </div>
              <div style={{clear: 'both' }}></div>
            </div>
        );
    }
});
module.exports = FeatureInfo;
