function Feature() {}
var Type = {
    Button: 'Button',
    CoolZone: 'CoolZone',
    HeatZone: 'HeatZone',
    LightZone: 'LightZone',
    Sensor: 'Sensor',
    Switch: 'Switch',
    Outlet: 'Outlet',
    WindowTreatment: 'WindowTreatment',
};

function LightZone() {}
LightZone.AttrIDs = {
    OnOff: 'onoff',
    Brightness: 'brightness',
    HSL: 'hsl'
};

function HeatZone() {}
HeatZone.AttrIDs = {
    CurrentTemp: 'currenttemp',
    TargetTemp: 'targettemp'
};

function Switch() {}
Switch.AttrIDs = {
    OnOff: 'onoff'
};

function Outlet() {}
Outlet.AttrIDs = {
    OnOff: 'onoff'
};

function WindowTreatment(){}
WindowTreatment.AttrIDs = {
    Offset: 'offset',
    OpenClose: 'openclose'
};

function cloneAttrs(sourceAttrs) {
    if (sourceAttrs == null) {
        return null;
    }

    var attrs = {};
    Object.keys(sourceAttrs).forEach(function(localId) {
        attrs[localId] = Object.assign({}, sourceAttrs[localId]);
    });
    return attrs;
}

module.exports = {
    Type: Type,
    Feature: Feature,
    LightZone: LightZone,
    HeatZone: HeatZone,
    Switch: Switch,
    Outlet: Outlet,
    WindowTreatment: WindowTreatment,
    cloneAttrs: cloneAttrs
};
