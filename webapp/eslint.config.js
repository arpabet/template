import pluginVue from 'eslint-plugin-vue'
import prettier from 'eslint-config-prettier'

export default [
  {
    ignores: ['node_modules', 'dist', '../assets', 'public'],
  },
  ...pluginVue.configs['flat/recommended'],
  prettier,
  {
    languageOptions: {
      ecmaVersion: 'latest',
      sourceType: 'module',
      globals: {
        window: 'readonly',
        document: 'readonly',
        localStorage: 'readonly',
        setTimeout: 'readonly',
        console: 'readonly',
      },
    },
    rules: {
      'vue/multi-word-component-names': 'off',
      'vue/no-v-html': 'off',
      'vue/attribute-hyphenation': 'off',
      // The template ships a <Footer> component; harmless since PascalCase
      // tags always resolve to components in Vue 3.
      'vue/no-reserved-component-names': 'off',
    },
  },
]
