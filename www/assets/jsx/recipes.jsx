(function() {

    var App = React.createClass({
        getInitialState: function() {
            return { cookBooks: [] }
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
        
        render: function() {
            return (
                <div>
                    <CookBookList cookBooks={this.state.cookBooks} />
                </div>
            );
        }
    });

    var CookBookList = React.createClass({
        render: function() {
            var self = this;
            var cookBookNodes = this.props.cookBooks.map(function(cookBook) {
                return (
                    <CookBook data={cookBook} />
                );
            });
            return (
                <div className="cookBookList">
                    <h1>CookBooks</h1>
                    {cookBookNodes}
                </div>
            );
        }
    });

    var CookBook = React.createClass({
        getInitialState: function() {
            return { triggers: [], actions: [] }
        },

        componentDidMount: function() {
            $.ajax({
                url: '/api/v1/cookbooks/' + this.props.data.id,
                dataType: 'json',
                cache: false,
                success: function(data) {
                    this.setState({triggers: data.triggers});
                }.bind(this),
                error: function(xhr, status, err) {
                    console.error(err.toString());
                }.bind(this)
            });
        },
        
        render: function() {
            //TODO: handle click, update main page to show a single cookbook
            //Show actions
            return (
                <div>
                    {this.props.data.name} : {this.props.data.description}
                    <TriggerList triggers={this.state.triggers} />
                </div>
            );
        }
    });

    var TriggerList = React.createClass({
        render: function() {
            var self = this;
            var triggerNodes = this.props.triggers.map(function(trigger) {
                return (
                    <Trigger data={trigger} />
                );
            });

            return (
                <div>
                    {triggerNodes}
                </div>
            );
        }
    });

    var Trigger = React.createClass({
        render: function() {
            return (
                <div>
                    <div>{this.props.data.name} : {this.props.data.description}</div>
                    <IngredientList ingredients={this.props.data.ingredients} />
                </div>
            );
        }
    });

    var ActionList = React.createClass({
        render: function() {
            return (
                <div></div>
            );
        }
    });

    var Action = React.createClass({
        render: function() {
            return (
                <div></div>
            );
        }
    });
    
    var IngredientList = React.createClass({
        render: function() {
            var self = this;
            var ingredientNodes = this.props.ingredients.map(function(ingredient) {
                return (
                    <Ingredient data={ingredient} />
                );
            });

            return (
                <div>
                    {ingredientNodes}
                </div>
            );
        }
    });

    var Ingredient = React.createClass({
        render: function() {
        //TODO: Render different depending on if type is number|string|boolean
            return (
                <div>{this.props.data.id} : {this.props.data.name} : {this.props.data.description} : {this.props.data.type}</div>
            );
        }
    });
    //NumberIngredient
    //StringIngredient
    //BooleanIngredient
    //ListIngredient?? max length, mixed types?
    //TODO: How to render button|scene|zone|device, need hierarchy
    //TODO: Recipe? Combo of name,desc,id,trigger,action
    
    React.render(<App />, document.body);
})()

/*
CookBooks
CookBook
Trigger
Action
Ingredient

CookBook
 - name
 - description
 - triggers
 - actions

Trigger
 - name
 - description
 - ingredients

Action
 - name
 - description
 - ingredients

Ingredient
 - id
 - name
 - description
 - type string|number|boolean|...

//TODO: UI for choosing scene | button | zone

//Flow
 - pick a cook book by clicking on it
 - renders all the triggers
 - pick a trigger
 - for each ingredient render one by one, with Next as user populates them
 - THEN
 - renders all actions
 - pick one by clicking on it
 - for each ingredient render one by one, with Next as user populates them
 - once finished - choose a name for your recipe, click save (can recipes be partially filled?)
*/