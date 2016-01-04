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
                <div className="cmp-ControlApp">
                    <ul className="nav nav-tabs" role="tablist">
                        <li role="presentation" className="active">
                            <a href="#scenes" role="tab" aria-controls="scenes" data-toggle="tab">Scenes</a>
                        </li>
                        <li role="presentation">
                            <a href="#zones" role="tab" aria-controls="zones" data-toggle="tab">Zones</a>
                        </li>
                        <li role="presentation">
                            <a href="#recipes" role="tab" aria-controls="recipes" data-toggle="tab">Recipes</a>
                        </li>
                    </ul>
                    <div className="tab-content">
                        <div role="tabpanel" className="tab-pane active" id="scenes">
                            <SceneList scenes={this.state.scenes} />
                        </div>
                        <div role="tabpanel" className="tab-pane fade" id="zones">
                            <ZoneList zones={this.state.zones} />
                        </div>
                        <div role="tabpanel" className="tab-pane fade" id="recipes">
                            <RecipeApp />
                        </div>
                    </div>
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
                    <Zone id={zone.id} name={zone.name} type={zone.type} key={id}/>
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
                    <ZoneControl ref="zoneControl" name={this.props.name} id={this.props.id} type={this.props.type}/>
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

            var uniqueId = this.cssSafeIdentifier('zoneControl' + this.props.id);
            return (
                <div id={uniqueId} className={"cmp-ZoneControl collapse " + uniqueId}>
                    <div className="well">
                        <div className="content">
                            <div className="pull-left">
                                <h4 className="level">N/A</h4>
                                <input className="valueSlider" type="text" data-slider-value="0" data-slider-min="00" data-slider-max="100" data-slider-step="1" data-slider-orientation="vertical"></input>
                            </div>
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
                description: ''
            };
        },

        triggerCookBookSelected: function(cookBookID) {
            this.setState({ triggerCookBookID: cookBookID });

            var self = this;
            this._loadCookBook(cookBookID, function(err, data) {
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
            this._loadCookBook(cookBookID, function(err, data) {
                if (err) {
                    console.error(err.toString());
                    return;
                }

                self.setState({ actions: data.actions });
            });
        },

        _loadCookBook: function(cookBookID, callback) {
            $.ajax({
                url: '/api/v1/cookbooks/' + cookBookID,
                dataType: 'json',
                cache: false,
                success: function(data) {
                    callback(null, data);
                }.bind(this),
                error: function(xhr, status, err) {
                    callback({ err: err });
                }.bind(this)
            });
        },

        triggerSelected: function(trigger) {
            this.setState({ trigger: trigger });
        },

        actionSelected: function(action) {
            this.setState({ action: action });
        },

        saveClicked: function(evt) {
            evt.preventDefault();
            evt.stopPropagation();

            var recipe = this.toJSON();
            var self = this;
            $.ajax({
                url: '/api/v1/recipes',
                type: 'POST',
                dataType: 'json',
                data: JSON.stringify(recipe),
                cache: false,
                success: function(data) {
                    self.props.onCreate(recipe);
                    console.log('success');
                }.bind(this),
                error: function(xhr, status, err) {
                    console.error(err);
                    //TODO: callback?
                }.bind(this)
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
            json.action = {
                id: this.state.action.id,
                ingredients: this.refs.actionIngredients.toJSON()
            }
            json.trigger = {
                id: this.state.trigger.id,
                ingredients: this.refs.triggerIngredients.toJSON()
            }
            return json;
        },

        cancelClicked: function(evt) {
            evt.preventDefault();
            evt.stopPropagation();
            this.props.onCancel();
        },

        render: function() {
            var triggerChild, actionChild;
            var spinner = <div className="text-center"><i className="fa fa-spinner fa-spin"></i></div>;
            if (this.state.trigger) {
                // Render the selected trigger
                triggerChild = <IngredientList ref="triggerIngredients" ingredients={this.state.trigger.ingredients} />
            } else if (this.state.triggers) {
                // Render the trigger list
                triggerChild = <TriggerList triggers={this.state.triggers} selected={this.triggerSelected}/>
            }
            else if (this.state.triggerCookBookID !== -1) {
                // Got a cook book, must be loading triggers
                triggerChild = spinner;
            }
            else {
                triggerChild = <CookBookList cookBooks={this.props.cookBooks} selected={this.triggerCookBookSelected}/>
            }

            if (this.state.action) {
                actionChild = <IngredientList ref="actionIngredients" ingredients={this.state.action.ingredients} />
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
                    <div className="form-group">
                        <label htmlFor="name">Name</label>
                        <input value={this.state.name} onChange={this.handleNameChange} className="name form-control" type="text" id="name"/>
                    </div>
                    <div className="form-group">
                        <label htmlFor="description">Description</label>
                        <input value={this.state.description} onChange={this.handleDescriptionChange} className="description form-control" type="text" id="description"/>
                    </div>
                    <div className="trigger">
                        <h3>Trigger</h3>
                        {triggerChild}
                    </div>
                    <div className="action">
                        <h3>Action</h3>
                        {actionChild}
                    </div>
                    <div className="clearfix footer">
                        <button className="btn btn-default pull-right" onClick={this.cancelClicked}>Cancel</button>
                        <button className="btn btn-primary pull-right" onClick={this.saveClicked}>Save</button>
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
                    <CookBook data={cookBook} selected={self.handleClick}/>
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
                    <Trigger data={trigger} selected={self.handleClick} />
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
                    <Action data={action} selected={self.handleClick}/>
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
                return (
                    <Ingredient data={ingredient} ref={ingredient.id} />
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
                json[key] = self.refs[key].value();
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
            //this.props.data.id

            var input;
            switch(this.props.data.type) {
            case 'string':
            case 'duration':
            case 'integer':
            case 'float':
                input = <input className="ingredientInput form-control" type="text" onChange={this.changeHandler} id={this.getNextIdAndIncrement()}/>;
                break;
            case 'boolean':
                input = <input className="ingredientInput" type="checkbox" value="true" onChange={this.changeHandler} id={this.getNextIdAndIncrement()}/>;
                break;
            case 'datetime':
                //TODO: show calendar
                break;
            default:
                throw 'unknown ingredient type: ' + this.props.data.type;
            }

            //TODO: htmlFor should be unique
            return (
                <div>
                    <div className="form-group">
                        <label htmlFor={this.getCurrentId()}>{this.props.data.name}</label>
                        <p>{this.props.data.description}</p>
                        {input}
                    </div>
                </div>
            );
        },

        value: function() {
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
