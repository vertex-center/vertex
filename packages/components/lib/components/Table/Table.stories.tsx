import { Meta, StoryObj } from "@storybook/react";
import { Table } from "./Table.tsx";

const meta: Meta<typeof Table> = {
    title: "Components/Table",
    component: Table,
    tags: ["autodocs"],
};

type Story = StoryObj<typeof Table>;

export const Normal: Story = {
    render: function Render(props) {
        return (
            <Table {...props}>
                <thead>
                    <tr>
                        <th>Header 1</th>
                        <th>Header 2</th>
                    </tr>
                </thead>
                <tbody>
                    <tr>
                        <td>Cell A1</td>
                        <td>Cell B1</td>
                    </tr>
                    <tr>
                        <td>Cell A2</td>
                        <td>Cell B2</td>
                    </tr>
                </tbody>
            </Table>
        );
    },
};

export default meta;
