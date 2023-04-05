/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{svelte,html,ts,js}'],
  theme: {
    extend: {},
  },
  daisyui: {
    themes:
    {
      "primary": "#65C3C8",
      "secondary": "#EF9FBC",
      "accent": "#EEAF3A",
      "neutral": "#291334",
      "base-100": "#FAF7F5",
      "info": "#3ABFF8",
      "success": "#36D399",
      "warning": "#FBBD23",
      "error": "#F87272",
    }
  },
  plugins: [require("daisyui")],
  
}
