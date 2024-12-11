/* to disable code block CSS from tailwind/typography, we use another code highlighter */
const disabledCss = {
  "code::before": false,
  "code::after": false,
  "blockquote p:first-of-type::before": false,
  "blockquote p:last-of-type::after": false,
  pre: false,
  code: false,
  "pre code": false,
  "code::before": false,
  "code::after": false,
};

module.exports = {
  content: ["internal/templates/**/*.html"],
  theme: {
    extend: {
      /* disable code block CSS */
      typography: {
        DEFAULT: { css: disabledCss },
        sm: { css: disabledCss },
        lg: { css: disabledCss },
        xl: { css: disabledCss },
        "2xl": { css: disabledCss },
      },
    },
  },
  /* @tailwind/typography plugin should be befire daisyui */
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: [
      {
        light: {
          primary: "#223843",
          accent: "#D77A61",
          neutral: "#DBD3D8",
          "base-100": "#EFF1F3",
          success: "#15803d",
          info: "#0369a1",
          error: "#be123c",
          warning: "#b45309",
          "--rounded-box": "0.50rem",
          "--rounded-btn": "0.50rem",
        },
      },
    ],
  },
};
