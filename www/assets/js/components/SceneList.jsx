var React = require('react');
var ReactDOM = require('react-dom');
var Scene = require('./Scene.jsx');
var SceneInfo = require('./SceneInfo.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');

var SceneList = React.createClass({
    mixins: [UniqueIdMixin],
    
    getInitialState: function() {
        return {
            editMode: false,
            scenes: this.props.scenes,
            zones: this.props.zones,
        };
    },

    componentWillReceiveProps: function(nextProps) {
        if (nextProps.scenes) {
            this.setState({ scenes: nextProps.scenes });
        }
        if (nextProps.zones) {
            this.setState({ zones: nextProps.zones });
        }
    },
    
    componentDidMount: function() {

        //TODO: Needed?
        $.ajax({
            url: '/api/v1/systems/123/zones',
            dataType: 'json',
            cache: false,
            success: function(data) {
                this.setState({zones: data});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(err.toString());
            }.bind(this)
        });

        return;
        //TODO: Enable as part of a mode
        var el = ReactDOM.findDOMNode(this).getElementsByClassName('sceneList')[0];
        Sortable.create(el);
    },

    edit: function() {
        this.setState({ editMode: true });
    },

    endEdit: function() {
        this.setState({ editMode: false });
    },

    sceneDeleted: function(sceneId) {
        var scenes = this.state.scenes;
        for (var i=0; i<scenes.length; ++i) {
            if (scenes[i].id === sceneId) {
                scenes.splice(i, 1);
                this.setState({ scenes: scenes });
                break;
            }
        }
    },

    newScene: function() {
        //TODO: This is not the way...
        var scenes = this.state.scenes;
        scenes.unshift({ clientId: 'scenelist_' + this.getNextIdAndIncrement() + '' });
        this.setState({ scenes: scenes });
    },

    render: function() {
        var body;
        var btns;
        var self = this;
        if (this.state.editMode) {
            body = this.state.scenes.map(function(scene) {
                return (
                    <SceneInfo
                      onDestroy={self.sceneDeleted}
                      scenes={self.state.scenes}
                      zones={self.state.zones}
                      buttons={self.props.buttons}
                      scene={scene}
                      readOnlyFields="id"
                      key={scene.id || scene.clientId} />
                );
            });
            btns = (
                <div className="clearfix buttonWrapper">
                  <button className="btn btn-primary btnNew pull-left" onClick={this.newScene}>New Scene</button>
                  <button className="btn btn-success btnDone pull-right" onClick={this.endEdit}>Done</button>
                </div>
            );
        } else {
            body = this.state.scenes.map(function(scene) {
                return (
                    <Scene scene={scene} key={scene.id}/>
                );
            });
            btns = (
                <div className="clearfix buttonWrapper">
                  <button className="btn btn-primary btnEdit pull-right" onClick={this.edit}>Edit</button>
                </div>
            );
        }
        
        return (
            <div className="cmp-SceneList">
              {btns}
              {body}
            </div>
        );
    }
});
module.exports = SceneList;