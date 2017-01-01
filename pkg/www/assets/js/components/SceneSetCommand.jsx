var React = require('react');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var ScenePicker = require('./ScenePicker.jsx');
var uuid = require('uuid');

var SceneSetCommand = module.exports = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        return {
            clientId: uuid.v4(),
            sceneId: this.props.command.attributes.SceneID || '',
        }
    },

    getDefaultProps: function() {
        return {
            scenes: []
        }
    },

    toJson: function() {
        return {
            type: 'sceneSet',
            clientId: this.state.clientId,
            attributes: {
                SceneID: this.state.sceneId
            }
        };
    },

    scenePickerChanged: function(sceneId) {
        this.setState({ sceneId: sceneId });
        this.props.onChanged && this.props.onChanged();
    },

    render: function() {
        return (
            <div className="cmp-SceneSetCommand">
              <h4>Scene Set</h4>
              <div className={this.addErr("form-group", "attributes_SceneID")}>
                <ScenePicker
                    disabled={this.props.disabled}
                    changed={this.scenePickerChanged}
                    scenes={this.props.scenes}
                    sceneId={this.state.sceneId}
                    parentSceneId={this.props.parentSceneId}/>
                {this.errMsg("attributes_SceneID")}
              </div>
            </div>
        );
    }
});
module.exports = SceneSetCommand;
