var _current = 0;

module.exports = {
    getNextIdAndIncrement: function() {
        _current += 1;
        return _current;
    },

    getCurrentId: function() {
        return _current;
    }
};

