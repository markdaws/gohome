var path = require('path');
var webpack = require('webpack');

module.exports = {
    entry: './assets/js/gohome.js',
    output: {
        path: './assets/js/',
        filename: 'gohome-out.js'
    },
    module: {
        loaders: [
            {
                test: /\.jsx?$/,
                loader: 'babel-loader',
                exclude: /node_modules/,
                query: {
                    cacheDirectory: true,
                    presets: ['es2015', 'react']
                }
            },
            {
                test: /\.less$/,
                loader: 'style-loader!css-loader!less-loader'
            },
        ]
    }
};
