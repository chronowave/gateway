const path = require('path');

module.exports = {
  entry: {
    main: { import: './src/main/index.tsx', filename: 'main.js' },
    //worker: { import: './src/worker/index.tsx', filename: 'worker.js' },
  },
  module: {
    rules: [
      {
        test: /\.(tsx)$/,
        exclude: /node_modules/,
        use: ['ts-loader']
      },
      {
        test: /\.worker\/index\.(tsx)$/i,
        use: [{
          loader: 'comlink-loader',
          options: {
            singleton: true
          }
        }]
      }
    ]
  },
  resolve: {
    extensions: ['.tsx', '.js']
  },
  output: {
    publicPath: '/',
    path: path.resolve(__dirname, 'dist'),
  },
  devServer: {
    contentBase: './dist'
  }
};