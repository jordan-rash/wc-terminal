const colors = require('tailwindcss/colors')

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.js"
  ],
  theme: {
    fontFamily: {
      'spacemono': ['Space Mono', 'monospace'],
    },
    colors:{
      'solarizefg': {
        light: '#657b83',
        DEFAULT: '#839496',
      },
      'solarizebg': {
        light: '#fdf6e3',
        DEFAULT: '#002b36',
      },
      'solarizehl': {
        light: '#eee8d5',
        DEFAULT: '#073642',
      },
      'solarizeemph': {
        light: '#586e75',
        DEFAULT: '#93a1a1',
      },
      'solarizecomment': {
        light: '#93a1a1',
        DEFAULT: '#586e75',
      },
      'yellow': '#b58900',
      'orange': '#cb4b16',
      'red'   : '#dc322f',
      'magenta':'#d33682',
      'violet': '#6c71c4',
      'blue':   '#268bd2',
      'cyan':   '#2aa198',
      'green':  '#859900',
    },
    extend: {
      colors: {},
    },
  },
  plugins: [],
}
