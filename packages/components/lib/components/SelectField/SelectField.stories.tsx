import { Meta, StoryObj } from "@storybook/react";
import { SelectField, SelectOption } from "./SelectField.tsx";
import { useState } from "react";

const meta: Meta<typeof SelectField> = {
    title: "Components/Fields/Select Field",
    component: SelectField,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof SelectField>;

export const Normal: Story = {
    args: {
        id: "select-field",
        required: true,
    },
    argTypes: {
        onChange: { action: "onChange" },
    },
    render: function Render(props) {
        const [value, setValue] = useState<string>();

        const onChange = (value: unknown) => {
            setValue(value as string);
            props.onChange?.(value);
        };

        return (
            <SelectField value={value} {...props} onChange={onChange}>
                <SelectOption value="1">One</SelectOption>
                <SelectOption value="2">Two</SelectOption>
                <SelectOption value="3">Three</SelectOption>
            </SelectField>
        );
    },
};

export default meta;
