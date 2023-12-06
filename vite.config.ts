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

function transformYaml(data, filePath) {
    if (filePath.endsWith("_category_.yml")) {
        return data;
    }
    // Commented for now as it causes issues.
    //
    // const regex = /\{"\$ref":"#\/(.*?)"}/g;
    // let dataString = JSON.stringify(data);
    // let match = null;
    // while ((match = regex.exec(dataString)) !== null) {
    //     if (match.index === regex.lastIndex) {
    //         regex.lastIndex++;
    //     }
    //     const ref = match[1];
    //     const refPath = ref.split("/");
    //     const refData = refPath.reduce((d, key) => d[key], data);
    //     dataString = dataString.replace(match[0], JSON.stringify(refData));
    //     data = JSON.parse(dataString);
    //     regex.lastIndex = 0;
    // }
    return JSON.parse(JSON.stringify(data));
}

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
        yaml({ transform: transformYaml }),
    ],
});
