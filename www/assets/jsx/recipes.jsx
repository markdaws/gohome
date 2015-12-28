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

    var AssetsMixin = {
        getImageUrl: function(imageName) {
            return 'assets/images/' + imageName;
        }
    }

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
                        <h2 className="pull-left">Recipes</h2>
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
                    <RecipeInfo recipe={recipe} onDestroy={self.recipeDestroyed}/>
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
                enabled: this.props.enabled,
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

    ReactDOM.render(<RecipeApp />, document.body.getElementsByClassName('content')[0]);
})()