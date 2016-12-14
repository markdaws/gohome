var path = require('path');
var webpack = require('webpack');
var CopyWebpackPlugin = require('copy-webpack-plugin');

var HtmlWebpackPlugin = require('html-webpack-plugin');
var HTMLWebpackPluginConfig = new HtmlWebpackPlugin({
    template: 'underscore-template!' + __dirname + '/assets/html/index.html',
    filename: 'index.html',
    inject: 'body',
    timestamp: new Date().getTime()
});

module.exports = {
    entry: './assets/js/gohome.js',
    output: {
        path: './dist/',
        filename: 'js/[name]-[hash].js'
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
            {
                test: /\.(jpe?g|png|gif|svg)$/i,
                loader: 'file?hash=sha512&digest=hex&name=images/[name]-[hash].[ext]'
            },
        ],
        postLoaders: [
            {
                test: /\.js$/,
                exclude: /node_modules/, // do not lint third-party code
                loader: 'jshint-loader'
            }
        ]
    },
    jshint: {
        strict: "global"
    },
    plugins: [
        HTMLWebpackPluginConfig,
        new CopyWebpackPlugin([
            { from: './assets/html/logout.html' },
            { from: './assets/css/**/*.css', to: './css', toType: 'dir', flatten: true },
            { from: './assets/fonts/', to: './fonts' },
            { from: './assets/images/', to: './images' },
            { from: './assets/js/ext/', to: './js' }
        ])
    ]
};
