var React = require('react');
var ReactDOM = require('react-dom');
var Attribute = require('../attribute.js');
var BrightnessAttr = require('./BrightnessAttr.jsx');
var OnOffAttr = require('./OnOffAttr.jsx');
var TempAttr = require('./TempAttr.jsx');
var HueAttr = require('./HueAttr.jsx');
var OffsetAttr = require('./OffsetAttr.jsx');
var OpenClosedAttr = require('./OpenClosedAttr.jsx');
var Feature = require('../feature.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'FeatureControl',
    prefix: 'b-'
});
require('../../css/components/FeatureControl.less')

var FeatureControl = React.createClass({
    getDefaultProps: function() {
        return {
            hideReadOnlyAttrs: false
        };
    },

    getInitialState: function() {
        this._modifiedAttrs = {};

        // Need to make sure we clone the attributes so we don't affect other
        // instances of the feature control that are using this feature
        return {
            attrs: this.props.attrs || Feature.cloneAttrs(this.props.feature.attrs)
        }
    },

    componentDidMount: function() {
        this.props.didMount && this.props.didMount(this);
    },

    componentWillUnmount: function() {
        this.props.willUnmount && this.props.willUnmount(this);
    },

    monitorData: function(data) {
        if (!data || !data.features) {
            return;
        }
        var attrs = data.features[this.props.id];
        if (attrs == undefined) {
            return;
        }
        this.setState({attrs: attrs});
    },

    setAttrs: function(attr, value) {
        //TODO: This should just keep track of the attributes, how to have an unset
        //state for all UI elements, including slider?
        var newAttr = Object.assign({}, attr)
        newAttr.value = value;

        // Keep track of all the changes
        this._modifiedAttrs[newAttr.localId] = newAttr;

        this.props.onAttrChanged && this.props.onAttrChanged(this.props.feature, newAttr);
    },

    modifiedAttrs: function() {
        return this._modifiedAttrs;
    },

    render: function() {
        var attributes = [];
        var onOff, openClosed;

        // Features can have multiple attributes. For example a lightzone feature can have
        // onoff, brightness, hue attributes, each one has a LocalID so that we can distinguish
        // between the various attributes
        var localIDs = Object.keys(this.state.attrs);

        localIDs.forEach(function(localID) {
            var attribute = this.state.attrs[localID];

            // Check to see if we should skip read only attributes
            if (this.props.hideReadOnlyAttrs && attribute.perms === Attribute.Perms.ReadOnly) {
                return
            }

            // For each type of attribute we have a component that can render
            // the specific type of data e.g. a slider for brightness, rgb
            // picker for the hue
            switch(attribute.type) {
                case Attribute.Type.OnOff:
                    onOff = (
                        <div {...classes('on-off', '', 'pull-right')}>
                            <OnOffAttr
                                onToggle={this.setAttrs}
                                key={localID}
                                attr={attribute} />
                        </div>
                    );
                    break;
                case Attribute.Type.OpenClose:
                    openClosed = (
                        <div {...classes('open-closed', '', 'pull-right')}>
                            <OpenClosedAttr
                                onToggle={this.setAttrs}
                                key={localID}
                                attr={attribute} />
                        </div>
                    );
                    break;
                case Attribute.Type.Brightness:
                    attributes.push(
                        <BrightnessAttr
                            onBrightnessChanged={this.setAttrs}
                            key={localID}
                            attr={attribute} />
                    );
                    break;

                case Attribute.Type.Hue:
                    attributes.push(
                        <HueAttr
                            onHueChanged={this.setAttrs}
                            key={localID}
                            attr={attribute} />
                    );
                    break;

                case Attribute.Type.Offset:
                    attributes.push(
                        <OffsetAttr
                            onOffsetChanged={this.setAttrs}
                            key={localID}
                            attr={attribute} />
                    );
                    break;

                case Attribute.Type.Temperature:
                    attributes.push(
                        <TempAttr
                            onTempChanged={this.setAttrs}
                            key={localID}
                            attr={attribute} />
                    );
                    break;

                default:
                    console.error('unknown attribute type: ' + attribute.type);
            }
        }.bind(this));

        return (
            <div {...classes('')}>
                <div className="clearfix">
                    <div {...classes('name', '', 'pull-left')}>
                        {this.props.feature.name}
                    </div>
                    {onOff}
                    {openClosed}
                </div>
                {attributes}
            </div>
        );
    }
});
module.exports = FeatureControl;
