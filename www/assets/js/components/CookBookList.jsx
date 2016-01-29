var React = require('react');
var CookBook = require('./CookBook.jsx');

module.exports = React.createClass({
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
