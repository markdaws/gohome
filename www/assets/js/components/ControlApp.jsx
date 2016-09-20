var React = require('react');
var System = require('./System.jsx');
var SceneList = require('./SceneList.jsx');
var ZoneList = require('./ZoneList.jsx');
var Logging = require('./Logging.jsx');
var RecipeApp = require('./RecipeApp.jsx');
var Constants = require('../constants.js');

var ControlApp = React.createClass({
    getInitialState: function() {
        return { devices: [], buttons: [] };
    },

    componentDidMount: function() {
        //TODO: remove, scenes can fetch buttons
        $.ajax({
            url: '/api/v1/systems/123/buttons',
            dataType: 'json',
            cache: false,
            success: function(data) {
                this.setState({buttons: data});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(err.toString());
            }.bind(this)
        });
        
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
                            buttons={this.state.buttons} />
                    </div>
                    <div role="tabpanel" className="tab-pane fade" id="zones">
                        <ZoneList />
                    </div>
                    <div role="tabpanel" className="tab-pane fade" id="system">
                        <System />
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
module.exports = ControlApp;
