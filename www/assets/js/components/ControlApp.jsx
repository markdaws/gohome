var React = require('react');
var ReactDOM = require('react-dom');
var ReactRedux = require('react-redux');
var System = require('./System.jsx');
var SceneList = require('./SceneList.jsx');
var ZoneSensorList = require('./ZoneSensorList.jsx');
var Logging = require('./Logging.jsx');
var RecipeApp = require('./RecipeApp.jsx');
var Constants = require('../constants.js');
var SceneActions = require('../actions/SceneActions.js');
var SensorActions = require('../actions/SensorActions.js');
var SystemActions = require('../actions/SystemActions.js');
var ZoneActions = require('../actions/ZoneActions.js');
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'ControlApp',
    prefix: 'b-'
});
require('../../css/components/ControlApp.less')

var ControlApp = React.createClass({
    getDefaultProps: function() {
        return {
            buttons: [],
            devices: [],
            zones: [],
            sensors: [],

            //TODO: Change to array
            scenes: { items: [] }
        };
    },

    componentDidMount: function() {
        //TODO: Have a loading screen until all of these have loaded
        this.props.loadAllDevices();
        this.props.loadAllZones();
        this.props.loadAllScenes();
        this.props.loadAllButtons();
        this.props.loadAllSensors();
    },

    render: function() {

        var zoneBody;
        if (this.props.zones.length === 0) {
            zoneBody = (
                <h5 className="emptyMessage">You don't have any zones. Go to the devices tab and import a Device, or manually edit the .json system file.</h5>
            );
        } else {
            zoneBody = (
                <ZoneSensorList
                    sensors={this.props.sensors}
                    zones={this.props.zones}
                    devices={this.props.devices}
                />
            );
        }

        var emptySceneBody;
        if (this.props.scenes.items.length === 0) {
            emptySceneBody = (
                <h5 className="emptyMessage">You don't have any scenes.  Click on the Edit button to add a new Scene.</h5>
            );
        }

        return (
            <div {...classes()}>
                <ul className="nav nav-tabs" role="tablist">
                    <li role="presentation" className="">
                        <a href="#scenes" role="tab" aria-controls="scenes" data-toggle="tab">
                            <i className="fa fa-sliders"></i>
                        </a>
                    </li>
                    <li role="presentation" className="active">
                        <a href="#zones" role="tab" aria-controls="zones" data-toggle="tab">
                            <i className="fa fa-code-fork"></i>
                        </a>
                    </li>
                    <li role="presentation" className="">
                        <a href="#system" role="tab" aria-controls="system" data-toggle="tab">
                            <i className="fa fa-tablet"></i>
                        </a>
                    </li>
                    {/*
                    //TODO: re-enable after v1.0
                    <li role="presentation">
                    <a href="#logging" role="tab" aria-controls="logging" data-toggle="tab">Logging</a>
                    </li>
                    <li role="presentation">
                    <a href="#recipes" role="tab" aria-controls="recipes" data-toggle="tab">Recipes</a>
                    </li>
                    */}
                </ul>
                <div className="tab-content">
                    <div role="tabpanel" className="tab-pane fade" id="scenes">
                        <div className={(this.props.appLoadStatus.scenesLoaded ? "" : "hideTabContent")}>
                            <SceneList scenes={this.props.scenes} buttons={this.props.buttons} zones={this.props.zones} />
                            {emptySceneBody}
                        </div>
                        <i className={"fa fa-spinner fa-spin " + (this.props.appLoadStatus.scenesLoaded ? "hidden" : "")}></i>
                    </div>
                    <div role="tabpanel" className="tab-pane active" id="zones">
                        <div className={(this.props.appLoadStatus.zonesLoaded ? "" : "hideTabContent")}>
                            {zoneBody}
                        </div>
                        <i className={"fa fa-spinner fa-spin " + (this.props.appLoadStatus.zonesLoaded ? "hidden" : "")}></i>
                    </div>
                    <div role="tabpanel" className="tab-pane fade" id="system">
                        <div className={(this.props.appLoadStatus.devicesLoaded ? "" : "hideTabContent")}>
                            <System devices={this.props.devices}/>
                        </div>
                        <i className={"fa fa-spinner fa-spin " + (this.props.appLoadStatus.scenesLoaded ? "hidden" : "")}></i>
                    </div>
            {/*
            <div role="tabpanel" className="tab-pane fade" id="logging">
            <Logging />
            </div>
            <div role="tabpanel" className="tab-pane fade" id="recipes">
                    <RecipeApp />
                    </div>
                    */}
                </div>
            </div>
        );
    }
});

function mapStateToProps(state) {
    return {
        devices: state.system.devices,
        zones: state.zones,
        scenes: state.scenes,
        buttons: state.buttons,
        sensors: state.sensors,
        appLoadStatus: state.appLoadStatus
    };
}

function mapDispatchToProps(dispatch) {
    return {
        loadAllButtons: function() {
            dispatch(SystemActions.loadAllButtons());
        },
        loadAllDevices: function() {
            dispatch(SystemActions.loadAllDevices());
        },
        loadAllScenes: function() {
            dispatch(SceneActions.loadAll());
        },
        loadAllZones: function() {
            dispatch(ZoneActions.loadAll());
        },
        loadAllSensors: function() {
            dispatch(SensorActions.loadAll());
        }
    }
}

module.exports = ReactRedux.connect(mapStateToProps, mapDispatchToProps)(ControlApp);
