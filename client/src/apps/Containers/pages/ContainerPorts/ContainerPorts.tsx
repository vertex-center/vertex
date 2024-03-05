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
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { APIError } from "../../../../components/Error/APIError";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import { Controller, useFieldArray, useForm } from "react-hook-form";
import { ArrowUUpLeft, FloppyDiskBack } from "@phosphor-icons/react";
import { API } from "../../backend/api";
import { useSaveContainerPorts } from "../../hooks/useContainer";

export default function ContainerPorts() {
    const { uuid } = useParams();
    const { control, handleSubmit, reset } = useForm({
        defaultValues: async () => {
            const ports = await API.getContainerPorts(uuid);
            return { ports };
        },
    });

    const { fields } = useFieldArray({
        control,
        name: "ports",
    });

    const { savePorts, isPending, error } = useSaveContainerPorts(uuid, {
        onSuccess: () => {
            reset({}, { keepValues: true });
        },
    });

    const onSubmit = handleSubmit((data) => {
        savePorts(data.ports);
    });

    const isLoading = isPending;

    return (
        <Vertical gap={24}>
            <Title variant="h2">Ports</Title>
            <APIError error={error} />
            <ProgressOverlay show={isLoading} />
            <form onSubmit={onSubmit}>
                <Vertical gap={24}>
                    <Table>
                        <TableHead>
                            <TableRow>
                                <TableHeadCell>
                                    Port inside container
                                </TableHeadCell>
                                <TableHeadCell>
                                    Port outside container
                                </TableHeadCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {fields?.map((port, i) => (
                                <TableRow key={port.id}>
                                    <TableCell>
                                        <Controller
                                            control={control}
                                            name={`ports.${i}.in`}
                                            render={({
                                                field,
                                                formState: { dirtyFields },
                                            }) => (
                                                <Input
                                                    {...field}
                                                    style={{
                                                        color:
                                                            dirtyFields
                                                                ?.ports?.[
                                                                `${i}`
                                                            ]?.in &&
                                                            "var(--blue)",
                                                    }}
                                                />
                                            )}
                                        />
                                    </TableCell>
                                    <TableCell>
                                        <Controller
                                            control={control}
                                            name={`ports.${i}.out`}
                                            render={({
                                                field,
                                                formState: { dirtyFields },
                                            }) => (
                                                <Input
                                                    {...field}
                                                    style={{
                                                        color:
                                                            dirtyFields
                                                                ?.ports?.[
                                                                `${i}`
                                                            ]?.out &&
                                                            "var(--blue)",
                                                    }}
                                                />
                                            )}
                                        />
                                    </TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                    <Horizontal justifyContent="flex-end" gap={10}>
                        <Button
                            type="reset"
                            variant="outlined"
                            onClick={() => reset()}
                            rightIcon={<ArrowUUpLeft />}
                            disabled={isLoading}
                        >
                            Cancel
                        </Button>
                        <Button
                            type="submit"
                            variant="colored"
                            rightIcon={<FloppyDiskBack />}
                            disabled={isLoading}
                        >
                            Save
                        </Button>
                    </Horizontal>
                </Vertical>
            </form>
        </Vertical>
    );
}
