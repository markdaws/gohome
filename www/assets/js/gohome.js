var React = require('react');
var ReactDOM = require('react-dom');
var ControlApp = require('./components/ControlApp.jsx');

//TODO: Clean up
var apiUrl = '/api/v1/systems/123/scenes';
var apiUrlZones = '/api/v1/systems/123/zones';
ReactDOM.render(<ControlApp url={apiUrl} zoneUrl={apiUrlZones}/>, document.getElementsByClassName('content')[0]);


