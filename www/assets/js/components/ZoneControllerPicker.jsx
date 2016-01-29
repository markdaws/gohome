var React = require('react');

var ZoneControllerPicker = React.createClass({
    getInitialState: function() {
        return {
            value: this.props.controller || ''
        };
    },

    selected: function(evt) {
        this.setState({ value: evt.target.value });
        this.props.changed && this.props.changed(evt.target.value);
    },
    
    render: function() {
        return (
            <div className="cmp-ZoneControllerPicker">
              <select className="form-control" onChange={this.selected} defaultValue={this.props.controller} value={this.state.value}>
                <option value="">Default</option>
                <option value="FluxWIFI">Flux WIFI</option>
              </select>
            </div>
        );
    }
});
module.exports = ZoneControllerPicker;