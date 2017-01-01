// These values match the attribute.go file, if you update that file you will
// need to update the values here too

function Attribute() {}

var Type = {
    OpenClose: 'OpenClose',
    OnOff: 'OnOff',
    Brightness: 'Brightness',
    HSL: 'HSL',
    Offset: 'Offset',
    Temperature: 'Temperature'
};

var Perms = {
    ReadOnly: 'r',
    ReadWrite: 'rw'
};

function OnOff(){}
OnOff.States = {
    0: 'Unknown',
    1: 'Off',
    2: 'On'
};

function OpenClose(){}
OpenClose.States = {
    0: 'Unknown',
    1: 'Closed',
    2: 'Open'
};

module.exports = {
    Type: Type,
    Attribute: Attribute,
    OnOff: OnOff,
    OpenClose: OpenClose,
    Perms: Perms
};
