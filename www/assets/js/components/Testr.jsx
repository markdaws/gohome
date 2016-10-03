//Test file - delete

var React = require('react');
var ReactRedux = require('react-redux');

var Testr = React.createClass({
    getInitialState: function() {
        return {
            items: [
                { id: 1, name: 'one' },
                { id: 2, name: 'two' }
            ]
        };
    },

    componentWillMount: function() {
        console.log('cmpWillMount');

        setTimeout(function() {
            var items = this.state.items;
            items.push({ id: 3, name: 'three' });

            console.log('\n\n\n\n\n');
            this.setState({ items: items });

            setTimeout(function() {
                items[items.length-1] = { id: 10, name: 'ten' };

                console.log('\n\n\n\n\n');
                this.setState({ items: items });
            }.bind(this), 1000);
        }.bind(this), 1000);
    },

    componentDidMount: function() {
        console.log('cmpDidMount');
    },

    componentWillUnmount: function() {
        console.log('cmpWillUnmount');
    },

    componentWillUpdate: function() {
        console.log('cmpWillUpdate');
    },

    /*
    shouldComponentUpdate: function() {
        console.log('shouldCmpUpdate');
        return true;
    },*/

    componentWillReceiveProps: function(newProps) {
        console.log('cmpWillReceiveProps');
    },

    render: function() {
        console.log('cmpRender');

        var nodes = this.state.items.map(function(item) {
            return <TestrChild key={item.id} item={item}/>
        });
        return (
            <div>
              <div>hi there</div>
              {nodes}
            </div>
        );
    }
});

var TestrChild = React.createClass({
    getDefaultProps: function() {
        return {
            item: {}
        }
    },

    /*
    shouldComponentUpdate: function() {
        console.log('shouldCmpUpdate' + this.props.item.id);
        return true;
    },*/

    componentWillMount: function() {
        console.log('c-cmpWillMount:' + this.props.item.id);
    },

    componentDidMount: function() {
        console.log('c-cmpDidMount' + this.props.item.id);
    },

    componentWillUnmount: function() {
        console.log('c-cmpWillUnmount' + this.props.item.id);
    },

    componentWillUpdate: function() {
        console.log('c-cmpWillUpdate' + this.props.item.id);
    },

    componentWillReceiveProps: function(newProps) {
        console.log('c-cmpWillReceiveProps' + this.props.item.id);
    },

    render: function() {
        console.log('c-cmpRender' + this.props.item.id);
        return (
            <div>
                { this.props.item.name }
            </div>
        );
    }
});
module.exports = Testr;
