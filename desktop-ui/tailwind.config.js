/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: 'oklch(0.98 0.01 270)',
          100: 'oklch(0.95 0.03 270)',
          200: 'oklch(0.88 0.06 270)',
          300: 'oklch(0.78 0.10 270)',
          400: 'oklch(0.68 0.15 270)',
          500: 'oklch(0.60 0.18 270)',
          600: 'oklch(0.52 0.16 270)',
          700: 'oklch(0.45 0.14 270)',
          800: 'oklch(0.38 0.12 270)',
          900: 'oklch(0.32 0.10 270)',
        },
        secondary: {
          50: 'oklch(0.95 0.02 200)',
          100: 'oklch(0.88 0.05 200)',
          200: 'oklch(0.80 0.08 200)',
          300: 'oklch(0.72 0.12 200)',
          400: 'oklch(0.65 0.15 200)',
          500: 'oklch(0.60 0.16 200)',
          600: 'oklch(0.52 0.14 200)',
          700: 'oklch(0.45 0.12 200)',
          800: 'oklch(0.38 0.10 200)',
          900: 'oklch(0.32 0.08 200)',
        },
      },
    },
  },
}
