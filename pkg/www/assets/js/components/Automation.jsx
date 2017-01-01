import React from 'react';
import Api from '../utils/API.js';
import BEMHelper from 'react-bem-helper';

var classes = new BEMHelper({
    name: 'Automation',
    prefix: 'b-'
});
require('../../css/components/Automation.less')

class Automation extends React.Component {
    constructor() {
        super();
        this.handleClick = this.handleClick.bind(this);
    }

    handleClick(event) {
        Api.automationTest(this.props.automation.tempId, function(err, data) {
            if (err) {
                //TODO: Show error/success
                console.error(err);
            }
        });
    }

    render() {
        return (
            <div {...classes()}>
                <div {...classes('name')}>
                    {this.props.automation.name}
                </div>
                <a role="button" {...classes('activate', '', 'btn btn-primary')} onClick={this.handleClick}>
                    Test
                </a>
            </div>
        )
    }
}
module.exports = Automation;
