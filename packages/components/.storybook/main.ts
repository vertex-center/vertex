import type { StorybookConfig } from "@storybook/react-vite";
import * as path from "path";

const config: StorybookConfig = {
    stories: [
        "../lib/components/**/*.mdx",
        "../lib/components/**/*.stories.@(js|jsx|mjs|ts|tsx)",
    ],
    addons: [
        "@storybook/addon-links",
        "@storybook/addon-essentials",
        "@storybook/addon-onboarding",
        "@storybook/addon-interactions",
    ],
    framework: path.resolve(
        require.resolve("@storybook/react-vite"),
        "..",
    ) as any,
    // This doesn't work with Yarn workspaces, so we use the above workaround instead.
    // framework: {
    //     name: "@storybook/react-vite",
    //     options: {},
    // },
    docs: {
        autodocs: "tag",
        defaultName: "Documentation",
    },
};

export default config;
