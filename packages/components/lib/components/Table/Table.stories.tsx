import { Meta, StoryObj } from "@storybook/react";
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeadCell,
    TableRow,
} from "./Table.tsx";

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
                <TableHead>
                    <TableRow>
                        <TableHeadCell>Header 1</TableHeadCell>
                        <TableHeadCell>Header 2</TableHeadCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    <TableRow>
                        <TableCell>Cell A1</TableCell>
                        <TableCell>Cell B1</TableCell>
                    </TableRow>
                    <TableRow>
                        <TableCell>Cell A2</TableCell>
                        <TableCell>Cell B2</TableCell>
                    </TableRow>
                </TableBody>
            </Table>
        );
    },
};

export default meta;
