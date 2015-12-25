(function() {

    var App = React.createClass({
        getInitialState: function() {
            return { scenes: [], zones: [] };
        },

        componentDidMount: function() {
            //TODO: Have a loading indicator for scenes + zones
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
                    <Zone id={zone.id} name={zone.name} type={zone.type}/>
                );
            })
            return (
                <div>
                    <a className="zoneListHeader" data-toggle="collapse" href=".zoneList">Zones</a>
                    <div className="collapse zoneList row">
                        {zoneNodes}
                    </div>
                </div>
            );
        }
    });

    var Zone = React.createClass({
        getInitialState: function() {
            return { value: 100 }
        },

        handleChange: function(event) {
            this.setState({ value: event.target.value });
        },

        handleClick: function(event) {
            //TODO: Modal on desktop?
            return;
            console.log(ReactDOM.findDOMNode(this.refs.zoneModal));
            $(ReactDOM.findDOMNode(this.refs.zoneModal)).modal();
        },

        render: function() {
            var value = this.state.value;
            return (
                <div className="col-xs-12 col-sm-3 col-md-3 col-lg-3">
                    <a data-toggle="collapse" href={".what" + this.props.id} className="btn btn-default zone" onClick={this.handleClick}>
                        <div>
                            <i className="fa fa-lightbulb-o"></i>
                        </div>
                        <span className="name">{this.props.name} ({this.props.type})</span>
                        <input style={{display: 'none'}} type="text" value={value} onChange={this.handleChange}></input>
                    </a>
                    <ZoneModal ref="zoneModal" name={this.props.name} id={this.props.id}/>
                </div>
            )
        }
    });

    var ZoneModal = React.createClass({
        componentDidMount: function() {
            var s = $(ReactDOM.findDOMNode(this)).find('.valueSlider');
            s.slider({ reversed: true});
            var i = 0;
            var self = this;
            s.on('slideStop', function(evt) {
                $.ajax({
                    url: '/api/v1/systems/1/zones/' + self.props.id,
                    type: 'POST',
                    dataType: 'json',
                    contentType: 'application/json; charset=utf-8',
                    data: JSON.stringify({ value: parseFloat(evt.value) }),
                    success: function(data) {
                        console.log('set the zone');
                    }.bind(self),
                    error: function(xhr, status, err) {
                        console.error(err.toString());
                    }.bind(self)
                });
            });
        },

        render: function() {
            return (
                <div className={"collapse zoneModal " + " what" + this.props.id}>
                    <div className="well">
                        <div className="content">
                            <h3>{this.props.name}</h3>
                            <input className="valueSlider" type="text" data-slider-value="0" data-slider-min="00" data-slider-max="100" data-slider-step="1" data-slider-orientation="vertical"></input>
                        </div>
                    </div>
                </div>
            );
        }
    });

    var SceneList = React.createClass({
        componentDidMount: function() {
            return;
            //TODO: Enable as part of a mode
            var el = ReactDOM.findDOMNode(this).getElementsByClassName('sceneList')[0];
            Sortable.create(el);
        },

        render: function() {
            var self = this;
            var sceneNodes = Object.keys(this.props.scenes).map(function(id) {
                var scene = self.props.scenes[id];
                return (
                    <Scene id={scene.id} name={scene.name} description={scene.description} />
                );
            });
            return (
                <div>
                    <a className="sceneListHeader" data-toggle="collapse" href=".sceneList">Scenes</a>
                    <div className="collapse sceneList row">
                        {sceneNodes}
                    </div>
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
                contentType: 'application/json; charset=utf-8',
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
                <div className="col-xs-6 col-sm-3 col-md-3 col-lg-3">
                    <a className="btn btn-default scene" onClick={this.handleClick}>
                <div>
                <i className="fa fa-sliders"></i>
                        </div>
                        <span className="name">{this.props.name}</span>
                    </a>
                </div>
            )
        }
    });

    var apiUrl = '/api/v1/systems/123/scenes';
    var apiUrlZones = '/api/v1/systems/123/zones';
    React.render(<App url={apiUrl} zoneUrl={apiUrlZones}/>, document.getElementsByClassName('content')[0]);
})();
