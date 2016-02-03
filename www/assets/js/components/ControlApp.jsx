var React = require('react');
var System = require('./System.jsx');
var SceneList = require('./SceneList.jsx');
var ZoneList = require('./ZoneList.jsx');
var Logging = require('./Logging.jsx');
var RecipeApp = require('./RecipeApp.jsx');

var ControlApp = React.createClass({
    getInitialState: function() {
        return { scenes: [], zones: [], devices: [], buttons: [] };
    },

    componentDidMount: function() {
        //TODO: Have a loading indicator for scenes + zones
        $.ajax({
            url: '/api/v1/systems/123/scenes',
            dataType: 'json',
            cache: false,
            success: function(data) {
                this.setState({scenes: data});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(err.toString());
            }.bind(this)
        });

        $.ajax({
            url: '/api/v1/systems/123/zones',
            dataType: 'json',
            cache: false,
            success: function(data) {
                this.setState({zones: data});
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(err.toString());
            }.bind(this)
        });

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
                  <a href="#system" role="tab" aria-controls="system" data-toggle="tab">System</a>
                </li>
                <li role="presentation">
                  <a href="#scenes" role="tab" aria-controls="scenes" data-toggle="tab">Scenes</a>
                </li>
                <li role="presentation">
                  <a href="#zones" role="tab" aria-controls="zones" data-toggle="tab">Zones</a>
                </li>
                <li role="presentation">
                  <a href="#logging" role="tab" aria-controls="logging" data-toggle="tab">Logging</a>
                </li>
                <li role="presentation">
                  <a href="#recipes" role="tab" aria-controls="recipes" data-toggle="tab">Recipes</a>
                </li>
              </ul>
              <div className="tab-content">
                <div role="tabpanel" className="tab-pane active" id="system">
                  <System />
                </div>
                <div role="tabpanel" className="tab-pane fade" id="scenes">
                  <SceneList
                    scenes={this.state.scenes}
                    zones={this.state.zones}
                    buttons={this.state.buttons} />
                </div>
                <div role="tabpanel" className="tab-pane fade" id="zones">
                  <ZoneList zones={this.state.zones} />
                </div>
                <div role="tabpanel" className="tab-pane fade" id="logging">
                  <Logging />
                </div>
                <div role="tabpanel" className="tab-pane fade" id="recipes">
                  <RecipeApp />
                </div>
              </div>
            </div>
        );
    }
});
module.exports = ControlApp;