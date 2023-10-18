import { StoryContext, StoryFn } from "@storybook/react";
import ThemeProvider from "../lib/theme";

export const withThemeProvider = (Story: StoryFn, context: StoryContext) => {
    const theme =
        context?.globals?.backgrounds?.value === "#ffffff"
            ? "theme-vertex-light"
            : "theme-vertex-dark";

    console.log(theme);

    return (
        <ThemeProvider theme={theme}>
            <div className={theme}>
                <Story
                    {...context}
                    globals={{
                        ...context.globals,
                        theme: theme,
                    }}
                />
            </div>
        </ThemeProvider>
    );
};
