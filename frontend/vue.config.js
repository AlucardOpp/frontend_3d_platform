const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  devServer: {
    port: 8080,
    proxy: {
      '/api': {
        target: 'http://localhost:8095',
        secure: false,
        changeOrigin: true,
      },
      '/upload': {
        target: 'http://localhost:8095',
        secure: false,
        changeOrigin: true,
      },
    }
  }
})
