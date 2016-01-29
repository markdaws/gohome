var React = require('react');
var UniqueIdMixin = require('./UniqueIdMixin.jsx')

var Ingredient = React.createClass({
    mixins: [UniqueIdMixin],
    getInitialState: function() {
        return {
            value: undefined
        };
    },

    changeHandler: function(evt) {
        this.setState({ value: evt.target.value });
    },

    render: function() {
        var input;
        switch(this.props.data.type) {
        case 'string':
        case 'duration':
        case 'integer':
        case 'float':
            input = <input className="ingredientInput form-control" type="text" onChange={this.changeHandler} id={this.getNextIdAndIncrement()}/>;
            break;
        case 'boolean':
            input = <input className="ingredientInput form-control" type="checkbox" value="true" onChange={this.changeHandler} id={this.getNextIdAndIncrement()}/>;
            break;
        case 'datetime':
            //TODO: show calendar
            break;
        default:
            throw 'unknown ingredient type: ' + this.props.data.type;
        }

        var err = this.props.err;
        var errDesc = err ? err.description : '';
        return (
            <div>
              <div className={"form-group" + (err ? " has-error" : "")}>
                <label className="control-label" htmlFor={this.getCurrentId()}>{this.props.data.name}</label>
                <p>{this.props.data.description}</p>
                {input}
                <span className={"help-block" + (err ? "" : " invisible")}>Error - {errDesc}</span>
              </div>
            </div>
        );
    },

    value: function() {
        if (this.state.value == undefined) {
            return undefined;
        }
        
        switch(this.props.data.type) {
        case 'string':
            return this.state.value;
        case 'integer':
        case 'duration':
            return parseInt(this.state.value, 10);
        case 'float':
            return parseFloat(this.state.value)
        case 'boolean':
            return this.state.value === true || this.state.value === 'true';
            break;
        case 'datetime':
            //TODO:
            break;
        default:
            throw 'Unknown data type: ' + this.props.data.type;
        }
    }
});
module.exports = Ingredient;