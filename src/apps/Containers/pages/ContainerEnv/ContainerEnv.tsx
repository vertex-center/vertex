import React, { ChangeEvent, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import {
    Button,
    Input,
    MaterialIcon,
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeadCell,
    TableRow,
    Title,
} from "@vertex-center/components";
import { Horizontal } from "../../../../components/Layouts/Layouts";
import { useContainerEnv } from "../../hooks/useContainer";
import { APIError } from "../../../../components/Error/APIError";
import { ProgressOverlay } from "../../../../components/Progress/Progress";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import Content from "../../../../components/Content/Content";
import { API } from "../../backend/api";
import { EnvVariables } from "../../backend/models";
import styles from "./ContainerEnv.module.sass";

export default function ContainerEnv() {
    const { uuid } = useParams();
    const queryClient = useQueryClient();

    const { env: currentEnv, isLoadingEnv, errorEnv } = useContainerEnv(uuid);
    const [env, setEnv] = useState<EnvVariables>();

    useEffect(() => {
        if (!currentEnv) return;
        setEnv(JSON.parse(JSON.stringify(currentEnv)));
        setSaved(true);
    }, [currentEnv]);

    const [saved, setSaved] = useState<boolean>(true);

    const mutationSaveEnv = useMutation({
        mutationFn: async (env: EnvVariables) => {
            await API.saveEnv(uuid, env);
        },
        onSuccess: () => setSaved(true),
        onSettled: () => {
            queryClient.invalidateQueries({
                queryKey: ["containers", uuid],
            });
            queryClient.invalidateQueries({
                queryKey: ["container_env", uuid],
            });
        },
    });
    const { isLoading: isUploading } = mutationSaveEnv;

    const save = () => {
        let patch = [...env];
        patch = patch.filter(
            (env, i) =>
                env.name !== currentEnv[i].name ||
                env.value !== currentEnv[i].value
        );
        mutationSaveEnv.mutate(patch);
    };

    const onNameChange = (i: number, e: ChangeEvent<HTMLInputElement>) => {
        const newEnv = [...env];
        newEnv[i].name = e.target.value;
        updateEnv(newEnv);
    };

    const onValueChange = (i: number, e: ChangeEvent<HTMLInputElement>) => {
        const newEnv = [...env];
        newEnv[i].value = e.target.value;
        updateEnv(newEnv);
    };

    const updateEnv = (env: EnvVariables) => {
        setEnv(env);
        setSaved(isSaved());
    };

    const isSaved = () => {
        for (let i = 0; i < env.length; i++) {
            if (env[i].value !== currentEnv[i].value) return false;
            if (env[i].name !== currentEnv[i].name) return false;
        }
        return true;
    };

    return (
        <Content>
            <Title variant="h2">Environment</Title>
            <Table>
                <TableHead>
                    <TableRow>
                        <TableHeadCell>Name</TableHeadCell>
                        <TableHeadCell>Value</TableHeadCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    {env?.map((env, i) => (
                        <TableRow key={env.name}>
                            <TableCell>
                                <Input
                                    value={env.name}
                                    name={currentEnv[i].name + "_name"}
                                    onChange={(e) => onNameChange(i, e)}
                                    disabled={isUploading}
                                    className={styles.input}
                                    style={{
                                        color:
                                            env.name !== currentEnv[i].name &&
                                            "var(--blue)",
                                    }}
                                />
                            </TableCell>
                            <TableCell>
                                <Input
                                    value={env.value}
                                    name={currentEnv[i].name}
                                    placeholder={env.default}
                                    onChange={(e) => onValueChange(i, e)}
                                    type={env.secret ? "password" : undefined}
                                    disabled={isUploading}
                                    className={styles.input}
                                    style={{
                                        color:
                                            env.value !== currentEnv[i].value &&
                                            "var(--blue)",
                                    }}
                                />
                            </TableCell>
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
            <ProgressOverlay show={isLoadingEnv ?? isUploading} />
            <Horizontal justifyContent="flex-end">
                <Button
                    variant="colored"
                    onClick={save}
                    rightIcon={<MaterialIcon icon="save" />}
                    disabled={isUploading || saved || saved === undefined}
                >
                    Save
                </Button>
            </Horizontal>
            <APIError error={errorEnv} />
        </Content>
    );
}
