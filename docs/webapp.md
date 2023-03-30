# nuxtjs h5 app

## https://www.nuxtjs.cn/guide/installation

## https://typescript.nuxtjs.org/cookbook/plugins

```
npx create-nuxt-app webapp

#### To get started:

cd nuxtjs-app
npm run dev

#### To build & start for production:

cd nuxtjs-app
npm run build
npm run start

#### To test:

cd nuxtjs-app
npm run test

```

## 优化.

```
#### 图片压缩.
npm install --save-dev imagemin
npm install --save-dev imagemin-mozjpeg imagemin-pngquant

npm run images


#### 去掉生产日志.
npm install babel-plugin-transform-remove-console -D


#### workbox.
cd /Users/tokenfun/taurus/bitdisk/gitlab/bytezero/workbox
node ./packages/workbox-cli/build/bin.js copyLibraries build/
cp -rf build/workbox-v6.5.4 ../nuxtjs-app/static/js/workbox

```

## 依赖.

```
npm install @vant/touch-emulator
npm install git+ssh://git@gitlab.cume.cc:bytezeroc/peerjs.git
npm install pako
npm install streamsaver @types/streamsaver
npm install @nuxtjs/device
npm install @nuxtjs/i18n
npm install less less-loader@7

npm install express
npm install qrcode

```
