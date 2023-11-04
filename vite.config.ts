import {defineConfig} from "vite";
import {resolve} from "path";
import react from "@vitejs/plugin-react";
import dts from "vite-plugin-dts";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        react(),
        dts({
            insertTypesEntry: true,
            exclude: ["**/*.stories.tsx", "**/*.test.tsx"],
        }),
    ],
    build: {
        lib: {
            entry: resolve(__dirname, "lib/index.ts"),
            name: "vertex-components",
            fileName: (format) => `vertex-components.${format}.js`,
        },
        copyPublicDir: false,
        rollupOptions: {
            external: ["react", "react-dom"],
            output: {
                globals: {
                    react: "React",
                    "react-dom": "ReactDOM",
                },
            },
        },
    },
});
