module.exports = {
    cssSafeIdentifier: function(value) {
        return value.replace(/:/g, '_');
    }
};
