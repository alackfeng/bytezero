module.exports = {
  root: true,
  env: {
    browser: true,
    node: true,
  },
  extends: [
    '@nuxtjs/eslint-config-typescript',
    'plugin:nuxt/recommended',
    'prettier',
  ],
  plugins: [],
  // add your custom rules here
  rules: {
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debug': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'vue/first-attribute-linebreak': [
      'error',
      { singleline: 'ignore', multiline: 'beside' },
    ],
  },
}
