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
        target: "https://alpha.tasuke.dev",
        changeOrigin: true,
      },
      "/frontendapi.FrontendService": {
        target: "https://alpha.tasuke.dev",
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
