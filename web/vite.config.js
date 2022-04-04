import { defineConfig } from "vite";
import tsConfigPaths from "vite-tsconfig-paths";
import mpa from "vite-plugin-mpa";

export default defineConfig({
  plugins: [tsConfigPaths(), mpa({ open: false })],
});
