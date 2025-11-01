/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/web/templates/**/*.{html,js}", "./src/web/templates/*.{html,js}"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["light", "dark", "dim", "nord", "winter", "corporate"],
  },
}
