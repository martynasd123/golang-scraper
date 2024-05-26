const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');

module.exports = function (env, argv) {
    const isProduction = argv.mode === "production";

    return {
        entry: path.join(__dirname, "src", "index.tsx"),
        devtool: 'inline-source-map',
        output: {
            publicPath: '/',
            filename: 'bundle.js',
            path: path.resolve(__dirname, "dist"),
        },
        module: {
            rules: [
                {
                    test: /(\.?js$)|(\.?jsx$)/,
                    exclude: /node_modules/,
                    use: {
                        loader: "babel-loader",
                        options: {
                            presets: ['@babel/preset-env', '@babel/preset-react'],
                            envName: isProduction ? "production" : "development"
                        }
                    },
                },
                {
                    test: /\.tsx?$/,
                    use: 'ts-loader',
                    exclude: /node_modules/,
                },
                {
                    test: /\.less$/i,
                    use: [
                        "style-loader",
                        "css-loader",
                        "less-loader"
                    ],
                },
                {
                    test: /\.(png|jp(e*)g|svg|gif)$/,
                    type: 'asset/resource',
                    exclude: /node_modules/,
                }
            ]
        },
        resolve: {
            extensions: [".js", ".jsx", ".ts", ".tsx", ".scss", ".css"]
        },
        plugins: [
            new HtmlWebpackPlugin({
                template: path.join(__dirname, "src", "index.html"),
                filename: "index.html",
            })
        ],
        devServer: {
            historyApiFallback: true,
            static: {
                directory: `dist/`,
            },
            allowedHosts: ["127.0.0.1"],
            port: 3000,
            proxy: [
                {
                    context: ['/api'],
                    target: 'http://localhost:8080',
                    changeOrigin: true,
                },
            ],
        }
    }
}