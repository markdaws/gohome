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
            var icon = this.props.type === 'light' ? 'fa fa-lightbulb-o' : 'fa fa-picture-o';
            return (
                <div className="col-xs-12 col-sm-3 col-md-3 col-lg-3">
                    <a data-toggle="collapse" href={".zoneControl" + this.props.id} className="btn btn-default zone" onClick={this.handleClick}>
                        <div>
                            <i className={icon}></i>
                        </div>
                        <span className="name">{this.props.name} ({this.props.type})</span>
                        <input style={{display: 'none'}} type="text" value={value} onChange={this.handleChange}></input>
                    </a>
                    <ZoneControl ref="zoneControl" name={this.props.name} id={this.props.id} type={this.props.type}/>
                </div>
            )
        }
    });

    var ZoneControl = React.createClass({
        componentDidMount: function() {
            var $el = $(ReactDOM.findDOMNode(this));
            var $value = $el.find('.level');
            var s = $el.find('.valueSlider');
            var slider = s.slider({ reversed: true});
            console.log(slider);
            this.setState({ slider: slider });
            var self = this;
            s.on('change', function(evt) {
                $value.text(evt.value.newValue + '%');
            });
            s.on('slideStop', function(evt) {
                self._setValue(evt.value, function(err) {
                    if (err) {
                        console.error(err);
                    }
                });
            });
        },

        _setValue: function(value, callback) {
            var $el = $(ReactDOM.findDOMNode(this));
            this.state.slider.slider('setValue', value, false, true);
            this._send({ value: parseFloat(value) }, callback);
        },

        _send: function(data, callback) {
            $.ajax({
                url: '/api/v1/systems/1/zones/' + this.props.id,
                type: 'POST',
                dataType: 'json',
                contentType: 'application/json; charset=utf-8',
                data: JSON.stringify(data),
                success: function(data) {
                    callback();
                }.bind(this),
                error: function(xhr, status, err) {
                    callback({ err: err });
                }.bind(this)
            });
        },

        handleOnClick: function(evt) {
            evt.stopPropagation();
            evt.preventDefault();
            this._setValue(100, function(err) {
                if (err) {
                    console.error(err);
                }
            });
        },

        handleOffClick: function(evt) {
            evt.stopPropagation();
            evt.preventDefault();
            this._setValue(0, function(err) {
                if (err) {
                    console.error(err);
                }
            });
        },

        render: function() {
            var onText = 'On';
            var offText = 'Off';
            if (this.props.type !== 'light') {
                onText = 'Open';
                offText = 'Close';
            }
            return (
                <div className={"collapse zoneControl " + " zoneControl" + this.props.id}>
                    <div className="well">
                        <div className="content">
                            <div className="left">
                                <h4 className="level">N/A</h4>
                                <input className="valueSlider" type="text" data-slider-value="0" data-slider-min="00" data-slider-max="100" data-slider-step="1" data-slider-orientation="vertical"></input>
                            </div>
                            <div className="right">
                                <a href="#" className="btn btn-default on" onClick={this.handleOnClick}>{onText}</a>
                                <a href="#" className="btn btn-default off" onClick={this.handleOffClick}>{offText}</a>
                            </div>
                           <div className="footer"></div>
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
