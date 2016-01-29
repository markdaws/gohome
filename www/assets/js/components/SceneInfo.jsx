var React = require('react');
var SaveBtn = require('./SaveBtn.jsx');
var InputValidationMixin = require('./InputValidationMixin.jsx');
var UniqueIdMixin = require('./UniqueIdMixin.jsx');
var CommandInfo = require('./CommandInfo.jsx');

var SceneInfo = React.createClass({
    mixins: [InputValidationMixin, UniqueIdMixin],

    getInitialState: function() {
        return {
            id: this.props.scene.id || '',
            name: this.props.scene.name || '',
            address: this.props.scene.address || '',
            managed: (this.props.scene.managed == undefined) ? true : this.props.scene.managed,
            commands: this.props.scene.commands || [],
            zones: this.props.zones || []
            //TODO: readonly id
        };
    },

    componentWillReceiveProps: function(nextProps) {
        //Needed?
        if (nextProps.zones) {
            this.setState({ zones: nextProps.zones });
        }
    },
    
    toJson: function() {
        var s = this.state;
        return {
            id: this.state.id,
            name: this.state.name,
            address: this.state.address,
            managed: this.state.managed,
        };
    },

    deleteScene: function() {
        var self = this;
        $.ajax({
            url: '/api/v1/systems/123/scenes/' + this.state.id,
            type: 'DELETE',
            cache: false,
            success: function(data) {
                self.props.onDestroy(self.state.id);
            }.bind(this),
            error: function(xhr, status, err) {
                console.error(err);
            }.bind(this)
        });
    },

    render: function() {
        var commands
        //TODO: remove
        this.state.managed = true;
        var self = this;
        if (this.state.managed) {
            //TODO: What is the key here? Commands don't have ids ...
            commands = this.state.commands.map(function(command) {
                return <CommandInfo zones={self.props.zones} command={command} />
            });
        } else {
            commands = <p>This is an unmanaged scene. The scene is controlled by a 3rd party device so we can&apos;t show the individual commands it will execute. To modify the scene you will need to use the app provided with the 3rd party device.</p>
        }
        return (
            <div className="cmp-SceneInfo well">
              <div className="clearfix">
                <button className="btn btn-danger pull-right" onClick={this.deleteScene}>Delete Scene</button>
              </div>
              <div className={this.addErr("form-group", "name")}>
                <label className="control-label" htmlFor={this.uid("name")}>Name*</label>
                <input value={this.state.name} data-statepath="name" onChange={this.changed} className="name form-control" type="text" id={this.uid("name")}/>
                {this.errMsg("name")}
              </div>
              <div className={this.addErr("form-group", "id")}>
                <label className="control-label" htmlFor={this.uid("id")}>ID</label>
                <input value={this.state.id} readOnly={this.isReadOnly("id")} data-statepath="id" onChange={this.changed} className="id form-control" type="text" id={this.uid("id")}/>
                {this.errMsg("id")}
            </div>
              {/*<!-- TODO: Only needed for unmanaged scenes -->*/}
              <div className={this.addErr("form-group", "address")}>
                <label className="control-label" htmlFor={this.uid("address")}>Address</label>
                <input value={this.state.adddress} data-statepath="address" onChange={this.changed} className="address form-control" type="text" id={this.uid("address")}/>
                {this.errMsg("address")}
              </div>
            
              <h3>Commands</h3>
              {commands}
              <button className="btn btn-primary btnAddCommand">Add Command</button>
            </div>
        );
    }
});
module.exports = SceneInfo;