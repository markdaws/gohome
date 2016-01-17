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
                        <li role="presentation">
                            <a href="#scenes" role="tab" aria-controls="scenes" data-toggle="tab">Scenes</a>
                        </li>
                        <li role="presentation">
                            <a href="#zones" role="tab" aria-controls="zones" data-toggle="tab">Zones</a>
                        </li>
                        {/*
                        <li role="presentation">
                            <a href="#devices" role="tab" aria-controls="devices" data-toggle="tab">Devices</a>
                            </li>*/}
                        <li role="presentation" className="active">
                            <a href="#logging" role="tab" aria-controls="logging" data-toggle="tab">Logging</a>
                        </li>
                        <li role="presentation">
                            <a href="#recipes" role="tab" aria-controls="recipes" data-toggle="tab">Recipes</a>
                        </li>
                    </ul>
                    <div className="tab-content">
                        <div role="tabpanel" className="tab-pane fade" id="scenes">
                            <SceneList scenes={this.state.scenes} />
                        </div>
                        <div role="tabpanel" className="tab-pane fade" id="zones">
                            <ZoneList zones={this.state.zones} />
                        </div>
                        <div role="tabpanel" className="tab-pane fade" id="devices">
                            <DeviceList devices={this.state.devices} />
                        </div>
                        <div role="tabpanel" className="tab-pane active" id="logging">
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
                <div className="cmp-Zone col-xs-12 col-sm-3 col-md-3 col-lg-3">
                    <a role="button" aria-expanded="false" aria-controls={"zoneControl" + this.props.id} data-toggle="collapse" href={"#zoneControl" + this.cssSafeIdentifier(this.props.id)} className="btn btn-primary zone">
                        <div>
                            <i className={icon}></i>
                        </div>
                        <span className="name">{this.props.name}</span>
                        <input style={{display: 'none'}} type="text" value={value} onChange={this.handleChange}></input>
                    </a>
                    {/* TODO: position:absolute if desktop/tablet vs phone */}
                    <ZoneControl ref="zoneControl" name={this.props.name} id={this.props.id} type={this.props.type} output={this.props.output}/>
                </div>
            )
        }
    });

    var ZoneControl = React.createClass({
        mixins: [CssMixin],
        componentDidMount: function() {
            var $el = $(ReactDOM.findDOMNode(this));
            var $value = $el.find('.level');
            var s = $el.find('.valueSlider');
            var slider = s.slider({ reversed: true});
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

            var left;
            if (this.props.output === 'continuous') {
                left = (
                    <div className="pull-left">
                        <h4 className="level">N/A</h4>
                        <input className="valueSlider" type="text" data-slider-value="0" data-slider-min="00" data-slider-max="100" data-slider-step="1" data-slider-orientation="vertical"></input>
                    </div>
                    );
            }

            var uniqueId = this.cssSafeIdentifier('zoneControl' + this.props.id);
            return (
                <div id={uniqueId} className={"cmp-ZoneControl collapse " + uniqueId}>
                    <div className="well">
                        <div className="content">
                            {left}
                            <div className="pull-right">
                                <a href="#" className="btn btn-default on" onClick={this.handleOnClick}>{onText}</a>
                                <a href="#" className="btn btn-default off" onClick={this.handleOffClick}>{offText}</a>
                            </div>
                           <div className="footer clearfix"></div>
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
                    console.log('set the scene');
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
            console.log(err);
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
                    console.log(val)
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
