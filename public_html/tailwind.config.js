/** @type {import('tailwindcss').Config} */
export default {
  darkMode: true, // or 'media' or 'class'
  content: [
    './src/**/*.{html,js,svelte,ts}',
    'node_modules/preline/dist/*.js',
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('preline/plugin'),
  ],
}

