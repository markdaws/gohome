var React = require('react');
var AssetsMixin = require('./AssetsMixin.jsx')

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
module.exports = CookBook;