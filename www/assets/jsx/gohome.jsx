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
        handleClick: function(event) {
            $.ajax({
                url: '/api/v1/systems/1/scenes/active',
                type: 'POST',
                dataType: 'json',
                contnetType: 'application/json; charset=utf-8',
                data: JSON.stringify({ id: this.props.id }),
                success: function(data) {
                    console.log('set the scene');
                }.bind(this),
                error: function(xhr, status, err) {
                    console.error(err.toString());
                }.bind(this)
            });
        },

        render: function() {
            return (
                <div>
                    <span>{this.props.id} : {this.props.name} : {this.props.description}</span>
                    <a onClick={this.handleClick}> [Set]</a>
                </div>
            )
        }
    });

    var apiUrl='/api/v1/systems/123/scenes';
    React.render(<App url={apiUrl}/>, document.body);
})();
