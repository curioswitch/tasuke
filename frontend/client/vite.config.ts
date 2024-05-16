import path from "node:path";
import react from "@vitejs/plugin-react";
import vike from "vike/plugin";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [react({}), vike({ prerender: true })],
  server: {
    port: 8080,
    proxy: {
      "/__/firebase": {
        target: "https://tasuke-dev.web.app",
        changeOrigin: true,
      },
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
