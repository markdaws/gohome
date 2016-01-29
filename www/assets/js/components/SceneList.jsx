var React = require('react');
var ReactDOM = require('react-dom');
var Scene = require('./Scene.jsx');

module.exports = React.createClass({
    componentDidMount: function() {
        return;
        //TODO: Enable as part of a mode
        var el = ReactDOM.findDOMNode(this).getElementsByClassName('sceneList')[0];
        Sortable.create(el);
    },

    render: function() {
        //TODO: Add loading
        var self = this;
        var sceneNodes = Object.keys(this.props.scenes).map(function(id) {
            var scene = self.props.scenes[id];
            return (
                <Scene scene={scene} key={id}/>
            );
        });
        return (
            <div className="cmp-SceneList row">
              {sceneNodes}
            </div>
        );
    }
});
