(function() {

    var App = React.createClass({
        getInitialState: function() {
            return { scenes: [], zones: [] };
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

            $.ajax({
                url: this.props.zoneUrl,
                dataType: 'json',
                cache: false,
                success: function(data) {
                    this.setState({zones: data});
                }.bind(this),
                error: function(xhr, status, err) {
                    console.error(err.toString());
                }.bind(this)
            });
        },
        
        render: function() {
            return (
                <div>
                    <SceneList scenes={this.state.scenes} />
                    <ZoneList zones={this.state.zones} />
                </div>
            );
        }
    });

    var ZoneList = React.createClass({
        render: function() {
            var self = this;
            var zoneNodes = Object.keys(this.props.zones).map(function(id) {
                var zone = self.props.zones[id];
                return (
                    <Zone id={zone.id} name={zone.name} />
                );
            })
            return (
                <div className="zoneList">
                    <h1>Zones</h1>
                    {zoneNodes}
                </div>
            );
        }
    });

    var Zone = React.createClass({
        getInitialState: function() {
            return { value: 100 }
        },

        clickHandler: function(event) {
            $.ajax({
                url: '/api/v1/systems/1/zones/' + this.props.id,
                type: 'POST',
                dataType: 'json',
                contnetType: 'application/json; charset=utf-8',
                data: JSON.stringify({ value: parseFloat(this.state.value) }),
                success: function(data) {
                    console.log('set the zone');
                }.bind(this),
                error: function(xhr, status, err) {
                    console.error(err.toString());
                }.bind(this)
            });
        },

        handleChange: function(event) {
            this.setState({ value: event.target.value });
        },

        render: function() {
            var value = this.state.value;
            return (
                <div>
                    <span>{this.props.id} : {this.props.name}</span>
                    <input type="text" value={value} onChange={this.handleChange}></input>
                    <a onClick={this.clickHandler}> [Set]</a>
                </div>
            )
        }
    });

    var SceneList = React.createClass({
        render: function() {
            var self = this;
            var sceneNodes = Object.keys(this.props.scenes).map(function(id) {
                var scene = self.props.scenes[id];
                return (
                    <Scene id={scene.id} name={scene.name} description={scene.description} />
                );
            });
            return (
                <div className="sceneList">
                    <h1>Scenes</h1>
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

    var apiUrl = '/api/v1/systems/123/scenes';
    var apiUrlZones = '/api/v1/systems/123/zones';
    React.render(<App url={apiUrl} zoneUrl={apiUrlZones}/>, document.body);
})();
