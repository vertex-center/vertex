import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import mdx from "@mdx-js/rollup";
import remarkFrontmatter from "remark-frontmatter";
import remarkGfm from "remark-gfm";
import remarkDirective from "remark-directive";
import remarkMermaid from "remark-mermaidjs";
import yaml from "@rollup/plugin-yaml";
import { visit } from "unist-util-visit";

const remarkDirectiveBlocks = () => {
    return (tree) => {
        visit(tree, (node) => {
            if (
                node.type === "containerDirective" ||
                node.type === "leafDirective" ||
                node.type === "textDirective"
            ) {
                const data = node.data || (node.data = {});
                data.hName = node.name;
            }
        });
    };
};

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        react(),
        mdx({
            remarkPlugins: [
                remarkFrontmatter,
                remarkGfm,
                remarkDirective,
                remarkDirectiveBlocks,
                [
                    remarkMermaid,
                    {
                        mermaidConfig: {
                            theme: "dark",
                        },
                    },
                ],
            ],
        }),
        yaml(),
    ],
});
