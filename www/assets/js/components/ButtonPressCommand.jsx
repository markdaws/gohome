var React = require('react');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var ButtonPicker = require('./ButtonPicker.jsx');
var uuid = require('uuid');

var ButtonPressCommand = module.exports = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],

    getInitialState: function() {
        return {
            id: uuid.v4(),
            buttonId: this.props.command.attributes.ButtonID || '',
            errors: null,
        }
    },

    getDefaultProps: function() {
        return {
            buttons: []
        }
    },

    toJson: function() {
        return {
            type: 'buttonPress',
            id: this.state.id,
            attributes: {
                ButtonID: this.state.buttonId
            }
        };
    },

    setErrors: function(errors) {
        this.setState({ errors: errors });
    },

    buttonPickerChanged: function(buttonId) {
        this.setState({ buttonId: buttonId });
    },

    render: function() {
        return (
            <div className="cmp-ButtonPressCommand">
              <h4>Button Press</h4>
              <div className={this.addErr("form-group", "attributes_ButtonID")}>
                <label className="control-label" htmlFor={this.uid("attributes_ButtonID")}>Button*</label>
                <ButtonPicker
                    disabled={this.props.disabled}
                    changed={this.buttonPickerChanged}
                    buttons={this.props.buttons}
                    buttonId={this.state.buttonId} />
                {this.errMsg("attributes_ButtonID")}
              </div>
            </div>
        );
    }
});
module.exports = ButtonPressCommand;
