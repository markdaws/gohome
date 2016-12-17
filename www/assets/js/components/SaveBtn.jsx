var React = require('react');
var Spinner = require('./Spinner.jsx');

var SaveBtn = React.createClass({
    getDefaultProps: function() {
        return {
            status: 'default'
        }
    },

    clicked: function() {
        this.props.clicked();
    },

    render: function() {
        var btnType, body;
        var disabled = true;

        switch(this.props.status) {
            case SaveBtn.STATUS.Saving:
                btnType = 'btn-primary';
                body = (
                    <Spinner />
                );
                break;
            case SaveBtn.STATUS.Error:
                btnType = "btn-danger";
                body = (
                    <div>
                        Error
                    </div>
                );
                break;
            case SaveBtn.STATUS.Success:
                btnType = "btn-success";
                body = (
                    <div>
                        <span className="glyphicon glyphicon-ok"></span>
                    </div>
                );
                break;
            default:
                btnType = "btn-primary";
                body = (
                    <div>
                        {this.props.text}
                    </div>
                );
                disabled = false;
                break;
        }

        var disabledClass = disabled ? " disabled" : "";
        return (
            <button className={"cmp-SaveBtn btn " + btnType + disabledClass} onClick={this.clicked}>
              {body}
            </button>
        );
    }
});

SaveBtn.STATUS = {
    Default: 'default',
    Saving: 'saving',
    Success: 'success',
    Error: 'error'
};
module.exports = SaveBtn;
