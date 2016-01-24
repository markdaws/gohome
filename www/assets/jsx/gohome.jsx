//TODO: webpack, split up
(function() {

    var UniqueIdMixin = {
        getNextIdAndIncrement: function() {
            UniqueIdMixin._current += 1;
            return UniqueIdMixin._current;
        },

        getCurrentId: function() {
            return UniqueIdMixin._current;
        }
    };
    UniqueIdMixin._current = 0;

    var CssMixin = {
        cssSafeIdentifier: function(value) {
            return value.replace(/:/g, '_');
        }
    }

    var AssetsMixin = {
        getImageUrl: function(imageName) {
            return 'assets/images/' + imageName;
        }
    }

    var ControlApp = React.createClass({
        getInitialState: function() {
            return { scenes: [], zones: [], devices: [] };
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
                <div className="cmp-ControlApp">
                    <ul className="nav nav-tabs" role="tablist">
                        <li role="presentation" className="active">
                            <a href="#system" role="tab" aria-controls="system" data-toggle="tab">System</a>
                        </li>
                        <li role="presentation">
                            <a href="#scenes" role="tab" aria-controls="scenes" data-toggle="tab">Scenes</a>
                        </li>
                        <li role="presentation">
                            <a href="#zones" role="tab" aria-controls="zones" data-toggle="tab">Zones</a>
                        </li>
                        <li role="presentation">
                            <a href="#logging" role="tab" aria-controls="logging" data-toggle="tab">Logging</a>
                        </li>
                        <li role="presentation">
                            <a href="#recipes" role="tab" aria-controls="recipes" data-toggle="tab">Recipes</a>
                        </li>
                    </ul>
                    <div className="tab-content">
                        <div role="tabpanel" className="tab-pane active" id="system">
                            <System />
                        </div>
                        <div role="tabpanel" className="tab-pane fade" id="scenes">
                            <SceneList scenes={this.state.scenes} />
                        </div>
                        <div role="tabpanel" className="tab-pane fade" id="zones">
                            <ZoneList zones={this.state.zones} />
                        </div>
                        <div role="tabpanel" className="tab-pane fade" id="logging">
                            <Logging />
                        </div>
                        <div role="tabpanel" className="tab-pane fade" id="recipes">
                            <RecipeApp />
                        </div>
                    </div>
                </div>
            );
        }
    });

    var System = React.createClass({
        getInitialState: function() {
            return {
                importing: false,
            };
        },

        importProduct: function() {
            this.setState({ importing: true });
        },

        cancelImport: function() {
            this.setState({ importing: false });
        },
        
        render: function() {
            var body, importBtn
            if (this.state.importing) {
                body = <Import/>
                importBtn = <button className="btn btn-danger" onClick={this.cancelImport}>Cancel</button>
            } else {
                body = <SystemDeviceList/>
                importBtn = <button className="btn btn-primary" onClick={this.importProduct}>Import</button>
            }
            return (
                <div className="cmp-System">
                    {importBtn}
                    {body}
                </div>
            );
        }
    });

    var Import = React.createClass({
        getInitialState: function() {
            return { selectedProduct: null };
        },

        productSelected: function(evt) {
            this.setState({ selectedProduct: evt.target.value });
        },

        render: function() {

            var body
            switch(this.state.selectedProduct) {
            case 'TCP600GWB':
                body = <ImportTCP600GWB />
                break;
            case 'FluxWIFI':
                body = <ImportFluxWIFI />
                break;
            default:
                body = null;
            }
            return (
                <div className="cmp-Import">
                    <h4>Select a product to import</h4>
                    <select className="form-control" onChange={this.productSelected} value={this.state.selectedProduct}>
                        <option value="">Choose ...</option>
                        <option value="LLL">Lutron</option>
                        <option value="TCP600GWB">Connected By TCP Hub</option>
                        <option value="FluxWIFI">Flux Wifi</option>
                    </select>
                    <div className="content">
                        {body}
                    </div>
                </div>
            )
        }
    });

    var ImportFluxWIFI = React.createClass({
        getInitialState: function() {
            return {
                discovering: false,
                zones: [],
                loading: true,
                devices: [],
            };
        },

        componentDidMount: function() {
            var self = this;
            $.ajax({
                url: '/api/v1/systems/123/devices',
                dataType: 'json',
                cache: false,
                success: function(data) {
                    self.filterDevices(data || []);
                },
                error: function(xhr, status, err) {
                    console.error(err.toString());
                }
            });
        },

        filterDevices: function(devices) {
            var filteredDevices = [];
            for (var i=0; i<devices.length; ++i) {
                switch(devices[i].modelNumber) {
                    default:
//                case 'GoHomeHub':
                    filteredDevices.push(devices[i]);
                    break;
                }
            }

            this.setState({
                devices: filteredDevices,
                loading: false
            });
        },
        
        discover: function() {
            this.setState({
                discovering: true,
                zones: []
            });

            var self = this;
            $.ajax({
                url: '/api/v1/discovery/FluxWIFI/zones',
                dataType: 'json',
                cache: false,
                success: function(data) {
                    self.setState({
                        discovering: false,
                        zones: data
                    });
                },
                error: function(xhr, status, err) {
                    self.setState({
                        discovering: false
                    });
                    console.error(err);
                }
            });
        },

        render: function() {

            var loading
            if (this.state.loading) {
                loading = <div className="spinnerWrapper">
                    <i className="fa fa-spinner fa-spin"></i></div>
            }

            var noDeviceBody
            if (!this.state.loading && this.state.devices.length === 0) {
                noDeviceBody = (
                    <div>
                        <h3>Import failed</h3>
                        <p>In order to import Flux WIFI bulbs, you must have a device in your system
                        that is capable of controlling them.  Please add one of the following devices
                        to your system first, then come back and try to import again:
                        </p>
                        <ul>
                           <li>GoHomeHub</li>
                        </ul>
                    </div>
                );
            }
            
            var zones
            if (this.state.zones.length > 0) {
                var self = this
                zones = this.state.zones.map(function(zone) {
                    return <ZoneInfo devices={self.state.devices} zone={zone} key={zone.address} />
                })
            }

            var importBody
            if (!this.state.loading && this.state.devices.length > 0) {
                importBody = (
                    <div>
                    <button className={"btn btn-primary" + (this.state.discovering ? " disabled" : "")}
                        onClick={this.discover}>Discover Zones</button>
                    <i className={"fa fa-spinner fa-spin" + (this.state.discovering ? "" : " hidden")}></i>
                    <h3 className={this.state.zones.length > 0 ? "" : " hidden"}>Zones</h3>
                    {zones}
                    </div>
                );
            }
            return (
                <div className="cmp-ImportFluxWIFI">
                    {loading}
                    {noDeviceBody}
                    {importBody}
                </div>
            );
        }
    });

    var DevicePicker = React.createClass({
        getInitialState: function() {
            return {
                value: ''
            };
        },
        
        selected: function(evt) {
            this.setState({ value: evt.target.value });
        },
        
        render: function() {
            var options = [];
            this.props.devices.forEach(function(device) {
                console.log(device)
                options.push(<option key={device.id} value={device.id}>{device.name}</option>);
            });
            return (
                <div className="cmp-DevicePicker">
                    <select className="form-control" onChange={this.selected} value={this.state.value}>
                        <option value="">Select a device...</option>
                        {options}
                    </select>
                </div>
            );
        }
    });
    
    var ZoneInfo = React.createClass({
        mixins: [UniqueIdMixin],
        getInitialState: function() {
            return {
                zone: this.props.zone
            }
        },

        nameChanged: function(evt) {
            var zone = this.state.zone;
            zone.name = evt.target.value;
            this.setState({ zone : zone });
        },

        descriptionChanged: function(evt) {
            var zone = this.state.zone;
            zone.description = evt.target.value;
            this.setState({ zone : zone });
        },

        addressChanged: function(evt) {
            var zone = this.state.zone;
            zone.address = evt.target.value;
            this.setState({ zone : zone });
        },

        devicePickerChanged: function(device) {
            var zone = this.state.zone;
            zone.deviceId = device.id;
            this.setState({ zone: zone });
        },
        
        render: function() {
            //TODO unique names for ids

            var zone = this.state.zone
            return (
                <div className="cmp-ZoneInfo well">
                    <div className="form-group">
                        <label className="control-label" htmlFor={"name" + this.getNextIdAndIncrement()}>Name</label>
                        <input value={zone.name} onChange={this.nameChanged} className="name form-control" type="text" id={"name" + this.getCurrentId()}/>
                    </div>
                    <div className="form-group">
                        <label className="control-label" htmlFor={"description" + this.getNextIdAndIncrement()}>Description</label>
                        <input value={zone.description} onChange={this.descriptionChanged} className="description form-control" type="text" id={"description" + this.getCurrentId()}/>
                    </div>
                    <div className="form-group">
                        <label className="control-label" htmlFor={"address" + this.getNextIdAndIncrement()}>Address</label>
                        <input value={zone.address} onChange={this.addressChanged} className="address form-control" type="text" id={"address" + this.getCurrentId()}/>
                    </div>
                    <div className="form-group">
                        <label className="control-label" htmlFor={"device" + this.getNextIdAndIncrement()}>Device</label>
                        <DevicePicker devices={this.props.devices} change={this.devicePickerChanged}/>
                    </div>
                    <div className="clearfix">
                        <button className="btn btn-primary pull-left" onClick={this.turnOn}>Turn On</button>
                        <button className="btn btn-primary btnOff pull-left" onClick={this.turnOff}>Turn Off</button>
                    </div>
                </div>
            );
        }
    });

    var ImportTCP600GWB = React.createClass({
        getInitialState: function() {
            return {
                location: "",
                locationFailed: false,
                discoveryInProgress: false,
                tokenInProgress: false,
                token: '',
                tokenError: false,
                tokenMissingAddress: false
            };
        },
        
        autoDiscover: function() {
            var self = this;
            this.setState({ discoveryInProgress: true });
            
            $.ajax({
                url: '/api/v1/discovery/TCP600GWB',
                dataType: 'json',
                cache: false,
                success: function(data) {
                    self.setState({
                        location: data.location,
                        discoveryInProgress: false
                    });
                },
                error: function(xhr, status, err) {
                    self.setState({
                        locationFailed: true,
                        discoveryInProgress: false
                    });
                }
            });
        },

        getToken: function() {
            var device = this.refs.devInfo.getValues();
            this.setState({
                tokenMissingAddress: false,
                tokenInProgress: true
            });
            
            if (device.address === '') {
                this.setState({
                    tokenMissingAddress: true,
                    tokenInProgress: false
                });
                return;
            }
            
            var self = this;
            $.ajax({
                url: '/api/v1/discovery/TCP600GWB/token?address=' + device.address,
                dataType: 'json',
                cache: false,
                success: function(data) {
                    self.setState({
                        tokenInProgress: false,
                        token: data.token,
                        tokenError: data.unauthorized
                    });
                },
                error: function(xhr, status, err) {
                    self.setState({
                        tokenError: true,
                        tokenInProgress: false
                    });
                }
            });
        },
        
        render: function() {
            return (
                <div className="cmp-ImportTCP600GWB">
                    <p>Click to automatically retrieve the network address for this device</p>
                    <div className="form-group has-error">
                        <button className={"btn btn-primary" + (this.state.discoveryInProgress ? " disabled" : "")} onClick={this.autoDiscover}>Discover Address</button>
                        <i className={"fa fa-spinner fa-spin" + (this.state.discoveryInProgress ? "" : " hidden")}></i>
                        <span className={"help-block" + (this.state.locationFailed ? "" : " hidden")}>Error - Auto discovery failed, verify your TCP device is connected to the same network. If this continues to fail, use the official TCP app to get the device address</span>
                </div>
                <p>Click to retrive the security token. Only click this after pressing the "sync" button on your physical ConnectedByTCP hub</p>
                <div className="form-group has-error">
                    <button className={"btn btn-primary" + (this.state.tokenInProgress ? " disabled" : "")} onClick={this.getToken}>Get Token</button>
                    <i className={"fa fa-spinner fa-spin" + (this.state.tokenInProgress ? "" : " hidden")}></i>
                    <span className={"help-block" + (this.state.tokenError ? "" : " hidden")}>Error - unable to get the token, make sure you press the physical "sync" button on the TCP hub device before clicking the "Get Token" button otherwise this will fail</span>
                    <span className={"help-block" + (this.state.tokenMissingAddress ? "" : " hidden")}>Error - you must put a valid network address in the "Address" field first before clicking this button</span>
                </div>
                <DeviceInfo showToken="true" token={this.state.token} tokenError={this.state.tokenError} address={this.state.location} ref="devInfo"/>
                </div>
            )
        }
    });

    var DeviceInfo = React.createClass({
        getInitialState: function() {
            return {
                device: {
                    name: this.props.name || '',
                    description: this.props.description || '',
                    address: this.props.address,
                    id: '',
                    modelNumber: '',
                    securityToken: this.props.token,
                    showToken: false
                }
            }
        },

        getValues: function() {
            return this.state.device;
        },

        componentWillReceiveProps: function(nextProps) {
            var device = this.state.device;
            if (nextProps.name != "") {
                device.name = nextProps.name;
            }
            if (nextProps.description != "") {
                device.description = nextProps.description;
            }
            if (nextProps.address != "") {
                device.address = nextProps.address;
            }
            if (nextProps.token != "") {
                device.securityToken = nextProps.token;
            }
            this.setState({ device: device });
        },

        nameChanged: function(evt) {
            var device = this.state.device;
            device.name = evt.target.value;
            this.setState({ device: device });
        },

        descriptionChanged: function(evt) {
            var device = this.state.device;
            device.description = evt.target.value;
            this.setState({ device: device });
        },

        addressChanged: function(evt) {
            var device = this.state.device;
            device.address = evt.target.value;
            this.setState({ device: device });
        },

        tokenChanged: function(evt) {
            var device = this.state.device;
            device.securityToken = evt.target.value;
            this.setState({ device: device });
        },

        testConnection: function() {
            //TODO: How to know what to call
        },
        
        render: function() {
            //TODO:need unique name for id and htmlFor
            var device = this.state.device;

            var token
            if (this.props.showToken) {
                token = (
                    <div className={"form-group" + (this.props.tokenError ? " has-error" : "")}>
                        <label className="control-label" htmlFor="securitytoken">Security Token</label>
                        <input value={device.securityToken} onChange={this.tokenChanged} className="securitytoken form-control" type="text" id="securitytoken"/>
                        <span className={"help-block" + (this.props.tokenError ? "" : " hidden")}>Error - failed to fetch token, make sure you pressed the sync button on the tcp hub device before requesting the token</span>
                    </div>
                );
            }
            
            return (
                <div className="cmp-DeviceInfo well">
                    <div className="form-group">
                        <label className="control-label" htmlFor="name">Name</label>
                        <input value={device.name} onChange={this.nameChanged} className="name form-control" type="text" id="name"/>
                        <span className={"help-block hidden"}>Error - TODO:</span>
                    </div>
                    <div className="form-group">
                        <label className="control-label" htmlFor="description">Description</label>
                        <input value={device.description} onChange={this.descriptionChanged} className="description form-control" type="text" id="description"/>
                        <span className={"help-block hidden"}>Error - TODO:</span>
                    </div>
                    <div className="form-group">
                        <label className="control-label" htmlFor="address">Address</label>
                        <input value={device.address} onChange={this.addressChanged} className="address form-control" type="text" id="address"/>
                        <span className={"help-block hidden"}>Error - TODO:</span>
                </div>
                {token}
                <button className="btn btn-primary" onClick={this.testConnection}>Test Connection</button>
                
                </div>
            );
        }
    });
    
    var SystemDeviceList = React.createClass({
        getInitialState: function() {
            return {
                loading: true,
                devices: [],
                addingNew: false
            };
        },

        componentDidMount: function() {
            var self = this;
            $.ajax({
                url: '/api/v1/systems/123/devices',
                dataType: 'json',
                cache: false,
                success: function(data) {
                    self.setState({devices: data, loading: false});
                },
                error: function(xhr, status, err) {
                    console.error(err.toString());
                }
            });
        },

        newClicked: function() {
            //TODO: Show new device UI
        },

        render: function() {
            var deviceNodes = this.state.devices.map(function(device) {
                return (
                    <DeviceInfo
                    name={device.name}
                    description={device.description}
                    address={device.address}
                    key={device.id}
                    />
                );
            })
            
            var body = this.state.loading
                ? <div className="text-center"><i className="fa fa-spinner fa-spin"></i></div>
                : deviceNodes;

            return (
                <div className="cmp-DeviceList">
                    <div className="header clearfix">
                        <button className="btn btn-primary pull-right" onClick={this.newClicked}>New Device</button>
                    </div>
                    <h3 className={this.state.devices.length > 0 ? "" : " hidden"}>Devices</h3>
                    {body}
                </div>
            );
        }
    });
    
    var Logging = React.createClass({
        getInitialState: function() {
            return {
                items: [],
                connectionStatus: 'connecting'
            };
        },

        componentDidMount: function() {
            this.reconnect();
        },

        componentDidUpdate: function() {
            var lastLi = this.refs.lastLi;
            if (!lastLi) {
                return;
            }

            //TODO: Shouldn't set the body element like this, use events
            //TODO: If the user has scrolled away from the bottom, don't do this
            //until they scroll back to the bottom again, annoying to jump away
            $('body')[0].scrollTop = ReactDOM.findDOMNode(lastLi).offsetTop;
        },

        componentWillUnmount: function() {
            var conn = this.state.conn;
            if (!conn) {
                return;
            }
            conn.Close();
        },

        reconnect: function() {
            var oldConn = this.state.conn;
            if (oldConn) {
                oldConn.close();
            }

            var conn = new WebSocket("ws://" + window.location.host + "/api/v1/events/ws");
            var self = this;
            conn.onopen = function(evt) {
                self.setState({
                    connectionStatus: 'connected'
                });
            };
            conn.onclose = function(evt) {
                conn = null;
                self.setState({
                    conn: null,
                    items: [],
                    connectionStatus: 'disconnected'
                });
            };
            conn.onmessage = function(evt) {
                var item = JSON.parse(evt.data);
                item.datetime = new Date(item.datetime);
                self.setState({ items: self.state.items.concat(item)});
            };
            this.setState({
                conn: conn,
                connectionStatus: 'connecting'
            });

            //TODO: Fetch X previous log items from server?
        },

        clearClicked: function() {
            this.setState({ items: [] });
        },

        render: function() {
            var body;

            switch(this.state.connectionStatus) {
            case 'connected':
                var itemCount = this.state.items.length;
                body = this.state.items.map(function(item, i) {
                    return <LogLine item={item} key={item.id} ref={itemCount === i+1 ? 'lastLi' : undefined}/>;
                });
                break;

            case 'connecting':
                body = <li className="spinner"><i className="fa fa-spinner fa-spin"></i></li>
                break;

            case 'disconnected':
                body = <li className="reconnect"><button className="btn btn-primary" onClick={this.reconnect}>Reconnect</button></li>
                break;
            }

            var hasEvents = this.state.items.length > 0;
            var waiting = !hasEvents && this.state.connectionStatus === 'connected';
            return (
                <div className="cmp-Logging">
                    <h3 className={!waiting ? 'hidden' : ''}>Waiting for events...</h3>
                    <ol className="list-unstyled">
                        {body}
                    </ol>
                    <div className="footer text-center">
                        <button className={(hasEvents ? '' : 'hidden') + ' btn btn-default'} onClick={this.clearClicked}>Clear</button>
                    </div>
                </div>
            );
        }
    });

    var LogLine = React.createClass({
        render: function() {
            return (
                <li className="cmp-LogLine">
                    <span className="datetime">{this.props.item.datetime.toLocaleString()}</span>
                    <span className="deviceName"> [{this.props.item.deviceName}]</span>
                    <span> : {this.props.item.friendlyMessage}</span>
                    <span className="rawMessage"> [Raw: {this.props.item.rawMessage}]</span>
                </li>
            );
        }
    });

    var DeviceList = React.createClass({
        render: function() {
            return (
                <div className="cmp-DeviceList">DEVICES!</div>
            );
        }
    });

    var ZoneList = React.createClass({
        render: function() {
            //TODO: Add loading
            var self = this;
            var zoneNodes = Object.keys(this.props.zones).map(function(id) {
                var zone = self.props.zones[id];
                return (
                    <Zone id={zone.id} name={zone.name} type={zone.type} output={zone.output} key={id}/>
                );
            })
            return (
                <div className="cmp-ZoneList row">
                    {zoneNodes}
                </div>
            );
        }
    });

    var Zone = React.createClass({
        mixins: [CssMixin],
        getInitialState: function() {
            return {
                value: 0,
                showSlider: false,
                slider: null
            }
        },

        componentDidMount: function() {
            var self = this;
            
            switch (this.props.output) {
            case 'binary':
            case 'continuous':
                var s = $(ReactDOM.findDOMNode(this)).find('.valueSlider');
                s.slider({ reversed: false });
                self.setState({ slider: s });
                s.on('change', function(evt) {
                    self.setState({ value: evt.value.newValue });
                });
                s.on('slideStop', function(evt) {
                    self.setValue('setLevel', evt.value, 0, 0, 0, function(err) {
                        if (err) {
                            //TODO:
                            console.error(err);
                        }
                    });
                    return false;
                });
                break;

            case 'rgb':
                var $el = $(ReactDOM.findDOMNode(this)).find('.zone-rgb .clickInfo span')
                $el.colorPicker({
                    doRender:false,
                    opacity: false,
                    margin: '0px 0px 0px -30px',
                    renderCallback: function($e, toggled) {
                        if (toggled !== undefined) {
                            // only send a value when the user actually interacts with the
                            // control not when it is first shown/hidden
                            return;
                        }
                        var rgb = this.color.colors.rgb;
                        self.setValue(
                            'setLevel',
                            0,
                            parseInt(rgb.r * 255),
                            parseInt(rgb.g * 255),
                            parseInt(rgb.b * 255),
                            function(err) {
                                if (err) {
                                    console.error(err);
                                }
                            }
                        );
                    }
                });
                break;
            }
        },

        infoClicked: function(evt) {
            evt.stopPropagation();
            evt.preventDefault();

            if (!this.isRgb()) {
                this.setState({ showSlider: true });
            }
        },

        isRgb: function() {
            return this.props.output === 'rgb';
        },
        
        setValue: function(cmd, value, r, g, b, callback) {
            if (!this.isRgb()) {
                this.state.slider.slider('setValue', value, false, true);
            }
            //TODO: Need rgb
            this.setState({ value: value });
            this.send({
                cmd: cmd,
                value: parseFloat(value),
                r: r,
                g: g,
                b: b
            }, callback);
        },

        toggleOn: function(evt) {
            evt.stopPropagation();
            evt.preventDefault();

            var cmd, level;
            if (this.state.value !== 0) {
                cmd = 'turnOff';
                level = 0;
            } else {
                cmd = 'turnOn';
                level = 100;
            }
            this.setValue(cmd, level, 0, 0, 0, function(err) {
                if (err) {
                    console.error(err);
                }
            });
        },

        send: function(data, callback) {
            $.ajax({
                url: '/api/v1/systems/1/zones/' + this.props.id,
                type: 'POST',
                dataType: 'json',
                contentType: 'application/json; charset=utf-8',
                data: JSON.stringify(data),
                success: function(data) {
                    callback();
                },
                error: function(xhr, status, err) {
                    callback(err);
                }
            });
        },

        render: function() {
            var value = this.state.value;
            var icon = this.props.type === 'light' ? 'fa fa-lightbulb-o' : 'fa fa-picture-o';

            var stepSize
            switch (this.props.output) {
            case 'continuous':
                stepSize = 1;
                break;
            case 'binary':
                stepSize = 100;
                break;
            case 'rgb':
                break;
            default:
                stepSize = 1;
            }
            
            return (
                <div className="cmp-Zone col-xs-12 col-sm-4 col-md-4 col-lg-4 clearfix">
                    <button className={"btn btn-primary zone" + (this.isRgb() ? " zone-rgb" : "")}>
                        <i className={icon}></i>
                        <span className="name">{this.props.name}</span>
                        <div className={"sliderWrapper pull-right" + (this.state.showSlider ? "" : " hidden")} >
                            <input className="valueSlider" type="text" data-slider-value="0" data-slider-min="00" data-slider-max="100" data-slider-step={stepSize} data-slider-orientation="horizontal"></input>
                            <span className="level pull-right">{this.state.value}%</span>
                        </div>
                        <div className="clearfix footer">
                            <div className={"clickInfo pull-right" + (this.state.showSlider ? " hidden" : "")}>
                                <span onClick={this.infoClicked}>Click to control</span>
                            </div>
                            <a className="btn btn-link pull-left" onClick={this.toggleOn}>
                                <i className="fa fa-power-off"></i>
                            </a>
                        </div>
                    </button>
                </div>
            )
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

    var Scene = React.createClass({
        handleClick: function(event) {
            $.ajax({
                url: '/api/v1/systems/1/scenes/active',
                type: 'POST',
                dataType: 'json',
                contentType: 'application/json; charset=utf-8',
                data: JSON.stringify({ id: this.props.scene.id }),
                success: function(data) {
                    //TODO: Common way in UI to display success/error
                }.bind(this),
                error: function(xhr, status, err) {
                    console.error(err.toString());
                }.bind(this)
            });
        },

        render: function() {
            return (
                <div className="cmp-Scene col-xs-6 col-sm-3 col-md-3 col-lg-3">
                    <a className="btn btn-primary scene" onClick={this.handleClick}>
                        <div>
                            <i className="fa fa-sliders"></i>
                        </div>
                        <span className="name">{this.props.scene.name}</span>
                    </a>
                </div>
            )
        }
    });

    var RecipeApp = React.createClass({
        getInitialState: function() {
            return {
                cookBooks: [],
                recipes: [],
                creatingRecipe: false
            }
        },

        recipeCreated: function(recipe) {
            this.setState({ creatingRecipe: false });
        },

        recipeCancel: function() {
            this.setState({ creatingRecipe: false });
        },

        componentDidMount: function() {
            $.ajax({
                url: '/api/v1/cookbooks',
                dataType: 'json',
                cache: false,
                success: function(data) {
                    this.setState({cookBooks: data});
                }.bind(this),
                error: function(xhr, status, err) {
                    console.error(err.toString());
                }.bind(this)
            });
        },

        newClicked: function() {
            this.setState({ creatingRecipe: true })
        },

        render: function() {
            var body, newButton;
            if (this.state.creatingRecipe) {
                body = <NewRecipe cookBooks={this.state.cookBooks} onCancel={this.recipeCancel} onCreate={this.recipeCreated}/>
            } else {
                newButton = <button className="btn btn-primary pull-right" onClick={this.newClicked}>New</button>
                body = <RecipeList />
            }

            return (
                <div className="cmp-RecipeApp">
                    <div className="clearfix header">
                        {newButton}
                    </div>
                    {body}
                </div>
            );
        }
    });

    var RecipeList = React.createClass({
        getInitialState: function() {
            return {
                recipes: [],
                loading: true
            };
        },

        addRecipe: function(recipe) {
            this.state.recipes.push(recipe);
            this.setState({ recipes: this.state.recipes });
        },

        componentDidMount: function() {
            var self = this;

            $.ajax({
                url: '/api/v1/recipes',
                dataType: 'json',
                cache: false,
                success: function(data) {
                    setTimeout(function() {
                        self.setState({recipes: data, loading: false});
                    }, 250);
                },
                error: function(xhr, status, err) {
                    console.error(err.toString());
                }
            });
        },

        recipeDestroyed: function(recipeId) {
            var recipes = this.state.recipes;
            for (var i=0; i<recipes.length; ++i) {
                if (recipes[i].id === recipeId) {
                    recipes.splice(i, 1);
                    this.setState({ recipes, recipes });
                    break;
                }
            }
        },

        render: function() {
            var self = this;
            var recipeNodes = this.state.recipes.map(function(recipe) {
                return (
                    <RecipeInfo recipe={recipe} key={recipe.id} onDestroy={self.recipeDestroyed}/>
                );
            });

            var body = this.state.loading
            ? <div className="text-center"><i className="fa fa-spinner fa-spin"></i></div>
            : recipeNodes;

            return (
                <div className="cmp-RecipeList">{body}</div>
            );
        }
    });

    var RecipeInfo = React.createClass({
        getInitialState: function() {
            return {
                enabled: this.props.recipe.enabled,
                checkboxDisabled: false
            };
        },

        deleteClicked: function(evt) {
            var self = this;
            $.ajax({
                url: '/api/v1/recipes/' + this.props.recipe.id,
                type: 'DELETE',
                cache: false,
                success: function(data) {
                    self.props.onDestroy(self.props.recipe.id);
                }.bind(this),
                error: function(xhr, status, err) {
                    console.error(err);
                }.bind(this)
            });
        },

        checkboxChange: function(evt) {
            var checkedState = evt.target.checked;

            this.setState({
                enabled: checkedState,
                checkboxDisabled: true
            });

            var self = this;
            $.ajax({
                url: '/api/v1/recipes/' + this.props.recipe.id,
                type: 'POST',
                dataType: 'json',
                data: JSON.stringify({ enabled: checkedState }),
                cache: false,
                success: function(data) {
                    self.setState({ checkboxDisabled: false })
                }.bind(this),
                error: function(xhr, status, err) {
                    console.error(err);
                    self.setState({
                        checkboxDisabled: false,
                        enabled: !checkedState
                    });
                }.bind(this)
            });
        },

        render: function() {
            var recipe = this.props.recipe;

            var inputTitle = this.state.enabled
                ? 'Click to disable'
                : 'Click to enabled';

            return (
                <div className="cmp-RecipeInfo well">
                    <h4>{recipe.name}</h4>
                    <p>{recipe.description}</p>
                    <div className="clearfix">
                        <input type="checkbox" title={inputTitle} className={this.state.checkboxDisabled ? 'disabled' : 'checkbox'} checked={this.state.enabled} onChange={this.checkboxChange}/>
                        <button className="btn btn-danger pull-right" onClick={this.deleteClicked} >Delete</button>
                    </div>
                </div>
            )
        }
    });

    var NewRecipe = React.createClass({
        getInitialState: function() {
            return {
                triggerCookBookID: -1,
                actionCookBookID: -1,
                triggers: null,
                actions: null,
                trigger: null,
                action: null,
                name: '',
                description: '',
                saveError: null,
                saving: false
            };
        },

        triggerCookBookSelected: function(cookBookID) {
            this.setState({ triggerCookBookID: cookBookID });

            var self = this;
            this.loadCookBook(cookBookID, function(err, data) {
                if (err) {
                    console.error(err.toString());
                    return;
                }

                self.setState({ triggers: data.triggers });
            });
        },

        actionCookBookSelected: function(cookBookID) {
            this.setState({ actionCookBookID: cookBookID });

            var self = this;
            this.loadCookBook(cookBookID, function(err, data) {
                if (err) {
                    console.error(err.toString());
                    return;
                }

                self.setState({ actions: data.actions });
            });
        },

        loadCookBook: function(cookBookID, callback) {
            $.ajax({
                url: '/api/v1/cookbooks/' + cookBookID,
                dataType: 'json',
                cache: false,
                success: function(data) {
                    callback(null, data);
                },
                error: function(xhr, status, err) {
                    callback({ err: err });
                }
            });
        },

        triggerSelected: function(trigger) {
            this.setState({ trigger: trigger });
        },

        actionSelected: function(action) {
            this.setState({ action: action });
        },

        saveClicked: function(evt) {
            this.setState({ saveError: null, saving: true });

            var recipe = this.toJSON();
            var self = this;
            $.ajax({
                url: '/api/v1/recipes',
                type: 'POST',
                dataType: 'json',
                data: JSON.stringify(recipe),
                cache: false,
                success: function(data) {
                    self.setState({ saving: false });
                    self.props.onCreate(recipe);
                },
                error: function(xhr, status, err) {
                    self.setState({ saving: false });
                    if (xhr.status === 400) {
                        self.setState({ saveError: JSON.parse(xhr.responseText) });
                    } else {
                        //Unknown error - todo
                    }
                }
            });
        },

        handleNameChange: function(evt) {
            this.setState({ name: evt.target.value });
        },

        handleDescriptionChange: function(evt) {
            this.setState({ description: evt.target.value });
        },

        toJSON: function() {
            var json = {};
            json.name = this.state.name;
            json.description = this.state.description;

            if (this.state.action) {
                json.action = {
                    id: this.state.action.id,
                    ingredients: this.refs.actionIngredients.toJSON()
                }
            }

            if (this.state.trigger) {
                json.trigger = {
                    id: this.state.trigger.id,
                    ingredients: this.refs.triggerIngredients.toJSON()
                }
            }
            return json;
        },

        cancelClicked: function(evt) {
            this.props.onCancel();
        },

        render: function() {
            var nameErr = false;
            var descErr = false;
            var triggerErr = false;
            var actionErr = false;
            var triggerIngredientErr;
            var actionIngredientErr;
            var err = this.state.saveError;
            var errDesc = '';
            if (err) {
                switch (err.paramId) {
                case 'name':
                    nameErr = true;
                    errDesc = err.description;
                    break;
                case 'description':
                    descErr = true;
                    errDesc = err.description;
                    break;
                case 'trigger':
                    triggerErr = true;
                    errDesc = err.description;
                    break;
                case 'action':
                    actionErr = true;
                    errDesc = err.description;
                    break;
                default:
                    if (err.paramId.startsWith('trigger.')) {
                        triggerIngredientErr = err;
                        triggerIngredientErr.paramId = triggerIngredientErr.paramId.replace('trigger.', '');
                    } else if (err.paramId.startsWith('action.')) {
                        actionIngredientErr = err;
                        actionIngredientErr.paramId = actionIngredientErr.paramId.replace('action.', '');
                    }
                }
            }

            var triggerChild, actionChild;
            var spinner = <div className="text-center"><i className="fa fa-spinner fa-spin"></i></div>;
            if (this.state.trigger) {
                // Render the selected trigger
                triggerChild = <IngredientList err={triggerIngredientErr} ref="triggerIngredients" ingredients={this.state.trigger.ingredients} />
            } else if (this.state.triggers) {
                // Render the trigger list
                triggerChild = <TriggerList triggers={this.state.triggers} selected={this.triggerSelected}/>
            }
            else if (this.state.triggerCookBookID !== -1) {
                // Got a cook book, must be loading triggers
                triggerChild = spinner;
            }
            else {
                //TODO: Only show cook books that have triggers
                triggerChild = <CookBookList cookBooks={this.props.cookBooks} selected={this.triggerCookBookSelected}/>
            }

            if (this.state.action) {
                actionChild = <IngredientList err={actionIngredientErr} ref="actionIngredients" ingredients={this.state.action.ingredients} />
            } else if (this.state.actions) {
                actionChild = <ActionList actions={this.state.actions} selected={this.actionSelected}/>
            }
            else if (this.state.actionCookBookID !== -1) {
                actionChild = spinner;
            }
            else {
                actionChild = <CookBookList cookBooks={this.props.cookBooks} selected={this.actionCookBookSelected}/>
            }

            return (
                <div className="cmp-NewRecipe">
                    <div className={"form-group" + (nameErr ? " has-error" : "")}>
                        <label className="control-label" htmlFor="name">Name</label>
                        <input value={this.state.name} onChange={this.handleNameChange} className="name form-control" type="text" id="name"/>
                        <span className={"help-block" + (nameErr ? "" : " invisible")}>Error - {errDesc}</span>
                    </div>
                    <div className={"form-group" + (descErr ? " has-error" : "")}>
                        <label className="control-label" htmlFor="description">Description</label>
                        <input value={this.state.description} onChange={this.handleDescriptionChange} className="description form-control" type="text" id="description"/>
                        <span className={"help-block" + (descErr ? "" : " invisible")}>Error - {errDesc}</span>
                    </div>
                    <div className={"trigger form-group" + (triggerErr ? " has-error" : "")}>
                        <h3>Trigger</h3>
                        {triggerChild}
                        <span className={"help-block" + (triggerErr ? "" : " invisible")}>Error - {errDesc}</span>
                    </div>
                    <div className={"action form-group" + (actionErr ? " has-error" : "")}>
                        <h3>Action</h3>
                        {actionChild}
                        <span className={"help-block" + (actionErr ? "" : " invisible")}>Error - {errDesc}</span>
                    </div>
                    <div className="clearfix footer">
                        <button className={"btn btn-default pull-right" + (this.state.saving ? " disabled" : "")} onClick={this.cancelClicked}>Cancel</button>
                        <button className={"btn btn-primary pull-right" + (this.state.saving ? " disabled" : "")} onClick={this.saveClicked}>Save</button>
                    </div>
                </div>
            );
        }
    });

    var CookBookList = React.createClass({
        handleClick: function(cookBookID) {
            this.props.selected(cookBookID);
        },

        render: function() {
            var self = this;
            var cookBookNodes = this.props.cookBooks.map(function(cookBook) {
                return (
                    <CookBook data={cookBook} selected={self.handleClick} key={cookBook.id}/>
                );
            });
            return (
                <div className="cmp-CookBookList clearfix">
                    {cookBookNodes}
                </div>
            );
        }
    });

    var CookBook = React.createClass({
        mixins: [AssetsMixin],
        handleClick: function(evt) {
            evt.preventDefault();
            evt.stopPropagation();
            this.props.selected(this.props.data.id);
        },

        render: function() {
            return (
                <div className="cmp-CookBook">
                    <button className="btn btn-default" onClick={this.handleClick}>
                        <img src={this.getImageUrl(this.props.data.logoUrl)} />
                        {this.props.data.name}
                    </button>
                </div>
            );
        }
    });

    var TriggerList = React.createClass({
        handleClick: function(trigger) {
            this.props.selected(trigger);
        },

        render: function() {
            var self = this;
            var triggerNodes = this.props.triggers.map(function(trigger) {
                return (
                    <Trigger data={trigger} selected={self.handleClick} key={trigger.name} />
                );
            });

            return (
                <div className="cmp-TriggerList clearfix">
                    {triggerNodes}
                </div>
            );
        }
    });

    var Trigger = React.createClass({
        handleClick: function(evt) {
            evt.preventDefault();
            evt.stopPropagation();
            this.props.selected(this.props.data);
        },

        render: function() {
            return (
                <div className="cmp-Trigger pull-left">
                    <button className="btn btn-primary" onClick={this.handleClick}>
                        <h4>{this.props.data.name}</h4>
                        <p>{this.props.data.description}</p>
                    </button>
                </div>
            );
        }
    });

    var ActionList = React.createClass({
        handleClick: function(action) {
            this.props.selected(action);
        },

        render: function() {
            var self = this;
            var actionNodes = this.props.actions.map(function(action) {
                return (
                    <Action data={action} selected={self.handleClick} key={action.name}/>
                );
            });
            return (
                <div className="cmp-ActionList clearfix">
                    {actionNodes}
                </div>
            );
        }
    });

    var Action = React.createClass({
        handleClick: function(evt) {
            evt.preventDefault();
            evt.stopPropagation();
            this.props.selected(this.props.data);
        },

        render: function() {
            return (
                <div className="cmp-Trigger pull-left">
                    <button className="btn btn-primary" onClick={this.handleClick}>
                        <h4>{this.props.data.name}</h4>
                        <p>{this.props.data.description}</p>
                    </button>
                </div>
            );
        }
    });

    var IngredientList = React.createClass({
        render: function() {
            var self = this;
            var ingredientNodes = this.props.ingredients.map(function(ingredient) {
                var err;
                if (self.props.err && self.props.err.paramId === ingredient.id) {
                    err = self.props.err;
                }
                return (
                    <Ingredient err={err} data={ingredient} ref={ingredient.id} key={ingredient.id} />
                );
            });

            return (
                <div className="cmp-IngredientList well">
                    {ingredientNodes}
                </div>
            );
        },

        toJSON: function() {
            var json = {};
            var self = this;
            Object.keys(this.refs).map(function(key) {
                var val = self.refs[key].value();
                if (val != undefined) {
                    json[key] = val;
                }
            });
            return json;
        }
    });

    var Ingredient = React.createClass({
        mixins: [UniqueIdMixin],
        getInitialState: function() {
            return {
                value: undefined
            };
        },

        changeHandler: function(evt) {
            this.setState({ value: evt.target.value });
        },

        render: function() {
            var input;
            switch(this.props.data.type) {
            case 'string':
            case 'duration':
            case 'integer':
            case 'float':
                input = <input className="ingredientInput form-control" type="text" onChange={this.changeHandler} id={this.getNextIdAndIncrement()}/>;
                break;
            case 'boolean':
                input = <input className="ingredientInput form-control" type="checkbox" value="true" onChange={this.changeHandler} id={this.getNextIdAndIncrement()}/>;
                break;
            case 'datetime':
                //TODO: show calendar
                break;
            default:
                throw 'unknown ingredient type: ' + this.props.data.type;
            }

            var err = this.props.err;
            var errDesc = err ? err.description : '';
            return (
                <div>
                    <div className={"form-group" + (err ? " has-error" : "")}>
                        <label className="control-label" htmlFor={this.getCurrentId()}>{this.props.data.name}</label>
                        <p>{this.props.data.description}</p>
                        {input}
                        <span className={"help-block" + (err ? "" : " invisible")}>Error - {errDesc}</span>
                    </div>
                </div>
            );
        },

        value: function() {
            if (this.state.value == undefined) {
                return undefined;
            }
            
            switch(this.props.data.type) {
            case 'string':
                return this.state.value;
            case 'integer':
            case 'duration':
                return parseInt(this.state.value, 10);
            case 'float':
                return parseFloat(this.state.value)
            case 'boolean':
                return this.state.value === true || this.state.value === 'true';
                break;
            case 'datetime':
                //TODO:
                break;
            default:
                throw 'Unknown data type: ' + this.props.data.type;
            }
        }
    });

    var apiUrl = '/api/v1/systems/123/scenes';
    var apiUrlZones = '/api/v1/systems/123/zones';
    ReactDOM.render(<ControlApp url={apiUrl} zoneUrl={apiUrlZones}/>, document.getElementsByClassName('content')[0]);
})();
