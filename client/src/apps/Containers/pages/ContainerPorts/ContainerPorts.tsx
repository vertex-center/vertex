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
    Trash,
} from "@phosphor-icons/react";
import { useContainerPorts } from "../../hooks/useContainer";
import Spacer from "../../../../components/Spacer/Spacer";
import { Fragment, ReactNode, useEffect } from "react";
import { Port } from "../../backend/models";
import NoItems from "../../../../components/NoItems/NoItems";
import { diffArrays, diffJson } from "diff";
import { usePatchPort } from "../../hooks/usePort";

type PortTableProps = {
    ports: Port[];
};

function PortTable(props: PortTableProps) {
    const { uuid } = useParams();
    const { ports } = props;

    if (!ports) return;

    const {
        control,
        handleSubmit,
        reset,
        formState: { isDirty },
    } = useForm({
        defaultValues: { ports },
    });

    useEffect(() => {
        reset({ ports });
    }, [ports]);

    const { fields, append, remove } = useFieldArray({
        control,
        name: "ports",
        keyName: "_id",
    });

    const { patchPortAsync, isPending, error } = usePatchPort({
        onSuccess: () => {
            // reset({}, { keepValues: true });
        },
    });

    const onAdd = () => {
        append({
            id: `TEMP_${Date.now()}`,
            container_id: uuid,
            in: "",
            out: "",
        });
    };

    const onSubmit = handleSubmit(async (d) => {
        let patch = diffArrays(ports, d.ports, {
            comparator: (a, b) => diffJson(a, b).length === 1,
        });

        const _deleted = new Set(
            patch
                .filter((p) => p.removed)
                .map((p) => p.value)
                .flat()
                .map((p) => p.id)
        );
        const _added = new Set(
            patch
                .filter((p) => p.added)
                .map((p) => p.value)
                .flat()
                .map((p) => p.id)
        );

        const modified = new Set([..._deleted].filter((x) => _added.has(x)));
        const deleted = new Set([..._deleted].filter((x) => !modified.has(x)));
        const added = new Set([..._added].filter((x) => !modified.has(x)));

        for (const p of d.ports) {
            if (modified.has(p.id)) {
                await patchPortAsync(p);
            }
        }
    });

    const isLoading = isPending;

    let table: ReactNode;
    if (fields?.length === 0) {
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
                        <TableHeadCell />
                    </TableRow>
                </TableHead>
                <TableBody>
                    {fields?.map((port, i) => (
                        <TableRow key={port.id}>
                            <TableCell>
                                <Controller
                                    control={control}
                                    name={`ports.${i}.in`}
                                    render={({ field }) => <Input {...field} />}
                                />
                            </TableCell>
                            <TableCell>
                                <Controller
                                    control={control}
                                    name={`ports.${i}.out`}
                                    render={({ field }) => <Input {...field} />}
                                />
                            </TableCell>
                            <TableCell right>
                                <Button
                                    type="button"
                                    onClick={() => remove(i)}
                                    variant="danger"
                                    borderless
                                    disabled={isLoading}
                                    rightIcon={<Trash />}
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
                                    Save changes
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

    const { ports, isLoadingPorts, errorPorts } = useContainerPorts(uuid);

    return (
        <Vertical gap={24}>
            <Title variant="h2">Ports</Title>
            <APIError error={errorPorts} />
            <ProgressOverlay show={isLoadingPorts} />
            <PortTable ports={ports} />
        </Vertical>
    );
}
