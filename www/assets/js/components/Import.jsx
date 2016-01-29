var React = require('react');
var ImportFluxWIFI = require('./ImportFluxWIFI.jsx');
var ImportTCP600GWB = require('./ImportTCP600GWB.jsx');

var Import = React.createClass({
    getInitialState: function() {
        return { selectedProduct: null };
    },

    productSelected: function(evt) {
        this.setState({ selectedProduct: evt.target.value });
    },

    render: function() {

        var body
        switch(this.state.selectedProduct) {
        case 'TCP600GWB':
            body = <ImportTCP600GWB />
            break;
        case 'FluxWIFI':
            body = <ImportFluxWIFI />
            break;
        default:
            body = null;
        }
        return (
            <div className="cmp-Import">
              <h3>Select a product to import</h3>
              <select className="form-control" onChange={this.productSelected} value={this.state.selectedProduct}>
                <option value="">Choose ...</option>
                <option value="LLL">Lutron</option>
                <option value="TCP600GWB">Connected By TCP Hub</option>
                <option value="FluxWIFI">Flux Wifi</option>
              </select>
              <div className="content">
                {body}
              </div>
            </div>
        )
    }
});
module.exports = Import;