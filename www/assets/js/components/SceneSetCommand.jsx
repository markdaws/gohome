var React = require('react');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var ScenePicker = require('./ScenePicker.jsx');

var SceneSetCommand = module.exports = React.createClass({
    mixins: [UniqueIdMixin, InputValidationMixin],
    getInitialState: function() {
        return {
            cid: this.getNextIdAndIncrement() + '',
            sceneId: this.props.command.attributes.SceneID || '',
            errors: null,
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
            clientId: this.state.cid,
            attributes: {
                SceneID: this.state.sceneId
            }
        };
    },

    setErrors: function(errors) {
        this.setState({ errors: errors });
    },

    scenePickerChanged: function(sceneId) {
        this.setState({ sceneId: sceneId });
    },
    
    render: function() {
        return (
            <div className="cmp-SceneSetCommand">
              <h4>Scene Set</h4>
              <div className={this.addErr("form-group", "attributes_SceneID")}>
                <label className="control-label" htmlFor={this.uid("attributes_SceneID")}>Scene*</label>
                <ScenePicker changed={this.scenePickerChanged} scenes={this.props.scenes} sceneId={this.state.sceneId} />
                {this.errMsg("attributes_SceneID")}
              </div>
            </div>
        );
    }
});
module.exports = SceneSetCommand;