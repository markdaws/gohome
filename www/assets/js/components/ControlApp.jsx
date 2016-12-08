var React = require('react');
var ReactRedux = require('react-redux');
var System = require('./System.jsx');
var SceneList = require('./SceneList.jsx');
var FeatureList = require('./FeatureList.jsx');
var AutomationList = require('./AutomationList.jsx');
var Logging = require('./Logging.jsx');
var SceneActions = require('../actions/SceneActions.js');
var SystemActions = require('../actions/SystemActions.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'ControlApp',
    prefix: 'b-'
});
require('../../css/components/ControlApp.less')

var ControlApp = React.createClass({
    getDefaultProps: function() {
        return {
            devices: [],
            automations: [],
            //TODO: Change to array
            scenes: { items: [] }
        };
    },

    componentDidMount: function() {
        this.props.loadAllDevices();
        this.props.loadAllScenes();
        this.props.loadAllAutomation();
    },

    render: function() {
        var featureBody;
        if (this.props.devices.length === 0) {
            featureBody = (
                <h5 {...classes('empty-message-zones')}>You don't have any devices or feature. Go to the devices tab to get started.</h5>
            );
        } else {
            featureBody = (
                <FeatureList devices={this.props.devices} />
            );
        }

        var emptySceneBody;
        if (this.props.scenes.items.length === 0) {
            emptySceneBody = (
                <h5 {...classes('empty-message-scenes')}>You don't have any scenes.  Click on the Edit button to add a new Scene.</h5>
            );
        }

        var emptyAutomationBody;
        if (this.props.automations.length === 0) {
            emptyAutomationBody = (
                <h5 {...classes('empty-message-automations')}>You don't have any automation.  Follow the instructions <a target="_blank" href="https://github.com/markdaws/gohome/blob/master/docs/automation.md">here</a> to create some automation rules.</h5>
            );
        }

        return (
            <div {...classes()}>
                <ul className="nav nav-tabs" role="tablist">
                    <li role="presentation" className="active">
                        <a href="#features" role="tab" aria-controls="features" data-toggle="tab">
                            <i className="fa fa-code-fork"></i>
                        </a>
                    </li>
                    <li role="presentation" className="">
                        <a href="#scenes" role="tab" aria-controls="scenes" data-toggle="tab">
                            <i className="fa fa-sliders"></i>
                        </a>
                    </li>
                    <li role="presentation" className="">
                        <a href="#automation" role="tab" aria-controls="automation" data-toggle="tab">
                            <i className="fa fa-cogs"></i>
                        </a>
                    </li>
                    <li role="presentation" className="">
                        <a href="#system" role="tab" aria-controls="system" data-toggle="tab">
                            <i className="fa fa-tablet"></i>
                        </a>
                    </li>
                </ul>
                <div className="tab-content">
                    <div role="tabpanel" className="tab-pane active" id="features">
                        <div className={(this.props.appLoadStatus.devicesLoaded ? "" : "hideTabContent")}>
                            {featureBody}
                        </div>
                        <div {...classes('spinner')}>
                            <i className={"fa fa-spinner fa-spin " + (this.props.appLoadStatus.devicesLoaded ? "hidden" : "")}></i>
                        </div>
                    </div>
                    <div role="tabpanel" className="tab-pane fade" id="scenes">
                        <div className={(this.props.appLoadStatus.scenesLoaded ? "" : "hideTabContent")}>
                            <SceneList scenes={this.props.scenes} devices={this.props.devices} />
                            {emptySceneBody}
                        </div>
                        <div {...classes('spinner')}>
                            <i className={"fa fa-spinner fa-spin " + (this.props.appLoadStatus.scenesLoaded ? "hidden" : "")}></i>
                        </div>
                    </div>
                    <div role="tabpanel" className="tab-pane fade" id="automation">
                        <div className={(this.props.appLoadStatus.automationLoaded ? "" : "hideTabContent")}>
                            <AutomationList automations={this.props.automations} />
                            {emptyAutomationBody}
                        </div>
                        <div {...classes('spinner')}>
                            <i className={"fa fa-spinner fa-spin " + (this.props.appLoadStatus.automationLoaded ? "hidden" : "")}></i>
                        </div>
                    </div>
                    <div role="tabpanel" className="tab-pane fade" id="system">
                        <div className={(this.props.appLoadStatus.devicesLoaded ? "" : "hideTabContent")}>
                            <System devices={this.props.devices}/>
                        </div>
                        <div {...classes('spinner')}>
                            <i className={"fa fa-spinner fa-spin " + (this.props.appLoadStatus.scenesLoaded ? "hidden" : "")}></i>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
});

function mapStateToProps(state) {
    return {
        devices: state.system.devices,
        automations: state.automations,
        scenes: state.scenes,
        appLoadStatus: state.appLoadStatus,
        errors: state.errors
    };
}

function mapDispatchToProps(dispatch) {
    return {
        loadAllDevices: function() {
            dispatch(SystemActions.loadAllDevices());
        },
        loadAllScenes: function() {
            dispatch(SceneActions.loadAll());
        },
        loadAllAutomation: function() {
            dispatch(SystemActions.loadAllAutomation());
        }
    }
}

module.exports = ReactRedux.connect(mapStateToProps, mapDispatchToProps)(ControlApp);
