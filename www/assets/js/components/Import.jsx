var React = require('react');
var DiscoverDevices = require('./DiscoverDevices.jsx');
var ImportTCP600GWB = require('./ImportTCP600GWB.jsx');

var Import = React.createClass({
    getInitialState: function() {
        return { selectedProduct: null };
    },

    productSelected: function(evt) {
        this.setState({ selectedProduct: evt.target.value });
    },

    render: function() {
        //TODO: Should get this list from the server, generate drop down automatically from
        //registered extensions
        
        var body
        switch(this.state.selectedProduct) {
        case 'tcp600gwb':
            body = <ImportTCP600GWB />
            break;
        case 'fluxwifi':
        case 'f7c029v2':
            body = <DiscoverDevices modelNumber={this.state.selectedProduct} />
            break;
        default:
            body = null;
        }
    
        return (
            <div className="cmp-Import">
                <h3>Select a product to import</h3>
                <select className="form-control" onChange={this.productSelected} value={this.state.selectedProduct}>
                    <option value="">Choose ...</option>
                    <option value="f7c029v2">Belkin WeMo Insight</option>          
                    <option value="l-bdgpro2-wh">Lutron</option>
                    <option value="tcp600gwb">Connected By TCP Hub</option>
                    <option value="fluxwifi">Flux WIFI Bulb</option>
                </select>
                <div className="content">
                    {body}
                </div>
            </div>
        );
    }
});
module.exports = Import;
