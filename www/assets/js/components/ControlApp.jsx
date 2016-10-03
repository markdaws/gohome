var React = require('react');
var ReactRedux = require('react-redux');
var System = require('./System.jsx');
var SceneList = require('./SceneList.jsx');
var ZoneList = require('./ZoneList.jsx');
var Logging = require('./Logging.jsx');
var RecipeApp = require('./RecipeApp.jsx');
var Constants = require('../constants.js');
var SceneActions = require('../actions/SceneActions.js');
var SystemActions = require('../actions/SystemActions.js');
var ZoneActions = require('../actions/ZoneActions.js');

var ControlApp = React.createClass({
    getDefaultProps: function() {
        return {
            buttons: [],
            devices: [],
            zones: [],
            scenes: {}
        };
    },

    componentDidMount: function() {
        //TODO: Have a loading screen until all of these have loaded
        this.props.loadAllDevices();
        this.props.loadAllZones();
        this.props.loadAllScenes();
        this.props.loadAllButtons();
    },

    render: function() {
        return (
            <div className="cmp-ControlApp">
                <ul className="nav nav-tabs" role="tablist">
                    <li role="presentation" className="active">
                        <a href="#scenes" role="tab" aria-controls="scenes" data-toggle="tab">Scenes</a>
                    </li>
                    <li role="presentation" className="">
                        <a href="#zones" role="tab" aria-controls="zones" data-toggle="tab">Zones</a>
                    </li>
                    <li role="presentation" className="">
                        <a href="#system" role="tab" aria-controls="system" data-toggle="tab">System</a>
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
                    <div role="tabpanel" className="tab-pane active" id="scenes">
                        <SceneList
                            scenes={this.props.scenes}
                            buttons={this.props.buttons}
                            zones={this.props.zones} />
                    </div>
                    <div role="tabpanel" className="tab-pane fade" id="zones">
                        <ZoneList zones={this.props.zones}/>
                    </div>
                    <div role="tabpanel" className="tab-pane fade" id="system">
                        <System devices={this.props.devices}/>
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
        buttons: state.buttons
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
        }
    }
}

module.exports = ReactRedux.connect(mapStateToProps, mapDispatchToProps)(ControlApp);
