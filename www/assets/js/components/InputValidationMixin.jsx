var React = require('react');

module.exports = {
    uid: function(field) {
        return this.state.id + '_' + field
    },

    getErr: function(field) {
        var errors = this.state.errors;
        if (!errors) {
            return null;
        }
        return errors[this.uid(field)];
    },

    hasErr: function(field) {
        return this.getErr(field) != null;
    },

    errMsg: function(field) {
        var err = this.getErr(field);
        if (!err) {
            return;
        }
        return <span className="help-block">{"Error - " + err.message}</span>
    },

    addErr: function(classes, field) {
        if (this.hasErr(field)) {
            return classes + " has-error";
        }
        return classes;
    },

    changed: function(evt, cb) {
        var statePath = evt.target.getAttribute('data-statepath');
        var s = {}
        s[statePath] = evt.target.value;
        s.dirty = true;

        var errors = this.state['errors'] || {};
        delete errors[this.uid(statePath)];
        s.errors = errors;
        this.setState(s, cb);
    },

    isReadOnly: function(field) {
        var fields = this.props.readOnlyFields || ''
        var items = fields.split(',');
        for (var i=0; i<items.length; ++i) {
            if (items[i] === field) {
                return true;
            }
        }
        return false;
    }
};

