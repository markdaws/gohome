module.exports = {
    uid: function(field) {
        return this.state.cid + '.' + field
    },
    getErr: function(field) {
        if (!this.state.errors) {
            return null;
        }
        return this.state.errors[this.uid(field)];            
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
    changed: function(evt) {
        var statePath = evt.target.getAttribute('data-statepath');
        var s = {}
        s[statePath] = evt.target.value;
        this.setState(s);
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

