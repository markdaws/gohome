var React = require('react');
var Ingredient = require('./Ingredient.jsx');

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
module.exports = IngredientList;