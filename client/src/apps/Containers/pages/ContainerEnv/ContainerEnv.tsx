import React, {
    ChangeEvent,
    Fragment,
    ReactNode,
    useEffect,
    useState,
} from "react";
import { useParams } from "react-router-dom";
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
import { Horizontal } from "../../../../components/Layouts/Layouts";
import { useRecreateContainer } from "../../hooks/useContainer";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useQueryClient } from "@tanstack/react-query";
import Content from "../../../../components/Content/Content";
import { EnvVariable } from "../../backend/models";
import {
    ArrowUUpLeft,
    FloppyDiskBack,
    Plus,
    Textbox,
    Trash,
} from "@phosphor-icons/react";
import { Controller, useFieldArray, useForm } from "react-hook-form";
import { diffArrays, diffJson } from "diff";
import NoItems from "../../../../components/NoItems/NoItems";
import {
    useCreateEnv,
    useDeleteEnv,
    useEnvironment,
    usePatchEnv,
} from "../../hooks/useEnvironment";
import Spacer from "../../../../components/Spacer/Spacer";

type EnvTableProps = {
    env: EnvVariable[];
};

function EnvTable(props: EnvTableProps) {
    const queryClient = useQueryClient();

    const { uuid } = useParams();
    const { env } = props;

    if (env === undefined) return;

    const {
        control,
        handleSubmit,
        reset,
        formState: { isDirty },
    } = useForm({
        defaultValues: { env },
    });

    useEffect(() => {
        reset({ env });
    }, [env]);

    const { fields, append, remove } = useFieldArray({
        control,
        name: "env",
        keyName: "_id",
    });

    const { patchEnvAsync } = usePatchEnv();
    const { deleteEnvAsync } = useDeleteEnv();
    const { createEnvAsync } = useCreateEnv();

    const { recreateContainer, isPendingRecreate, errorRecreate } =
        useRecreateContainer();

    const onAdd = () => {
        append({
            id: `TEMP_${Date.now()}`,
            container_id: uuid,
            type: "string",
            name: "",
            value: "",
        });
    };

    const onSubmit = handleSubmit(async (d) => {
        const _env = env === null ? [] : env;
        let patch = diffArrays(_env, d.env, {
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

        const requests = [];
        for (const p of d.env) {
            console.log(p);
            if (modified.has(p.id)) {
                requests.push(patchEnvAsync(p));
            } else if (added.has(p.id)) {
                requests.push(createEnvAsync(p));
            }
        }
        for (const p of _env) {
            if (deleted.has(p.id)) {
                requests.push(deleteEnvAsync(p.id));
            }
        }
        await Promise.all(requests);
        await queryClient.invalidateQueries({
            queryKey: ["environments"],
        });
        recreateContainer(uuid);
    });

    const isLoading = isPendingRecreate;
    const error = errorRecreate;

    let table: ReactNode;
    if (fields?.length === 0) {
        table = (
            <NoItems
                icon={<Textbox />}
                text="This container has no environment variables."
            />
        );
    } else {
        table = (
            <Table>
                <TableHead>
                    <TableRow>
                        <TableHeadCell>Name</TableHeadCell>
                        <TableHeadCell>Value</TableHeadCell>
                        <TableHeadCell />
                    </TableRow>
                </TableHead>
                <TableBody>
                    {fields?.map((port, i) => (
                        <TableRow key={port.id}>
                            <TableCell>
                                <Controller
                                    control={control}
                                    name={`env.${i}.name`}
                                    render={({ field }) => <Input {...field} />}
                                />
                            </TableCell>
                            <TableCell>
                                <Controller
                                    control={control}
                                    name={`env.${i}.value`}
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
                            Add variable
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

export default function ContainerEnv() {
    const { uuid } = useParams();

    const { env, isLoadingEnv, errorEnv } = useEnvironment({
        container_id: uuid,
    });

    return (
        <Content>
            <Title variant="h2">Environment</Title>
            <ProgressOverlay show={isLoadingEnv} />
            <APIError error={errorEnv} />
            <EnvTable env={env} />
        </Content>
    );
}
