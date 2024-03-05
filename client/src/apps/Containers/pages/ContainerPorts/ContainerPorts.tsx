import {
    Button,
    Input,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeadCell,
    TableRow,
    Title,
    Vertical,
} from "@vertex-center/components";
import { useParams } from "react-router-dom";
import { useQueryClient } from "@tanstack/react-query";
import { useContainerPorts } from "../../hooks/useContainer";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { APIError } from "../../../../components/Error/APIError";
import { Horizontal } from "../../../../components/Layouts/Layouts";

export default function ContainerPorts() {
    const { uuid } = useParams();
    // const queryClient = useQueryClient();

    const { ports, isLoadingPorts, errorPorts } = useContainerPorts(uuid);

    const error = errorPorts;
    const isLoading = isLoadingPorts;

    return (
        <Vertical gap={24}>
            <Title variant="h2">Ports</Title>
            <APIError error={error} />
            <ProgressOverlay show={isLoading} />
            <Table>
                <TableHead>
                    <TableRow>
                        <TableHeadCell>Port inside container</TableHeadCell>
                        <TableHeadCell>Port outside container</TableHeadCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {ports?.map((port, i) => (
                        <TableRow>
                            <TableCell>
                                <Input value={port?.in} disabled />
                            </TableCell>
                            <TableCell>
                                <Input value={port?.out} disabled />
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </Vertical>
    );
}
