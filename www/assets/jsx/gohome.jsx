(function() {

    var App = React.createClass({
        getInitialState: function() {
            return { scenes: [] };
        },

        componentDidMount: function() {
            $.ajax({
                url: this.props.url,
                dataType: 'json',
                cache: false,
                success: function(data) {
                    this.setState({scenes: data});
                }.bind(this),
                error: function(xhr, status, err) {
                    console.error(err.toString());
                }.bind(this)
            });
        },
        
        render: function() {
            return (
                <SceneList scenes={this.state.scenes} />
            );
        }
    });

    var SceneList = React.createClass({
        render: function() {
            var sceneNodes = this.props.scenes.map(function(scene) {
                return (
                    <Scene id={scene.id} name={scene.name} description={scene.description} />
                );
            });
            return (
                <div className="sceneList">
                    {sceneNodes}
                </div>
            );
        }
    });

    var Scene = React.createClass({
        render: function() {
            return (
                <div>{this.props.id} : {this.props.name} : {this.props.description}</div>
            )
        }
    });

    var scenes = [
        {id:1, name:"one", description:"desc1"},
        {id:2, name:"two", description:"desc2"},
        {id:3, name:"three", description:"desc3"}
    ];
    var apiUrl='/api/systems/123/scenes';
    React.render(<App url={apiUrl}/>, document.body);
})();
