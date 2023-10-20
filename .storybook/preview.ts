import type { Preview } from "@storybook/react";

import "../lib/index.sass";
import { withTheme } from "./theme";

const preview: Preview = {
    parameters: {
        actions: { argTypesRegex: "^on[A-Z].*" },
        layout: "centered",
        controls: {
            matchers: {
                color: /(background|color)$/i,
                date: /Date$/i,
            },
        },
        backgrounds: {
            default: "theme-vertex-dark",
            values: [
                {
                    name: "theme-vertex-dark",
                    value: "#111111",
                },
                {
                    name: "theme-vertex-light",
                    value: "#ffffff",
                },
            ],
        },
    },
    decorators: [withTheme],
};

export default preview;
