var React = require('react');
var ReactDOM = require('react-dom');
var ControlApp = require('./components/ControlApp.jsx');
var Provider = require('react-redux').Provider;
var store = require('./store');

ReactDOM.render(
    <Provider store={store}>
        <ControlApp />
    </Provider>,
    document.getElementsByClassName('content')[0]
);


