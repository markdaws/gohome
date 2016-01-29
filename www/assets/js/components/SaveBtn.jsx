var React = require('react');

module.exports = React.createClass({
    getInitialState: function() {
        return {
            current: 'default',
            timeout: -1
        };
    },
    reset: function() {
        this.setState({ current: 'default' });
    },
    saving: function() {
        this.setState({ current: 'saving' });
    },
    success: function() {
        this.setState({ current: 'success' });
    },
    failure: function() {
        this.clearTimeout();
        var self = this;
        var timeout = setTimeout(function() {
            self.reset();
        }, 1500);
        this.setState({
            current: 'failure',
            timeout: timeout
        });
    },
    clicked: function() {
        this.props.clicked();
    },
    clearTimeout: function() {
        clearTimeout(this.state.timeout);
    },
    render: function() {
        var btnType, body;
        var disabled = true;
        switch (this.state.current) {
        case 'default':
            btnType = "btn-primary";
            body = (
                <div>
                  {this.props.text}
                </div>
            );
            disabled = false;
            break;
        case 'saving':
            btnType = 'btn-primary';
            body = (
                <div>
                  <i className="fa fa-spinner fa-spin"></i>
                </div>
            );                
            break;
        case 'success':
            btnType = "btn-success";
            body = (
                <div>
                  <span className="glyphicon glyphicon-ok"></span>
                </div>
            );
            break;
        case 'failure':
            btnType = "btn-danger";
            body = (
                <div>
                  Error
                </div>
            );
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
