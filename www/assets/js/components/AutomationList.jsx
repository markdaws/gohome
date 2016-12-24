import React from 'React';
import AutomationCell from './AutomationCell.jsx';
import Automation from './Automation.jsx';
import Grid from './Grid.jsx';
import BEMHelper from 'react-bem-helper';

var classes = new BEMHelper({
    name: 'AutomationList',
    prefix: 'b-'
});
import '../../css/components/AutomationList.less'

const AutomationList = ({automations = []}) => {
    const gridCells = automations.map(function(automation) {
        return {
            key: automation.tempId,
            cell: <AutomationCell automation={automation} />,
            content: <Automation automation={automation} key={automation.tempId}/>
        };
    });

    return (
        <div {...classes()}>
            <h2 {...classes('header')}>Automation</h2>
            <Grid cells={gridCells} />
        </div>
    );
}
module.exports = AutomationList;
