import { Meta, StoryObj } from "@storybook/react";
import Code from "./Code.tsx";

const exampleCode = `import React from "react";
import { render } from "react-dom";
import { App } from "./App";
import "./index.css";
`;

const meta: Meta = {
    title: "Components/Code",
    component: Code,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Code>;

export const Default: Story = {
    args: {
        code: exampleCode,
        language: "javascript",
    },
    render: function Render(props) {
        return <Code {...props} />;
    },
};

export default meta;
