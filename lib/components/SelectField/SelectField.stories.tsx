import { Meta, StoryObj } from "@storybook/react";
import { SelectField, SelectOption } from "./SelectField.tsx";

const meta: Meta<typeof SelectField> = {
    title: "Components/Fields/Select Field",
    component: SelectField,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof SelectField>;

const children = [
    <SelectOption key="A">Option 1</SelectOption>,
    <SelectOption key="B">Option 2</SelectOption>,
    <SelectOption key="C">Option 3</SelectOption>,
];

export const Normal: Story = {
    args: {
        id: "select-field",
        label: "Label",
        required: true,
        description: "A short description",
        children,
    },
    argTypes: {
        onChange: { action: "onChange" },
    },
    render: (props) => <SelectField {...props} />,
};

export default meta;
