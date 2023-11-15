import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import mdx from "@mdx-js/rollup";
import remarkFrontmatter from "remark-frontmatter";
import remarkGfm from "remark-gfm";
import yaml from "@rollup/plugin-yaml";

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        react(),
        mdx({
            remarkPlugins: [remarkFrontmatter, remarkGfm],
        }),
        yaml(),
    ],
    base: "./",
});
