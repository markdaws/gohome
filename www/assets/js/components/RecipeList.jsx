var React = require('react');
var RecipeInfo = require('./RecipeInfo.jsx');

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
module.exports = RecipeList;