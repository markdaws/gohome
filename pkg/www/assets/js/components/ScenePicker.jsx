var React = require('react');

var ScenePicker = React.createClass({
    getInitialState: function() {
        return {
            value: this.props.sceneId || ''
        };
    },

    selected: function(evt) {
        this.setState({ value: evt.target.value });
        this.props.changed && this.props.changed(evt.target.value);
    },

    render: function() {
        var options = [];
        this.props.scenes.forEach(function(scene) {
            if (!scene.id) {
                // If this scene has not been saved it can't be used
                return;
            }

            // Can't set itself
            if (scene.id === this.props.parentSceneId) {
                return;
            }

            options.push(<option key={scene.id} value={scene.id}>{scene.name}</option>);
        }.bind(this));

        return (
            <div className="cmp-ScenePicker">
                <select
                    className="form-control"
                    disabled={this.props.disabled}
                    onChange={this.selected}
                    value={this.state.value}>
                <option value="">Select a Scene...</option>
                {options}
              </select>
            </div>
        );
    }
});
module.exports = ScenePicker;
