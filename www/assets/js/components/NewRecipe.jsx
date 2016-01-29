var React = require('react');
var IngredientList = require('./IngredientList.jsx');
var TriggerList = require('./TriggerList.jsx');
var ActionList = require('./ActionList.jsx');
var CookBookList = require('./CookBookList.jsx');

module.exports = React.createClass({
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
