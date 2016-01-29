var React = require('react');

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
module.exports = RecipeInfo;