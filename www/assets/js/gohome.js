var React = require('react');
var ReactDOM = require('react-dom');
var ControlApp = require('./components/ControlApp.jsx');
var Provider = require('react-redux').Provider;
var store = require('./store');

var C1 = require('./components/Testr.jsx');

//TODO: Remove - testing
/*
ReactDOM.render(
    <Provider store={store}>
        <C1 />
    </Provider>,
    document.getElementsByClassName('content')[0]
);
*/

ReactDOM.render(
    <Provider store={store}>
        <ControlApp />
    </Provider>,
    document.getElementsByClassName('content')[0]
);

