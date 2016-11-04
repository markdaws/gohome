var React = require('react');
var DiscoverDevices = require('./DiscoverDevices.jsx');
//TODO: Delete
var ImportTCP600GWB = require('./ImportTCP600GWB.jsx');
var Api = require('../utils/API.js');

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
            switch (discoverer.type) {
                case 'ScanDevices':
                    // Need to scan the network
                    body = <DiscoverDevices discoverer={discoverer} key={discoverer.id} />;
                    break;

                case 'FromString':
                    // This importer imports from a user provided string
                    //TODO:
                    break;
            }
        }

        var options = this.state.discoverers.map(function(discoverer) {
            return <option key={discoverer.id} value={discoverer.id}>{discoverer.name}</option>;
        });
        
        return (
            <div className="cmp-Import">
                <h3>Select a product to import</h3>
                <select className="form-control" onChange={this.productSelected} value={this.state.selectedProduct}>
                    <option value="">Choose ...</option>
                    {options}
                    {/*TODO: Remove
                    <option value="f7c029v2">Belkin WeMo Insight</option>
                    <option value="f7c043fc">Belkin WeMo Maker</option>
                    <option value="l-bdgpro2-wh">Lutron</option>
                    <option value="tcp600gwb">Connected By TCP Hub</option>
                    <option value="fluxwifi">Flux WIFI Bulb</option>*/}
                </select>
                <div className="content">
                    {body}
                </div>
            </div>
        );
    }
});
module.exports = Import;
