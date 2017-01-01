var React = require('react');
var DiscoverDevices = require('./DiscoverDevices.jsx');
var Api = require('../utils/API.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'Import',
    prefix: 'b-'
});
require('../../css/components/Import.less')

var Import = React.createClass({
    getInitialState: function() {
        return {
            discoverers: [],
            selectedDiscoverer: ''
        };
    },

    productSelected: function(evt) {
        this.setState({ selectedDiscoverer: evt.target.value });
    },

    componentDidMount: function() {
        Api.discoverersList(function(err, data) {
            if (err) {
                //TODO:
                console.error(err);
                return;
            }

            this.setState({ discoverers: data });
        }.bind(this));
    },
    
    render: function() {
        var body;

        var discoverer;
        for (var i=0; i<this.state.discoverers.length; ++i) {
            if (this.state.discoverers[i].id === this.state.selectedDiscoverer) {
                discoverer = this.state.discoverers[i];
                break;
            }
        }
        if (discoverer) {
            body = <DiscoverDevices discoverer={discoverer} key={discoverer.id} />;
        }

        var options = this.state.discoverers.map(function(discoverer) {
            return <option key={discoverer.id} value={discoverer.id}>{discoverer.name}</option>;
        });
        
        return (
            <div {...classes()}>
                <h2 {...classes('header')}>Add Hardware</h2>
                <select className="form-control" onChange={this.productSelected} value={this.state.selectedProduct}>
                    <option value="">Choose a product ...</option>
                    {options}
                </select>
                <div {...classes('content')}>
                    {body}
                </div>
            </div>
        );
    }
});
module.exports = Import;
