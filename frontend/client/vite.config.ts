import react from "@vitejs/plugin-react";
import vike from "vike/plugin";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react({}), vike({ prerender: true })],
  server: {
    proxy: {
      "/__/firebase": {
        target: "https://tasuke-dev.web.app",
        changeOrigin: true,
      },
    },
  },
});
