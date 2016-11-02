var React = require('react');
var ReactDOM = require('react-dom');
var ControlApp = require('./components/ControlApp.jsx');
var Provider = require('react-redux').Provider;
var store = require('./store');
var Login = require('./components/Login.jsx');
var Api = require('./utils/API.js');

/*
var C1 = require('./components/Testr.jsx');

//TODO: Remove - testing
ReactDOM.render(
    <Provider store={store}>
        <C1 />
    </Provider>,
    document.getElementsByClassName('content')[0]
);
*/

// If the user has logged in there is a session cookie, if not we show the login screen
var cookies = document.cookie.split(';');
var sid = '';
for (var i=0; i<cookies.length; ++i) {
    var index = cookies[i].indexOf('=');
    if (index === -1) {
        continue;
    }

    var name = cookies[i].substr(0, index);
    if (name === 'sid') {
        sid = cookies[i].substr(index+1);
        break;
    }
}

if (sid !== '') {
    // Need to set the SID so that we can call our APIs
    Api.setSID(sid);

    ReactDOM.render(
            <Provider store={store}>
            <ControlApp />
            </Provider>,
        document.getElementsByClassName('content')[0]
    );
} else {    
    ReactDOM.render(
            <Provider store={store}>
            <Login />
            </Provider>,
        document.getElementsByClassName('content')[0]
    );
}
