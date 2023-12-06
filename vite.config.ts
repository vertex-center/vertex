import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import mdx from "@mdx-js/rollup";
import remarkFrontmatter from "remark-frontmatter";
import remarkGfm from "remark-gfm";
import remarkDirective from "remark-directive";
import remarkMermaid from "remark-mermaidjs";
import yaml from "@rollup/plugin-yaml";
import { visit } from "unist-util-visit";
import SwaggerParser from "@apidevtools/swagger-parser";

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

async function transformYaml(data, filePath) {}

const openapi = async () => {
    return {
        name: "openapi",
        transform: async (src, id) => {
            if (!id.endsWith(".yml") && !id.endsWith(".yaml")) {
                return src;
            }
            if (id.endsWith("_category_.yml")) {
                return src;
            }
            const api = await SwaggerParser.dereference(id);
            return `export default ${JSON.stringify(api)}`;
        },
    };
};

export default defineConfig(async () => ({
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
        await openapi(),
    ],
}));
