var React = require('react');
var ReactDOM = require('react-dom');
var Api = require('../utils/API.js')
var BEMHelper = require('react-bem-helper');

var classes = new BEMHelper({
    name: 'Login',
    prefix: 'b-'
});
require('../../css/components/Login.less')

var Login = React.createClass({
    getInitialState: function() {
        return {
            error: null
        };
    },
    
    loginClicked: function(evt) {
        evt.preventDefault();
        
        this.setState({error: null});

        var $el = $(ReactDOM.findDOMNode(this))
        var login = $el.find('#login').val();
        var password = $el.find('#password').val();

        if (login === '') {
            this.setState({error: {msg: 'login cannot be empty' }});
            return;
        }
        if (password === '') {
            this.setState({error: {msg: 'password cannot be empty' }});
            return;
        }

        Api.sessionCreate(login, password, function(err, data) {
            if (err) {
                this.setState({error: err});
                return
            }

            window.location = '/';

            //expire cookie: document.cookie = 'sid=; expires=Thu, 01 Jan 1970 00:00:01 GMT;'
        }.bind(this));
    },
    
    render: function() {
        return (
            <div {...classes()}>
                <div {...classes('header')}>
                    <img {...classes('header-logo')} src="/assets/images/logo.png"></img>
                </div>

                <div {...classes('login-form')}>
                    <form>
                    <div className="form-group">
                        <input
                            type="text"
                            className="form-control"
                            id="login"
                            autoCapitalize="none"
                            autoCorrect="off"
                            placeholder="Login"></input>
                    </div>

                    <div className="form-group">
                        <input type="password" className="form-control" id="password" placeholder="Password"></input>
                    </div>
                    <button type="submit" onClick={this.loginClicked} className="btn btn-primary">Submit</button>
                    </form>
                </div>

                <div {...classes('error', this.state.error !== null ? '' : 'hidden')}>
                    The login/password combination is not valid, follow the instructions below to add a user or update your password
                </div>
                
                <div {...classes('need-credentials')}>
                    Don't have a login/password, or forgot your details, click <a href="https://github.com/markdaws/gohome">here</a> and follow the instructions.
                </div>
            </div>
        );
    }
});

module.exports = Login;
