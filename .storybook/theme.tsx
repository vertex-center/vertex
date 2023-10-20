import { StoryContext, StoryFn } from "@storybook/react";

export const withTheme = (Story: StoryFn, context: StoryContext) => {
    const theme =
        context?.globals?.backgrounds?.value === "#ffffff"
            ? "theme-vertex-light"
            : "theme-vertex-dark";

    return (
        <div className={theme}>
            <Story
                {...context}
                globals={{
                    ...context.globals,
                    theme: theme,
                }}
            />
        </div>
    );
};
