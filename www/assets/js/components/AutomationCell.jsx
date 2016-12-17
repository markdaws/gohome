import React from 'react';
import BEMHelper from 'react-bem-helper';

var classes = new BEMHelper({
    name: 'AutomationCell',
    prefix: 'b-'
});
require('../../css/components/AutomationCell.less')

const AutomationCell = ({automation}) => {
    return (
        <div {...classes()}>
            <div {...classes('icon')}>
                <i className="icomoon-ion-ios-cog-outline"></i>
            </div>
            <div {...classes('name')}>
                {automation.name}
            </div>
        </div>
    );
};
module.exports = AutomationCell;
