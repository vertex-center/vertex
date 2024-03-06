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
import {
    ArrowUUpLeft,
    FloppyDiskBack,
    Plus,
    ShareNetwork,
} from "@phosphor-icons/react";
import {
    useContainerPorts,
    useSaveContainerPorts,
} from "../../hooks/useContainer";
import Spacer from "../../../../components/Spacer/Spacer";
import { Fragment } from "react";
import { Port } from "../../backend/models";
import NoItems from "../../../../components/NoItems/NoItems";

type PortTableProps = {
    ports: Port[];
};

function PortTable(props: PortTableProps) {
    const { uuid } = useParams();
    const { ports } = props;

    if (!ports) {
        return;
    }

    const {
        control,
        handleSubmit,
        reset,
        formState: { isDirty },
    } = useForm({
        defaultValues: { ports },
    });

    const { fields, append } = useFieldArray({
        control,
        name: "ports",
    });

    const { savePorts, isPending, error } = useSaveContainerPorts(uuid, {
        onSuccess: () => {
            reset({}, { keepValues: true });
        },
    });

    const onAdd = () => append({ container_id: uuid, in: "", out: "" });
    const onSubmit = handleSubmit((d) => savePorts(d.ports));

    const isLoading = isPending;

    let table;
    if (ports && ports?.length === 0) {
        table = (
            <NoItems
                icon={<ShareNetwork />}
                text="This container has no exposed ports."
            />
        );
    } else {
        table = (
            <Table>
                <TableHead>
                    <TableRow>
                        <TableHeadCell>Port inside container</TableHeadCell>
                        <TableHeadCell>Port outside container</TableHeadCell>
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
                                                    dirtyFields?.ports?.[`${i}`]
                                                        ?.in && "var(--blue)",
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
                                                    dirtyFields?.ports?.[`${i}`]
                                                        ?.out && "var(--blue)",
                                            }}
                                        />
                                    )}
                                />
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        );
    }

    return (
        <Fragment>
            <APIError error={error} />
            <ProgressOverlay show={isLoading} />
            <form onSubmit={onSubmit}>
                <Vertical gap={12}>
                    <Horizontal justifyContent="flex-end" gap={10}>
                        <Button
                            type="button"
                            variant="outlined"
                            onClick={onAdd}
                            rightIcon={<Plus />}
                            disabled={isLoading}
                        >
                            Add port
                        </Button>
                        <Spacer />
                        {isDirty && (
                            <Fragment>
                                <Button
                                    type="reset"
                                    variant="outlined"
                                    onClick={() => reset()}
                                    rightIcon={<ArrowUUpLeft />}
                                    disabled={isLoading}
                                />
                                <Button
                                    type="submit"
                                    variant="colored"
                                    rightIcon={<FloppyDiskBack />}
                                    disabled={isLoading}
                                >
                                    Save
                                </Button>
                            </Fragment>
                        )}
                    </Horizontal>
                    {table}
                </Vertical>
            </form>
        </Fragment>
    );
}

export default function ContainerPorts() {
    const { uuid } = useParams();

    const { ports, dataUpdatedAt, isLoadingPorts, errorPorts } =
        useContainerPorts(uuid);

    return (
        <Vertical gap={24}>
            <Title variant="h2">Ports</Title>
            <APIError error={errorPorts} />
            <ProgressOverlay show={isLoadingPorts} />
            <PortTable key={dataUpdatedAt} ports={ports} />
        </Vertical>
    );
}
