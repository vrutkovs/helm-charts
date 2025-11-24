import { defineConfig } from "allure";

export default defineConfig({
  name: "Allure Report Example",
  output: "/tmp/allure-report",
  plugins: {
    awesome: {
      options: {
        singleFile: true,
        reportLanguage: "en",
      },
    },
  },
});
