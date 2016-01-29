var React = require('react');
var NewRecipe = require('./NewRecipe.jsx');
var RecipeList = require('./RecipeList.jsx');

module.exports = React.createClass({
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
