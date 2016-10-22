var React = require('react');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var ZonePicker = require('./ZonePicker.jsx');
var Api = require('../utils/API.js');
var ClassNames = require('classnames');

var ZoneSetLevelCommand = module.exports = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        var attr = this.props.command.attributes;
        return {
            clientId: this.getNextIdAndIncrement() + '',
            level: attr.Level || 0,
            r: attr.R || 0,
            g: attr.G || 0,
            b: attr.B || 0,
            zoneId: this.props.command.attributes.ZoneID || '',
            zoneOutput: '',
            errors: null,
        }
    },

    toJson: function() {
        return {
            type: 'zoneSetLevel',
            clientId: this.state.clientId,
            //TODO: correctly capitalize json values
            attributes: {
                Level: parseFloat(this.state.level),
                R: parseInt(this.state.r, 10),
                G: parseInt(this.state.g, 10),
                B: parseInt(this.state.b, 10),
                ZoneID: this.state.zoneId
            }
        };
    },

    setErrors: function(errors) {
        this.setState({ errors: errors });
    },

    zonePickerChanged: function(zone) {
        this.setState({
            zoneId: zone.id,
            zoneOutput: zone.output });
    },

    testLevel: function() {
        if (!this.state.zoneId) {
            return;
        }

        Api.zoneSetLevel(
            this.state.zoneId,
            'setLevel',
            parseFloat(this.state.level),
            parseInt(this.state.r, 10),
            parseInt(this.state.g, 10),
            parseInt(this.state.b, 10),
            function(err, data) {
                //TODO: error
            });
    },

    render: function() {
        //TODO: Only show RGB if this is an OTRGB
        //TODO: Insert RGB Picker in UI as well
        //TODO: For binary outputs should have a picker on/off not 0-100
        return (
            <div className="cmp-ZoneSetLevelCommand">
                <h4>Zone Set Level</h4>
                <div className={this.addErr("form-group", "attributes_ZoneID")}>
                    <label className="control-label" htmlFor={this.uid("attributes_ZoneID")}>Zone*</label>
                    <ZonePicker
                        disabled={this.props.disabled}
                        changed={this.zonePickerChanged}
                        zones={this.props.zones}
                        zoneId={this.state.zoneId} />
                    {this.errMsg("attributes_ZoneID")}
                </div>
                <div className={this.addErr("form-group", "attributes_Level")}>
                    <label className="control-label" htmlFor={this.uid("attributes_Level")}>Level [0-100]</label>
                    <div className="input-group">
                        <input
                            disabled={this.props.disabled}
                            value={this.state.level}
                            data-statepath="level"
                            onChange={this.changed}
                            className="level form-control"
                            type="number"
                            id={this.uid("attributes_Level")}/>
                        <span className="input-group-btn">
                            <button className="btn btn-primary" onClick={this.testLevel}>
                                Test Level
                            </button>
                        </span>
                        {this.errMsg("attributes_Level")}
                    </div>
                </div>
                <div className={ClassNames({
                               clearfix: true,
                               rgbExpander: true,
                               hidden: this.state.zoneOutput !== 'rgb' })}>
                    <a data-toggle="collapse" href={"#" + this.uid("rgbExpand")}>
                        RGB
                        <i className="glyphicon glyphicon-menu-down"></i>
                    </a>
                </div>
                <div className="collapse rbgExpand" id={this.uid("rgbExpand")}>
                    <p><strong>NOTE:</strong> To set R/G/B values, leave the "Value" field set to 0. If "Value" is non-zero then the R/G/B values are ignored and instead R/G/B will all be set to 255 * (Value/100)</p>
                    <div className={this.addErr("form-group", "attributes_R")}>
                        <label className="control-label" htmlFor={this.uid("attributes_R")}>Level - Red [0-255]</label>
                        <input
                            disabled={this.props.disabled}
                            value={this.state.r}
                            data-statepath="r"
                            onChange={this.changed}
                            className="r form-control"
                            type="number"
                            id={this.uid("attributes_R")}/>
                        {this.errMsg("attributes_R")}
                    </div>
                    <div className={this.addErr("form-group", "attributes_G")}>
                        <label className="control-label" htmlFor={this.uid("attributes_G")}>Level - Green [0-255]</label>
                        <input
                            disabled={this.props.disabled}
                            value={this.state.g}
                            data-statepath="g"
                            onChange={this.changed}
                            className="g form-control"
                            type="number"
                            id={this.uid("attributes_G")}/>
                        {this.errMsg("attributes_G")}
                    </div>
                    <div className={this.addErr("form-group", "attributes_B")}>
                        <label className="control-label" htmlFor={this.uid("attributes_B")}>Level - Blue [0-255]</label>
                        <input
                            disabled={this.props.disabled}
                            value={this.state.b}
                            data-statepath="b"
                            onChange={this.changed}
                            className="b form-control"
                            type="number"
                            id={this.uid("attributes_B")}/>
                        {this.errMsg("attributes_B")}
                    </div>
                </div>
            </div>
        );
    }
});
module.exports = ZoneSetLevelCommand;
